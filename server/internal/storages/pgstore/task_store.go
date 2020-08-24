package pgstore

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/jackc/pgx/v4"
	"github.com/jackc/pgx/v4/pgxpool"

	"github.com/kulti/task-list/server/internal/models"
	"github.com/kulti/task-list/server/internal/storages"
)

type TaskStore struct {
	conn *pgxpool.Pool
}

func New(dbURL string) (*TaskStore, error) {
	conn, err := pgxpool.Connect(context.Background(), dbURL)
	if err != nil {
		return nil, fmt.Errorf("unable to connect to database: %w", err)
	}

	s := &TaskStore{
		conn: conn,
	}

	return s, nil
}

func (s *TaskStore) Close() {
	if s.conn != nil {
		s.conn.Close()
	}
}

func (s *TaskStore) NewSprint(ctx context.Context, opts storages.SprintOpts) error {
	_, err := s.conn.Exec(ctx,
		"INSERT INTO task_lists (type, title, created_at, begin, \"end\") VALUES ($1, $2, $3, $4, $5)",
		"sprint", opts.Title, time.Now(), opts.Begin, opts.End)
	return err
}

func (s *TaskStore) CreateTask(ctx context.Context, task storages.Task, sprintIDStr string) (int64, error) {
	row := s.conn.QueryRow(ctx,
		"INSERT INTO tasks (text, points, burnt, state) VALUES($1, $2, $3, $4) RETURNING id",
		task.Text, task.Points, task.Burnt, task.State)
	var taskID int64
	err := row.Scan(&taskID)
	if err != nil {
		return -1, fmt.Errorf("failed to create task: %w", err)
	}

	var sprintID int64
	if sprintIDStr == "current" {
		row = s.conn.QueryRow(ctx,
			"SELECT id FROM task_lists WHERE type = 'sprint' ORDER BY created_at DESC LIMIT 1")

		err = row.Scan(&sprintID)
		if err != nil {
			return -1, fmt.Errorf("failed to find sprint: %w", err)
		}
	}

	_, err = s.conn.Exec(ctx,
		"INSERT INTO task_list_map (task_id, list_id) VALUES ($1, $2)",
		taskID, sprintID)
	if err != nil {
		return -1, fmt.Errorf("failed to add task to list: %w", err)
	}

	return taskID, nil
}

func (s *TaskStore) DeleteTask(ctx context.Context, taskID int64) error {
	_, err := s.conn.Exec(ctx, "DELETE FROM tasks WHERE id = $1", taskID)
	return err
}

func (s *TaskStore) ListTasks(ctx context.Context, sprintIDStr string) (storages.TaskList, error) {
	if sprintIDStr != "current" {
		return storages.TaskList{}, nil
	}

	var taskList storages.TaskList
	var sprintID int64

	row := s.conn.QueryRow(ctx,
		"SELECT id, title FROM task_lists WHERE type = 'sprint' ORDER BY created_at DESC LIMIT 1")

	err := row.Scan(&sprintID, &taskList.Title)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return storages.TaskList{
				Title: "No any sprint. Create one first.",
			}, nil
		}
		return storages.TaskList{}, fmt.Errorf("failed to find list: %w", err)
	}

	rows, _ := s.conn.Query(ctx,
		`SELECT tasks.id, tasks.text, tasks.points, tasks.burnt, tasks.state
		FROM tasks, task_list_map
		WHERE task_list_map.list_id = $1
			AND tasks.id = task_list_map.task_id
		ORDER BY tasks.id`,
		sprintID)
	defer rows.Close()

	err = rows.Err()
	for err == nil && rows.Next() {
		var task storages.Task
		err = rows.Scan(&task.ID, &task.Text, &task.Points, &task.Burnt, &task.State)
		if err == nil {
			taskList.Tasks = append(taskList.Tasks, task)
		}
	}

	if err != nil {
		return storages.TaskList{}, err
	}

	return taskList, nil
}

