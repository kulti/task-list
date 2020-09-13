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

func (s *TaskStore) UpdateTask(
	ctx context.Context, taskID int64, fn storages.UpdateTaskFn,
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

	row := tx.QueryRow(ctx,
		`SELECT text, points, burnt, state
		 FROM tasks
		 WHERE id = $1
		 FOR NO KEY UPDATE`, taskID)
	var task storages.Task
	err = row.Scan(&task.Text, &task.Points, &task.Burnt, &task.State)
	if err != nil {
		if !errors.Is(err, pgx.ErrNoRows) {
			return fmt.Errorf("failed to scan: %w", err)
		}
	} else {
		task, err := fn(task)
		if err != nil {
			return err
		}

		_, err = tx.Exec(ctx,
			`UPDATE tasks
			 SET text = $2, points = $3, burnt = $4, state = $5
			 WHERE id = $1`,
			taskID, task.Text, task.Points, task.Burnt, task.State)
		if err != nil {
			return fmt.Errorf("failed to execute update: %w", err)
		}
	}

	return tx.Commit(ctx)
}

func (s *TaskStore) PostponeTask(
	ctx context.Context, taskID int64, fn storages.PostponeTaskFn,
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

	row := tx.QueryRow(ctx,
		`SELECT text, points, burnt, state
		 FROM tasks
		 WHERE id = $1`, taskID)

	var task storages.Task
	err = row.Scan(&task.Text, &task.Points, &task.Burnt, &task.State)
	if err != nil {
		if !errors.Is(err, pgx.ErrNoRows) {
			return fmt.Errorf("failed to scan: %w", err)
		}
		return tx.Commit(ctx)
	}

	postponedTask, updatedTask, err := fn(task)
	if err != nil {
		return err
	}

	if updatedTask.Points == 0 {
		_, err = tx.Exec(ctx, "DELETE FROM tasks WHERE id = $1", taskID)
		if err != nil {
			return fmt.Errorf("failed to execute delete: %w", err)
		}
	} else {
		_, err = tx.Exec(ctx,
			`UPDATE tasks
		 SET text = $2, points = $3, burnt = $4, state = $5
		 WHERE id = $1`,
			taskID, updatedTask.Text, updatedTask.Points, updatedTask.Burnt, updatedTask.State)
		if err != nil {
			return fmt.Errorf("failed to execute update: %w", err)
		}
	}

	_, err = tx.Exec(ctx,
		"INSERT INTO postponed_tasks (text, points) VALUES($1, $2)",
		postponedTask.Text, postponedTask.Points)
	if err != nil {
		return fmt.Errorf("failed to execute insert: %w", err)
	}

	return tx.Commit(ctx)
}

func (s *TaskStore) PopPostponedTasks(ctx context.Context) ([]models.PostponedTask, error) {
	var tasks []models.PostponedTask

	rows, _ := s.conn.Query(ctx, `DELETE FROM postponed_tasks RETURNING text, points`)
	err := rows.Err()
	for err == nil && rows.Next() {
		var task models.PostponedTask
		err = rows.Scan(&task.Text, &task.Points)
		if err == nil {
			tasks = append(tasks, task)
		}
	}
	rows.Close()

	if err != nil {
		return nil, err
	}

	return tasks, nil
}

func (s *TaskStore) GetSprintTemplate(ctx context.Context) (models.SprintTemplate, error) {
	var tasks []models.TaskTemplate

	rows, _ := s.conn.Query(ctx, `SELECT text, points
		FROM new_sprint_task_tempate
		ORDER BY id`)
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

	return models.SprintTemplate{Tasks: tasks}, nil
}

func (s *TaskStore) SetSprintTemplate(ctx context.Context, tmpl models.SprintTemplate) (resultErr error) {
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

	_, err = tx.Exec(ctx, `DELETE FROM new_sprint_task_tempate `)
	if err != nil {
		return fmt.Errorf("failed to cleanup new sprint template: %w", err)
	}

	for _, task := range tmpl.Tasks {
		_, err := tx.Exec(ctx, `INSERT INTO new_sprint_task_tempate (text, points) VALUES ($1, $2)`,
			task.Text, task.Points)
		if err != nil {
			return fmt.Errorf("failed to insert task %q with points %d: %w", task.Text, task.Points, err)
		}
	}

	return tx.Commit(ctx)
}
