package main

import (
	"context"
	"errors"
	"strings"
)

var ErrInvalidTitle = errors.New("title is required")

type TaskService struct {
	repository TaskRepository
}

func NewTaskService(repository TaskRepository) *TaskService {
	return &TaskService{repository: repository}
}

func (s *TaskService) List(ctx context.Context) ([]Task, error) {
	return s.repository.List(ctx)
}

func (s *TaskService) Get(ctx context.Context, id int64) (Task, error) {
	return s.repository.Get(ctx, id)
}

func (s *TaskService) Create(ctx context.Context, request TaskRequest) (Task, error) {
	request.Title = strings.TrimSpace(request.Title)
	if request.Title == "" {
		return Task{}, ErrInvalidTitle
	}

	return s.repository.Create(ctx, request)
}

func (s *TaskService) Update(ctx context.Context, id int64, request TaskRequest) (Task, error) {
	request.Title = strings.TrimSpace(request.Title)
	if request.Title == "" {
		return Task{}, ErrInvalidTitle
	}

	return s.repository.Update(ctx, id, request)
}

func (s *TaskService) Delete(ctx context.Context, id int64) error {
	return s.repository.Delete(ctx, id)
}
