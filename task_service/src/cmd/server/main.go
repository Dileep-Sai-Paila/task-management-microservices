package main

import (
	"log"
	"net/http"
	"task_service/internal/adaptors/grpcclient"
	"task_service/internal/adaptors/persistance"
	"task_service/internal/adaptors/redis"
	"task_service/internal/config"
	"task_service/internal/interfaces/input/api/rest/handler"
	"task_service/internal/usecase"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

func main() {
	cfg, err := config.LoadConfig(".")
	if err != nil {
		log.Fatalf("could not load config: %v", err)
	}

	dbStore, err := persistance.NewDBStore(cfg.DBSource)
	if err != nil {
		log.Fatalf("could not connect to database: %v", err)
	}

	userClient, err := grpcclient.NewUserClient(cfg.UserServiceGRPCAddress)
	if err != nil {
		log.Fatalf("could not create user service client: %v", err)
	}

	redisCache, err := redis.NewCache(cfg.RedisAddress)
	if err != nil {
		log.Fatalf("could not connect to redis: %v", err)
	}

	taskUsecase := usecase.NewTaskUsecase(dbStore, userClient, redisCache)

	taskHandler := handler.NewTaskHandler(taskUsecase)

	r := chi.NewRouter()
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)

	r.Post("/tasks", taskHandler.CreateTask)
	r.Get("/tasks", taskHandler.ListTasks)
	r.Put("/tasks/{id}", taskHandler.UpdateTask)

	log.Printf("Task Service starting on %s", cfg.ServerAddress)
	if err := http.ListenAndServe(cfg.ServerAddress, r); err != nil {
		log.Fatalf("failed to start server: %v", err)
	}
}
