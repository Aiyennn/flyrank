package main

import (
	"database/sql"
	"encoding/json"
	"log"
	"net/http"
	"strconv"
	"strings"

	_ "modernc.org/sqlite"
)

type Task struct {
	ID    int    `json:"id"`
	Title string `json:"title"`
	Done  bool   `json:"done"`
}

type TaskRequest struct {
	Title string `json:"title"`
	Done  bool   `json:"done"`
}

var db *sql.DB

func main() {
	var err error

	db, err = sql.Open("sqlite", "tasks.db")
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	initDatabase()

	http.HandleFunc("/tasks", tasksHandler)
	http.HandleFunc("/tasks/", taskHandler)

	log.Println("Server running on http://localhost:8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}

func initDatabase() {
	_, err := db.Exec(`
	CREATE TABLE IF NOT EXISTS tasks(
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		title TEXT NOT NULL,
		done BOOLEAN NOT NULL
	)
	`)
	if err != nil {
		log.Fatal(err)
	}

	var count int
	err = db.QueryRow("SELECT COUNT(*) FROM tasks").Scan(&count)
	if err != nil {
		log.Fatal(err)
	}

	if count == 0 {
		db.Exec("INSERT INTO tasks(title, done) VALUES(?, ?)", "Buy groceries", false)
		db.Exec("INSERT INTO tasks(title, done) VALUES(?, ?)", "Study Go", false)
		db.Exec("INSERT INTO tasks(title, done) VALUES(?, ?)", "Exercise", true)
	}
}

func tasksHandler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {

	case http.MethodGet:
		rows, err := db.Query("SELECT id, title, done FROM tasks")
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		defer rows.Close()

		var tasks []Task

		for rows.Next() {
			var t Task
			rows.Scan(&t.ID, &t.Title, &t.Done)
			tasks = append(tasks, t)
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(tasks)

	case http.MethodPost:
		var req TaskRequest

		err := json.NewDecoder(r.Body).Decode(&req)
		if err != nil {
			http.Error(w, `{"error":"Invalid JSON"}`, http.StatusBadRequest)
			return
		}

		if strings.TrimSpace(req.Title) == "" {
			http.Error(w, `{"error":"Title is required"}`, http.StatusBadRequest)
			return
		}

		result, err := db.Exec(
			"INSERT INTO tasks(title, done) VALUES(?, ?)",
			req.Title,
			req.Done,
		)

		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		id, _ := result.LastInsertId()

		task := Task{
			ID:    int(id),
			Title: req.Title,
			Done:  req.Done,
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusCreated)
		json.NewEncoder(w).Encode(task)

	default:
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
	}
}

func taskHandler(w http.ResponseWriter, r *http.Request) {
	idStr := strings.TrimPrefix(r.URL.Path, "/tasks/")
	id, err := strconv.Atoi(idStr)

	if err != nil {
		http.Error(w, `{"error":"Invalid ID"}`, http.StatusBadRequest)
		return
	}

	switch r.Method {

	case http.MethodGet:
		var task Task

		err := db.QueryRow(
			"SELECT id, title, done FROM tasks WHERE id=?",
			id,
		).Scan(&task.ID, &task.Title, &task.Done)

		if err == sql.ErrNoRows {
			http.Error(w, `{"error":"Task not found"}`, http.StatusNotFound)
			return
		}

		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(task)

	case http.MethodPut:
		var req TaskRequest

		err := json.NewDecoder(r.Body).Decode(&req)
		if err != nil {
			http.Error(w, `{"error":"Invalid JSON"}`, http.StatusBadRequest)
			return
		}

		if strings.TrimSpace(req.Title) == "" {
			http.Error(w, `{"error":"Title is required"}`, http.StatusBadRequest)
			return
		}

		result, err := db.Exec(
			"UPDATE tasks SET title=?, done=? WHERE id=?",
			req.Title,
			req.Done,
			id,
		)

		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		rows, _ := result.RowsAffected()

		if rows == 0 {
			http.Error(w, `{"error":"Task not found"}`, http.StatusNotFound)
			return
		}

		task := Task{
			ID:    id,
			Title: req.Title,
			Done:  req.Done,
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(task)

	case http.MethodDelete:
		result, err := db.Exec(
			"DELETE FROM tasks WHERE id=?",
			id,
		)

		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		rows, _ := result.RowsAffected()

		if rows == 0 {
			http.Error(w, `{"error":"Task not found"}`, http.StatusNotFound)
			return
		}

		w.WriteHeader(http.StatusNoContent)

	default:
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
	}
}