func (s *TaskStore) UpdateTask(ctx context.Context, taskID int64, opts models.UpdateOptions) error {
	_, err := s.conn.Exec(ctx,
		"UPDATE tasks SET text = $2, points = $3, burnt = $4 WHERE id = $1",
		taskID, opts.Text, opts.Points, opts.Burnt)
	return err
}

func (s *TaskStore) DoneTask(ctx context.Context, taskID int64) error {
	return s.updateTaskStateWithStmt(ctx, taskID, models.DoneTaskEvent,
		"UPDATE tasks SET state = $2, burnt=points WHERE id = $1")
}

func (s *TaskStore) UndoneTask(ctx context.Context, taskID int64) error {
	return s.updateTaskState(ctx, taskID, models.UndoneTaskEvent)
}

func (s *TaskStore) TodoTask(ctx context.Context, taskID int64) error {
	return s.updateTaskState(ctx, taskID, models.TodoTaskEvent)
}

func (s *TaskStore) CancelTask(ctx context.Context, taskID int64) error {
	return s.updateTaskState(ctx, taskID, models.CancelTaskEvent)
}

func (s *TaskStore) BackTaskToWork(ctx context.Context, taskID int64) error {
	return s.updateTaskState(ctx, taskID, models.ToWorkTaskEvent)
}

func (s *TaskStore) PostponeTask(ctx context.Context, taskID int64) (resultErr error) {
	return s.updateTaskStateWithStmt(ctx, taskID, models.PostponeTaskEvent,
		`WITH task AS (DELETE FROM tasks WHERE id = $1 AND $2 = $2 RETURNING *)
		INSERT INTO postponed_tasks (text, points)
		SELECT text, points
		FROM task`)
}

func (s *TaskStore) updateTaskState(ctx context.Context, taskID int64,
	event models.SwitchTaskStateEvent,
) error {
	return s.updateTaskStateWithStmt(ctx, taskID, event, "UPDATE tasks SET state = $2 WHERE id = $1")
}

func (s *TaskStore) updateTaskStateWithStmt(ctx context.Context, taskID int64,
	event models.SwitchTaskStateEvent, stmt string,
) (resultErr error) {
	tx, err := s.conn.Begin(ctx)
	if err != nil {
		return fmt.Errorf("failed to begin transaction: %w", err)
	}

	defer func() {
		rollbackErr := tx.Rollback(ctx)
		if resultErr == nil && !errors.Is(rollbackErr, pgx.ErrTxClosed) {
			resultErr = fmt.Errorf("failed to rollback: %w", rollbackErr)
		}
	}()

	row := tx.QueryRow(ctx, "SELECT state FROM tasks WHERE id = $1 FOR NO KEY UPDATE", taskID)
	var curState models.TaskState
	err = row.Scan(&curState)
	if err == nil || errors.Is(err, pgx.ErrNoRows) {
		state, err := curState.NextState(event)
		if err != nil {
			return err
		}

		if _, err := tx.Exec(ctx, stmt, taskID, state); err != nil {
			return fmt.Errorf("failed to execute update: %w", err)
		}
	} else if err != nil {
		return fmt.Errorf("failed to scan: %w", err)
	}

	return tx.Commit(ctx)
}

func (s *TaskStore) GetSprintTemplate(ctx context.Context) (models.SprintTemplate, error) {
	stmts := []string{
		`SELECT text, points
		FROM new_sprint_task_tempate
		ORDER BY id`,
		`DELETE FROM postponed_tasks RETURNING text, points`,
	}

	var tasks []models.TaskTemplate

	for _, stmt := range stmts {
		rows, _ := s.conn.Query(ctx, stmt)
		err := rows.Err()
		for err == nil && rows.Next() {
			var task models.TaskTemplate
			err = rows.Scan(&task.Text, &task.Points)
			if err == nil {
				tasks = append(tasks, task)
			}
		}
		rows.Close()

		if err != nil {
			return models.SprintTemplate{}, err
		}
	}

	return models.SprintTemplate{Tasks: tasks}, nil
}
