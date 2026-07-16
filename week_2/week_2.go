package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"strings"
)

import (
    httpSwagger "github.com/swaggo/http-swagger"

    _ "week_2/docs"
)

// Task represents a single to-do item.
type Task struct {
	ID    int    `json:"id"`
	Title string `json:"title"`
	Done  bool   `json:"done"`
}

var tasks = []Task{
	{
		ID:    1,
		Title: "Buy groceries",
		Done:  false,
	},
	{
		ID:    2,
		Title: "Walk the dog",
		Done:  true,
	},
	{
		ID:    3,
		Title: "Learn Go",
		Done:  false,
	},
}

// apiDetails godoc
// @Summary      Get API details
// @Description  Returns basic metadata about the Task API, including its name, version, and available endpoints.
// @Tags         meta
// @Produce      json
// @Success      200  {object}  map[string]any
// @Router       / [get]
func apiDetails(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	w.WriteHeader(http.StatusOK)

	response := map[string]any{
		"name":      "Task API",
		"version":   "1.0",
		"endpoints": []string{"/tasks"},
	}

	json.NewEncoder(w).Encode(response)
}

// healthCheck godoc
// @Summary      Health check
// @Description  Returns the health status of the API.
// @Tags         meta
// @Produce      json
// @Success      200  {object}  map[string]string
// @Router       /health [get]
func healthCheck(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	w.WriteHeader(http.StatusOK)

	response := map[string]any{
		"status": "ok",
	}

	json.NewEncoder(w).Encode(response)
}

// getTask godoc
// @Summary      Get a task by ID
// @Description  Returns a single task matching the given ID.
// @Tags         tasks
// @Produce      json
// @Param        id   path      int  true  "Task ID"
// @Success      200  {object}  Task
// @Failure      400  {object}  map[string]string  "Invalid task ID"
// @Failure      404  {object}  map[string]string  "Task not found"
// @Router       /tasks/{id} [get]
func getTask(w http.ResponseWriter, r *http.Request) {
	idStr := r.PathValue("id")

	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "Invalid task ID", http.StatusBadRequest)
		return
	}

	for _, task := range tasks {
		if task.ID == id {
			w.Header().Set("Content-Type", "application/json")
			json.NewEncoder(w).Encode(task)
			return
		}
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusNotFound)

	json.NewEncoder(w).Encode(map[string]string{
		"error": "Task not found",
	})
}

// getTasks godoc
// @Summary      List all tasks
// @Description  Returns the full list of tasks.
// @Tags         tasks
// @Produce      json
// @Success      200  {array}  Task
// @Router       /tasks [get]
func getTasks(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(tasks)
}

// createTask godoc
// @Summary      Create a task
// @Description  Creates a new task with the given title.
// @Tags         tasks
// @Accept       json
// @Produce      json
// @Param        task  body      object{title=string}  true  "Task to create"
// @Success      201   {object}  Task
// @Failure      400   {object}  map[string]string  "Invalid JSON or missing title"
// @Router       /tasks [post]
func createTask(w http.ResponseWriter, r *http.Request) {
	var input struct {
		Title string `json:"title"`
	}

	err := json.NewDecoder(r.Body).Decode(&input)
	if err != nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)

		json.NewEncoder(w).Encode(map[string]string{
			"error": "Invalid JSON",
		})
		return
	}

	if strings.TrimSpace(input.Title) == "" {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)

		json.NewEncoder(w).Encode(map[string]string{
			"error": "title is required",
		})
		return
	}

	nextID := 1
	for _, task := range tasks {
		if task.ID >= nextID {
			nextID = task.ID + 1
		}
	}

	newTask := Task{
		ID:    nextID,
		Title: input.Title,
		Done:  false,
	}

	tasks = append(tasks, newTask)

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)

	json.NewEncoder(w).Encode(newTask)
}

// updateTask godoc
// @Summary      Update a task
// @Description  Updates the title and done status of an existing task.
// @Tags         tasks
// @Accept       json
// @Produce      json
// @Param        id    path      int   true  "Task ID"
// @Param        task  body      Task  true  "Updated task data"
// @Success      200   {object}  Task
// @Failure      400   {object}  map[string]string  "Invalid task ID or request body"
// @Failure      404   {object}  map[string]string  "Task not found"
// @Router       /tasks/{id} [put]
func updateTask(w http.ResponseWriter, r *http.Request) {
	idStr := r.PathValue("id")

	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "Invalid task ID", http.StatusBadRequest)
		return
	}

	var updated Task
	if err := json.NewDecoder(r.Body).Decode(&updated); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	for i := range tasks {
		if tasks[i].ID == id {
			tasks[i].Title = updated.Title
			tasks[i].Done = updated.Done

			w.Header().Set("Content-Type", "application/json")
			json.NewEncoder(w).Encode(tasks[i])
			return
		}
	}

	w.WriteHeader(http.StatusNotFound)
	json.NewEncoder(w).Encode(map[string]string{
		"error": "Task not found",
	})
}

// deleteTask godoc
// @Summary      Delete a task
// @Description  Deletes the task matching the given ID.
// @Tags         tasks
// @Produce      json
// @Param        id   path  int  true  "Task ID"
// @Success      204  "No Content"
// @Failure      400  {object}  map[string]string  "Invalid task ID"
// @Failure      404  {object}  map[string]string  "Task not found"
// @Router       /tasks/{id} [delete]
func deleteTask(w http.ResponseWriter, r *http.Request) {
	idStr := r.PathValue("id")

	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "Invalid task ID", http.StatusBadRequest)
		return
	}

	for i, task := range tasks {
		if task.ID == id {
			tasks = append(tasks[:i], tasks[i+1:]...)

			w.WriteHeader(http.StatusNoContent)
			return
		}
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusNotFound)
	json.NewEncoder(w).Encode(map[string]string{
		"error": "Task not found",
	})
}

// @title           Task API
// @version         1.0
// @description     A simple task management API built with Go's net/http.
// @termsOfService  http://swagger.io/terms/

// @contact.name   API Support
// @contact.email  support@example.com

// @license.name  MIT
// @license.url   https://opensource.org/licenses/MIT

// @host      localhost:8000
// @BasePath  /

func main() {
	http.HandleFunc("/", apiDetails)
	http.HandleFunc("/health", healthCheck)
	http.HandleFunc("GET /tasks", getTasks)
	http.HandleFunc("GET /tasks/{id}", getTask)
	http.HandleFunc("POST /tasks", createTask)
	http.HandleFunc("PUT /tasks/{id}", updateTask)
	http.HandleFunc("DELETE /tasks/{id}", deleteTask)
	
	http.Handle("/swagger/", httpSwagger.WrapHandler)

	fmt.Println("Server running on http://localhost:8000")


	err := http.ListenAndServe(":8000", nil)
	if err != nil {
		fmt.Println("Sever failed:", err)
	}
}