package persistance

import (
	"context"
	"database/sql"
	"fmt"
	"io/ioutil"
	"log"
	"strings"
	"task_service/internal/core/tasks"
	"time"

	_ "github.com/lib/pq"
)

type DBStore struct {
	DB *sql.DB
}

func NewDBStore(dbSource string) (*DBStore, error) {
	db, err := sql.Open("postgres", dbSource)
	if err != nil {
		return nil, fmt.Errorf("could not open database: %w", err)
	}
	if err := db.Ping(); err != nil {
		return nil, fmt.Errorf("could not ping database: %w", err)
	}
	log.Println("Task service database connection successful.")
	if err := runMigrations(db); err != nil {
		return nil, fmt.Errorf("could not run migrations: %w", err)
	}
	return &DBStore{DB: db}, nil
}

// reads and executes the sql file
func runMigrations(db *sql.DB) error {
	migration, err := ioutil.ReadFile("migrations/init.sql")
	if err != nil {
		return fmt.Errorf("could not read migration file: %w", err)
	}
	_, err = db.Exec(string(migration))
	if err != nil {
		return fmt.Errorf("could not execute migration: %w", err)
	}
	log.Println("Task service database migration successful.")
	return nil
}

func (store *DBStore) CreateTask(ctx context.Context, task *tasks.Task) error {
	query := `INSERT INTO tasks (title, description, user_id) VALUES ($1, $2, $3) RETURNING id, status, created_at, updated_at`
	err := store.DB.QueryRowContext(ctx, query, task.Title, task.Description, task.UserID).Scan(
		&task.ID,
		&task.Status,
		&task.CreatedAt,
		&task.UpdatedAt,
	)
	if err != nil {
		return fmt.Errorf("could not create task: %w", err)
	}
	return nil
}

// retrieves a list of tasks
func (store *DBStore) ListTasks(ctx context.Context, userID, status string) ([]tasks.Task, error) {
	query := "SELECT id, title, description, status, user_id, created_at, updated_at FROM tasks"
	var conditions []string
	var args []interface{}
	argID := 1
	if userID != "" {
		conditions = append(conditions, fmt.Sprintf("user_id = $%d", argID))
		args = append(args, userID)
		argID++
	}
	if status != "" {
		conditions = append(conditions, fmt.Sprintf("status = $%d", argID))
		args = append(args, status)
		argID++
	}
	if len(conditions) > 0 {
		query += " WHERE " + strings.Join(conditions, " AND ")
	}
	rows, err := store.DB.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, fmt.Errorf("could not query tasks: %w", err)
	}
	defer rows.Close()
	var taskList []tasks.Task
	for rows.Next() {
		var task tasks.Task
		if err := rows.Scan(&task.ID, &task.Title, &task.Description, &task.Status, &task.UserID, &task.CreatedAt, &task.UpdatedAt); err != nil {
			return nil, fmt.Errorf("could not scan task row: %w", err)
		}
		taskList = append(taskList, task)
	}
	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating task rows: %w", err)
	}
	return taskList, nil
}

func (store *DBStore) UpdateTask(ctx context.Context, task *tasks.Task) (*tasks.Task, error) {
	var setClauses []string
	var args []interface{}
	argID := 1

	// Check each field and add it to the SET clause if it's not a zero value.
	if task.Title != "" {
		setClauses = append(setClauses, fmt.Sprintf("title = $%d", argID))
		args = append(args, task.Title)
		argID++
	}
	if task.Description != "" {
		setClauses = append(setClauses, fmt.Sprintf("description = $%d", argID))
		args = append(args, task.Description)
		argID++
	}
	if task.Status != "" {
		setClauses = append(setClauses, fmt.Sprintf("status = $%d", argID))
		args = append(args, task.Status)
		argID++
	}

	// if no fields aree provided to update
	if len(setClauses) == 0 {
		return nil, fmt.Errorf("no fields to update")
	}

	setClauses = append(setClauses, fmt.Sprintf("updated_at = $%d", argID))
	args = append(args, time.Now())
	argID++

	args = append(args, task.ID)

	query := fmt.Sprintf("UPDATE tasks SET %s WHERE id = $%d RETURNING id, title, description, status, user_id, created_at, updated_at",
		strings.Join(setClauses, ", "), argID)

	updatedTask := &tasks.Task{}
	err := store.DB.QueryRowContext(ctx, query, args...).Scan(
		&updatedTask.ID,
		&updatedTask.Title,
		&updatedTask.Description,
		&updatedTask.Status,
		&updatedTask.UserID,
		&updatedTask.CreatedAt,
		&updatedTask.UpdatedAt,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("task with id %d not found", task.ID)
		}
		return nil, fmt.Errorf("could not update task: %w", err)
	}
	return updatedTask, nil
}
