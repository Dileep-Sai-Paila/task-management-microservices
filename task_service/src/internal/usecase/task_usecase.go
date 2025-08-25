package usecase

import (
	"context"
	"fmt"
	"log"
	"task_service/internal/core/tasks"
)

type taskUsecase struct {
	taskRepo   TaskRepository
	userClient UserServiceClient
	cache      Cache
}

func NewTaskUsecase(repo TaskRepository, client UserServiceClient, cache Cache) TaskUsecase {
	return &taskUsecase{
		taskRepo:   repo,
		userClient: client,
		cache:      cache,
	}
}

// to checkk the cache before making a grpc call
func (uc *taskUsecase) CreateTask(ctx context.Context, task *tasks.Task) error {
	// checking if the user is already validated in the cache.
	isValidated, err := uc.cache.GetUserValidation(ctx, int32(task.UserID))
	if err != nil {
		log.Printf("Cache error: %v", err)
	}

	if isValidated {
		log.Printf("Cache HIT for user ID: %d", task.UserID)
	} else {
		log.Printf("Cache MISS for user ID: %d. Calling User Service.", task.UserID)
		// if it reacxhes here, it means it is not in cache, so we'll validate the user via grpc
		_, err := uc.userClient.GetUser(ctx, int32(task.UserID))
		if err != nil {
			return fmt.Errorf("invalid user ID: %w", err)
		}

		// if the user is valid, store the validation in the cache for next time.
		if err := uc.cache.SetUserValidation(ctx, int32(task.UserID)); err != nil {
			log.Printf("Could not set user validation in cache: %v", err)
		}
	}

	err = uc.taskRepo.CreateTask(ctx, task) //create task in database
	if err != nil {
		return fmt.Errorf("could not create task in repository: %w", err)
	}

	notificationMsg := fmt.Sprintf("Task '%s' created for user %d.", task.Title, task.UserID)
	if err := uc.cache.PublishTaskNotification(ctx, notificationMsg); err != nil {
		log.Printf("Failed to publish task creation notification: %v", err)
	}

	return nil
}

func (uc *taskUsecase) ListTasks(ctx context.Context, userID, status string) ([]tasks.Task, error) {
	taskList, err := uc.taskRepo.ListTasks(ctx, userID, status)
	if err != nil {
		return nil, fmt.Errorf("could not list tasks: %w", err)
	}
	return taskList, nil
}

func (uc *taskUsecase) UpdateTask(ctx context.Context, task *tasks.Task) (*tasks.Task, error) {
	updatedTask, err := uc.taskRepo.UpdateTask(ctx, task)
	if err != nil {
		return nil, fmt.Errorf("could not update task: %w", err)
	}

	// to publish the notification message after updating the task succesfully
	notificationMsg := fmt.Sprintf("Task %d updated. New status: %s.", updatedTask.ID, updatedTask.Status)
	if err := uc.cache.PublishTaskNotification(ctx, notificationMsg); err != nil {
		log.Printf("Failed to publish task update notification: %v", err)
	}
	return updatedTask, nil
}
