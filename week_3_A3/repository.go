package main

import (
	"context"
	"errors"
)

var ErrTaskNotFound = errors.New("task not found")

type TaskRepository interface {
	List(context.Context) ([]Task, error)
	Get(context.Context, int64) (Task, error)
	Create(context.Context, TaskRequest) (Task, error)
	Update(context.Context, int64, TaskRequest) (Task, error)
	Delete(context.Context, int64) error
}
