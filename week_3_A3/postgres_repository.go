package main

import (
	"context"
	"database/sql"
	"errors"
)

// PostgresTaskRepository is the production storage implementation.
type PostgresTaskRepository struct {
	db *sql.DB
}

func NewPostgresTaskRepository(db *sql.DB) *PostgresTaskRepository {
	return &PostgresTaskRepository{db: db}
}

func (r *PostgresTaskRepository) List(ctx context.Context) ([]Task, error) {
	rows, err := r.db.QueryContext(ctx, `SELECT id, title, done FROM tasks ORDER BY id`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	tasks := make([]Task, 0)
	for rows.Next() {
		var task Task
		if err := rows.Scan(&task.ID, &task.Title, &task.Done); err != nil {
			return nil, err
		}
		tasks = append(tasks, task)
	}

	return tasks, rows.Err()
}

func (r *PostgresTaskRepository) Get(ctx context.Context, id int64) (Task, error) {
	var task Task
	err := r.db.QueryRowContext(ctx, `SELECT id, title, done FROM tasks WHERE id = $1`, id).
		Scan(&task.ID, &task.Title, &task.Done)
	if errors.Is(err, sql.ErrNoRows) {
		return Task{}, ErrTaskNotFound
	}
	return task, err
}

func (r *PostgresTaskRepository) Create(ctx context.Context, request TaskRequest) (Task, error) {
	var task Task
	err := r.db.QueryRowContext(ctx,
		`INSERT INTO tasks (title, done) VALUES ($1, $2) RETURNING id, title, done`,
		request.Title, request.Done,
	).Scan(&task.ID, &task.Title, &task.Done)
	return task, err
}

func (r *PostgresTaskRepository) Update(ctx context.Context, id int64, request TaskRequest) (Task, error) {
	var task Task
	err := r.db.QueryRowContext(ctx,
		`UPDATE tasks SET title = $1, done = $2 WHERE id = $3 RETURNING id, title, done`,
		request.Title, request.Done, id,
	).Scan(&task.ID, &task.Title, &task.Done)
	if errors.Is(err, sql.ErrNoRows) {
		return Task{}, ErrTaskNotFound
	}
	return task, err
}

func (r *PostgresTaskRepository) Delete(ctx context.Context, id int64) error {
	result, err := r.db.ExecContext(ctx, `DELETE FROM tasks WHERE id = $1`, id)
	if err != nil {
		return err
	}

	count, err := result.RowsAffected()
	if err != nil {
		return err
	}
	if count == 0 {
		return ErrTaskNotFound
	}
	return nil
}
