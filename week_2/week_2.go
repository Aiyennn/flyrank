package main

import (
	"fmt"
	"net/http"
	"encoding/json"
)

type Task struct {
	ID int `json:"id"`
	Title string `json:"title"`
	Done bool `json:"done"`
}

var tasks = []Task{
	{
		ID: 1,
		Title: "Buy groceries",
		Done: false,
	},
	{
		ID: 2,
		Title: "Walk the dog",
		Done: true,
	},
	{
		ID: 3,
		Title: "Learn Go",
		Done: false,
	},
}

func apiDetails(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	w.WriteHeader(http.StatusOK)

	response := map[string]any{
		"name": "Task API",
		"version": "1.0",
		"endpoints": []string{"/tasks"},
	}

	json.NewEncoder(w).Encode(response)
}

func healthCheck(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	w.WriteHeader(http.StatusOK)

	response := map[string]any{
		"status": "ok",
	}

	json.NewEncoder(w).Encode(response)
}


func main() {
	http.HandleFunc("/", apiDetails)
	http.HandleFunc("/health", healthCheck)

	fmt.Println("Server running on http://localhost:8000")

	err := http.ListenAndServe(":8000", nil)
	if err != nil {
		fmt.Println("Sever failed:", err)
	}
}