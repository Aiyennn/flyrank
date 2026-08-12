package main

import (
	"context"
	"sync"
)

type MemoryTaskRepository struct {
	mu     sync.Mutex
	nextID int64
	tasks  map[int64]Task
}

func NewMemoryTaskRepository() *MemoryTaskRepository {
	return &MemoryTaskRepository{nextID: 1, tasks: make(map[int64]Task)}
}

func (r *MemoryTaskRepository) List(context.Context) ([]Task, error) {
	r.mu.Lock()
	defer r.mu.Unlock()
	tasks := make([]Task, 0, len(r.tasks))
	for id := int64(1); id < r.nextID; id++ {
		if task, exists := r.tasks[id]; exists {
			tasks = append(tasks, task)
		}
	}
	return tasks, nil
}

func (r *MemoryTaskRepository) Get(_ context.Context, id int64) (Task, error) {
	r.mu.Lock()
	defer r.mu.Unlock()
	task, exists := r.tasks[id]
	if !exists {
		return Task{}, ErrTaskNotFound
	}
	return task, nil
}

func (r *MemoryTaskRepository) Create(_ context.Context, request TaskRequest) (Task, error) {
	r.mu.Lock()
	defer r.mu.Unlock()
	task := Task{ID: r.nextID, Title: request.Title, Done: request.Done}
	r.tasks[task.ID] = task
	r.nextID++
	return task, nil
}

func (r *MemoryTaskRepository) Update(_ context.Context, id int64, request TaskRequest) (Task, error) {
	r.mu.Lock()
	defer r.mu.Unlock()
	if _, exists := r.tasks[id]; !exists {
		return Task{}, ErrTaskNotFound
	}
	task := Task{ID: id, Title: request.Title, Done: request.Done}
	r.tasks[id] = task
	return task, nil
}

func (r *MemoryTaskRepository) Delete(_ context.Context, id int64) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	if _, exists := r.tasks[id]; !exists {
		return ErrTaskNotFound
	}
	delete(r.tasks, id)
	return nil
}
