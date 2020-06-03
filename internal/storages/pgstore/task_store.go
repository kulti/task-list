package pgstore

import (
	"context"
	"errors"
	"fmt"
	"strconv"
	"time"

	"github.com/jackc/pgx/v4/pgxpool"
	"github.com/kulti/task-list/internal/models"
)

type TaskStore struct {
	conn       *pgxpool.Pool
	todoListID int64
}

func New(dbURL string) (*TaskStore, error) {
	conn, err := pgxpool.Connect(context.Background(), dbURL)
	if err != nil {
		return nil, fmt.Errorf("unable to connect to database: %v", err)
	}

	s := &TaskStore{
		conn: conn,
	}

	row := conn.QueryRow(context.Background(), "SELECT id FROM task_lists WHERE type = $1", "todo")
	err = row.Scan(&s.todoListID)
	if err != nil {
		return nil, fmt.Errorf("unable to find todo list: %v", err)
	}

	return s, nil
}

func (s *TaskStore) Close() {
	if s.conn != nil {
		s.conn.Close()
	}
}

func (s *TaskStore) NewSprint(ctx context.Context, title string) error {
	_, err := s.conn.Exec(ctx,
		"INSERT INTO task_lists (type, title, created_at) VALUES ($1, $2, $3)",
		"sprint", title, time.Now())
	if err != nil {
		return err
	}

	_, err = s.conn.Exec(ctx, "DELETE FROM task_list_map WHERE list_id = $1", s.todoListID)
	return err
}

func (s *TaskStore) CreateTask(ctx context.Context, task models.Task, listIDs []string,
) (string, error) {
	row := s.conn.QueryRow(ctx,
		"INSERT INTO tasks (text, points, burnt, state) VALUES($1, $2, $3, $4) RETURNING id",
		task.Text, task.Points, task.Burnt, task.State)
	var taskID int64
	err := row.Scan(&taskID)
	if err != nil {
		return "", fmt.Errorf("failed to create task: %v", err)
	}

	for _, listType := range listIDs {
		row = s.conn.QueryRow(ctx,
			"SELECT id FROM task_lists WHERE type = $1 ORDER BY created_at DESC LIMIT 1",
			listType)

		var listID int
		err = row.Scan(&listID)
		if err != nil {
			return "", fmt.Errorf("failed to find list: %v", err)
		}

		_, err = s.conn.Exec(ctx,
			"INSERT INTO task_list_map (task_id, list_id) VALUES ($1, $2)",
			taskID, listID)
		if err != nil {
			return "", fmt.Errorf("failed to add task to list: %v", err)
		}
	}

	return strconv.FormatInt(taskID, 16), nil
}

func (s *TaskStore) TakeTaskToList(ctx context.Context, taskID, listIDs string) error {
	id, err := strconv.ParseInt(taskID, 16, 8)
	if err != nil {
		return err
	}

	_, err = s.conn.Exec(ctx,
		"INSERT INTO task_list_map (task_id, list_id) VALUES ($1, $2)",
		id, s.todoListID)
	if err != nil {
		return err
	}

	return s.updateTaskState(ctx, taskID, "todo")
}

func (s *TaskStore) DeleteTaskFromList(ctx context.Context, taskID, listType string) error {
	id, err := strconv.ParseInt(taskID, 16, 8)
	if err != nil {
		return err
	}
	switch listType {
	case "sprint":
		_, err = s.conn.Exec(ctx, "DELETE FROM tasks WHERE id = $1", id)
	case "todo":
		_, err = s.conn.Exec(ctx,
			"DELETE FROM task_list_map WHERE task_id = $1 AND list_id = $2",
			id, s.todoListID)
	default:
		err = errors.New("unknown list type")
	}
	return err
}

func (s *TaskStore) ListTasks(ctx context.Context, listType string) (models.TaskList, error) {
	var taskList models.TaskList
	var listID int64

	switch listType {
	case "todo":
		listID = s.todoListID
		taskList.Title = "Todo"
	default:
		row := s.conn.QueryRow(ctx,
			"SELECT id, title FROM task_lists WHERE type = $1 ORDER BY created_at DESC LIMIT 1",
			listType)

		err := row.Scan(&listID, &taskList.Title)
		if err != nil {
			return models.TaskList{}, fmt.Errorf("failed to find list: %v", err)
		}
	}

	rows, _ := s.conn.Query(ctx,
		`SELECT tasks.id, tasks.text, tasks.points, tasks.burnt, tasks.state
		FROM tasks, task_list_map
		WHERE task_list_map.list_id = $1
			AND tasks.id = task_list_map.task_id
		ORDER BY tasks.id`,
		listID)
	defer rows.Close()

	err := rows.Err()
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
	id, err := strconv.ParseInt(taskID, 16, 8)
	if err != nil {
		return err
	}

	_, err = s.conn.Exec(ctx,
		"UPDATE tasks SET text = $2, points = $3, burnt = $4 WHERE id = $1",
		id, opts.Text, opts.Points, opts.Burnt)
	return err
}

func (s *TaskStore) DoneTask(ctx context.Context, taskID string) error {
	return s.updateTaskState(ctx, taskID, "done")
}

func (s *TaskStore) CancelTask(ctx context.Context, taskID string) error {
	return s.updateTaskState(ctx, taskID, "canceled")
}

func (s *TaskStore) updateTaskState(ctx context.Context, taskID, state string) error {
	id, err := strconv.ParseInt(taskID, 16, 8)
	if err != nil {
		return err
	}
	_, err = s.conn.Exec(ctx, "UPDATE tasks SET state = $2 WHERE id = $1", id, state)
	return err
}
