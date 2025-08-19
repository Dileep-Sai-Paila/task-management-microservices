package usecase

import (
	"context"
	"task_service/internal/core/tasks"
	pb "task_service/proto"
)

// business logic for tasks
type TaskUsecase interface {
	CreateTask(ctx context.Context, task *tasks.Task) error
	ListTasks(ctx context.Context, userID, status string) ([]tasks.Task, error)
	UpdateTask(ctx context.Context, task *tasks.Task) (*tasks.Task, error)
}

// persistence operations for tasks
type TaskRepository interface {
	CreateTask(ctx context.Context, task *tasks.Task) error
	ListTasks(ctx context.Context, userID, status string) ([]tasks.Task, error)
	UpdateTask(ctx context.Context, task *tasks.Task) (*tasks.Task, error)
}

// for communicating with the User Service
type UserServiceClient interface {
	GetUser(ctx context.Context, userID int32) (*pb.GetUserResponse, error)
}

// for the caching layer
type Cache interface {
	SetUserValidation(ctx context.Context, userID int32) error
	GetUserValidation(ctx context.Context, userID int32) (bool, error)
}
