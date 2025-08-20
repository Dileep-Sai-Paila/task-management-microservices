package handler

import (
	"encoding/json"
	"net/http"
	"strconv"
	"task_service/internal/core/tasks"
	"task_service/internal/usecase"

	"github.com/go-chi/chi/v5"
)

type TaskHandler struct {
	taskUsecase usecase.TaskUsecase
}

func NewTaskHandler(uc usecase.TaskUsecase) *TaskHandler {
	return &TaskHandler{
		taskUsecase: uc,
	}
}

type CreateTaskRequest struct {
	Title       string `json:"title"`
	Description string `json:"description"`
	UserID      int    `json:"user_id"`
}

// for POST /tasks endpoint
func (h *TaskHandler) CreateTask(w http.ResponseWriter, r *http.Request) {
	var req CreateTaskRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}
	task := &tasks.Task{
		Title:       req.Title,
		Description: req.Description,
		UserID:      req.UserID,
	}
	if err := h.taskUsecase.CreateTask(r.Context(), task); err != nil {
		if err.Error()[:17] == "invalid user ID:" {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(task)
}

// for GET /tasks endpoint
func (h *TaskHandler) ListTasks(w http.ResponseWriter, r *http.Request) {
	userID := r.URL.Query().Get("user_id")
	status := r.URL.Query().Get("status")
	tasks, err := h.taskUsecase.ListTasks(r.Context(), userID, status)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(tasks)
}

type UpdateTaskRequest struct {
	Title       string `json:"title,omitempty"`
	Description string `json:"description,omitempty"`
	Status      string `json:"status,omitempty"`
}

// for PUT /tasks/{id}.
func (h *TaskHandler) UpdateTask(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "Invalid task ID", http.StatusBadRequest)
		return
	}

	var req UpdateTaskRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	task := &tasks.Task{
		ID:          id,
		Title:       req.Title,
		Description: req.Description,
		Status:      req.Status,
	}

	updatedTask, err := h.taskUsecase.UpdateTask(r.Context(), task)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(updatedTask)
}
