package pgstore

import (
	"context"
	"errors"
	"fmt"
	"strconv"
	"time"

	"github.com/jackc/pgx/v4"
	"github.com/jackc/pgx/v4/pgxpool"
	"github.com/kulti/task-list/server/internal/models"
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

func (s *TaskStore) NewSprint(ctx context.Context, opts models.SprintOpts) error {
	_, err := s.conn.Exec(ctx,
		"INSERT INTO task_lists (type, title, created_at, begin, \"end\") VALUES ($1, $2, $3, $4, $5)",
		"sprint", opts.Title, time.Now(), opts.Begin, opts.End)
	return err
}

func (s *TaskStore) CreateTask(ctx context.Context, task models.Task, listType string) (string, error) {
	row := s.conn.QueryRow(ctx,
		"INSERT INTO tasks (text, points, burnt, state) VALUES($1, $2, $3, $4) RETURNING id",
		task.Text, task.Points, task.Burnt, task.State)
	var taskID int64
	err := row.Scan(&taskID)
	if err != nil {
		return "", fmt.Errorf("failed to create task: %w", err)
	}

	row = s.conn.QueryRow(ctx,
		"SELECT id FROM task_lists WHERE type = $1 ORDER BY created_at DESC LIMIT 1",
		listType)

	var listID int
	err = row.Scan(&listID)
	if err != nil {
		return "", fmt.Errorf("failed to find list: %w", err)
	}

	_, err = s.conn.Exec(ctx,
		"INSERT INTO task_list_map (task_id, list_id) VALUES ($1, $2)",
		taskID, listID)
	if err != nil {
		return "", fmt.Errorf("failed to add task to list: %w", err)
	}

	return strconv.FormatInt(taskID, 16), nil
}

func (s *TaskStore) DeleteTaskFromList(ctx context.Context, taskID, listType string) error {
	id, err := strconv.ParseInt(taskID, 16, 64)
	if err != nil {
		return err
	}
	switch listType {
	case "sprint":
		_, err = s.conn.Exec(ctx, "DELETE FROM tasks WHERE id = $1", id)
	default:
		err = errUnknownListType
	}
	return err
}

func (s *TaskStore) ListTasks(ctx context.Context, listType string) (models.TaskList, error) {
	var taskList models.TaskList
	var listID int64

	row := s.conn.QueryRow(ctx,
		"SELECT id, title FROM task_lists WHERE type = $1 ORDER BY created_at DESC LIMIT 1",
		listType)

	err := row.Scan(&listID, &taskList.Title)
	if err != nil {
		return models.TaskList{}, fmt.Errorf("failed to find list: %w", err)
	}

	rows, _ := s.conn.Query(ctx,
		`SELECT tasks.id, tasks.text, tasks.points, tasks.burnt, tasks.state
		FROM tasks, task_list_map
		WHERE task_list_map.list_id = $1
			AND tasks.id = task_list_map.task_id
		ORDER BY tasks.id`,
		listID)
	defer rows.Close()

	err = rows.Err()
	for err == nil && rows.Next() {
		var task models.Task
		var taskID int64
		err = rows.Scan(&taskID, &task.Text, &task.Points, &task.Burnt, &task.State)
		if err == nil {
			task.ID = strconv.FormatInt(taskID, 16)
			taskList.Tasks = append(taskList.Tasks, task)
		}
	}

	if err != nil {
		return models.TaskList{}, err
	}

	return taskList, nil
}

func (s *TaskStore) UpdateTask(ctx context.Context, taskID string, opts models.UpdateOptions) error {
	id, err := strconv.ParseInt(taskID, 16, 64)
	if err != nil {
		return err
	}

	_, err = s.conn.Exec(ctx,
		"UPDATE tasks SET text = $2, points = $3, burnt = $4 WHERE id = $1",
		id, opts.Text, opts.Points, opts.Burnt)
	return err
}

func (s *TaskStore) DoneTask(ctx context.Context, taskID string) error {
	return s.updateTaskStateWithStmt(ctx, taskID, models.DoneTaskEvent,
		"UPDATE tasks SET state = $2, burnt=points WHERE id = $1")
}

func (s *TaskStore) UndoneTask(ctx context.Context, taskID string) error {
	return s.updateTaskState(ctx, taskID, models.UndoneTaskEvent)
}

func (s *TaskStore) TodoTask(ctx context.Context, taskID string) error {
	return s.updateTaskState(ctx, taskID, models.TodoTaskEvent)
}

func (s *TaskStore) CancelTask(ctx context.Context, taskID string) error {
	return s.updateTaskState(ctx, taskID, models.CancelTaskEvent)
}

func (s *TaskStore) updateTaskState(ctx context.Context, taskID string,
	event models.SwitchTaskStateEvent,
) error {
	return s.updateTaskStateWithStmt(ctx, taskID, event, "UPDATE tasks SET state = $2 WHERE id = $1")
}

func (s *TaskStore) updateTaskStateWithStmt(ctx context.Context, taskID string,
	event models.SwitchTaskStateEvent, stmt string,
) (resultErr error) {
	id, err := strconv.ParseInt(taskID, 16, 64)
	if err != nil {
		return err
	}

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

	row := tx.QueryRow(ctx, "SELECT state FROM tasks WHERE id = $1 FOR NO KEY UPDATE", id)
	var curState models.TaskState
	err = row.Scan(&curState)
	if err == nil || errors.Is(err, pgx.ErrNoRows) {
		state, err := curState.NextState(event)
		if err != nil {
			return err
		}

		if _, err := tx.Exec(ctx, stmt, id, state); err != nil {
			return fmt.Errorf("failed to execute update: %w", err)
		}
	} else if err != nil {
		return fmt.Errorf("failed to scan: %w", err)
	}

	return tx.Commit(ctx)
}

func (s *TaskStore) GetSprintTemplate(ctx context.Context) (models.SprintTemplate, error) {
	rows, _ := s.conn.Query(ctx,
		`SELECT id, text, points
		FROM new_sprint_task_tempate
		ORDER BY id`)
	defer rows.Close()

	var tasks []models.TaskTemplate
	err := rows.Err()
	for err == nil && rows.Next() {
		var task models.TaskTemplate
		var taskID int64
		err = rows.Scan(&taskID, &task.Text, &task.Points)
		if err == nil {
			task.ID = strconv.FormatInt(taskID, 16)
			tasks = append(tasks, task)
		}
	}

	if err != nil {
		return models.SprintTemplate{}, err
	}

	return models.SprintTemplate{Tasks: tasks}, nil
}
