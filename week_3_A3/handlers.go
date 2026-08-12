package main

import (
	"encoding/json"
	"errors"
	"net/http"
	"strconv"
	"strings"
)

func tasksHandler(service *TaskService) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodGet:
			tasks, err := service.List(r.Context())
			if err != nil {
				writeInternalError(w, err)
				return
			}
			writeJSON(w, http.StatusOK, tasks)

		case http.MethodPost:
			request, ok := decodeTaskRequest(w, r)
			if !ok {
				return
			}
			task, err := service.Create(r.Context(), request)
			if errors.Is(err, ErrInvalidTitle) {
				writeError(w, http.StatusBadRequest, err.Error())
				return
			}
			if err != nil {
				writeInternalError(w, err)
				return
			}
			writeJSON(w, http.StatusCreated, task)

		default:
			w.Header().Set("Allow", "GET, POST")
			writeError(w, http.StatusMethodNotAllowed, "method not allowed")
		}
	}
}

func taskHandler(service *TaskService) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		id, ok := taskID(w, r)
		if !ok {
			return
		}

		switch r.Method {
		case http.MethodGet:
			task, err := service.Get(r.Context(), id)
			writeTaskResult(w, task, err, http.StatusOK)

		case http.MethodPut:
			request, ok := decodeTaskRequest(w, r)
			if !ok {
				return
			}
			task, err := service.Update(r.Context(), id, request)
			writeTaskResult(w, task, err, http.StatusOK)

		case http.MethodDelete:
			err := service.Delete(r.Context(), id)
			if errors.Is(err, ErrTaskNotFound) {
				writeError(w, http.StatusNotFound, err.Error())
				return
			}
			if err != nil {
				writeInternalError(w, err)
				return
			}
			w.WriteHeader(http.StatusNoContent)

		default:
			w.Header().Set("Allow", "GET, PUT, DELETE")
			writeError(w, http.StatusMethodNotAllowed, "method not allowed")
		}
	}
}

func taskID(w http.ResponseWriter, r *http.Request) (int64, bool) {
	value := strings.TrimPrefix(r.URL.Path, "/tasks/")
	if value == "" || strings.Contains(value, "/") {
		writeError(w, http.StatusBadRequest, "invalid ID")
		return 0, false
	}

	id, err := strconv.ParseInt(value, 10, 64)
	if err != nil || id < 1 {
		writeError(w, http.StatusBadRequest, "invalid ID")
		return 0, false
	}
	return id, true
}

func decodeTaskRequest(w http.ResponseWriter, r *http.Request) (TaskRequest, bool) {
	defer r.Body.Close()
	var request TaskRequest
	decoder := json.NewDecoder(r.Body)
	decoder.DisallowUnknownFields()
	if err := decoder.Decode(&request); err != nil {
		writeError(w, http.StatusBadRequest, "invalid JSON")
		return TaskRequest{}, false
	}
	return request, true
}

func writeTaskResult(w http.ResponseWriter, task Task, err error, status int) {
	if errors.Is(err, ErrInvalidTitle) {
		writeError(w, http.StatusBadRequest, err.Error())
		return
	}
	if errors.Is(err, ErrTaskNotFound) {
		writeError(w, http.StatusNotFound, err.Error())
		return
	}
	if err != nil {
		writeInternalError(w, err)
		return
	}
	writeJSON(w, status, task)
}

func writeJSON(w http.ResponseWriter, status int, value any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(value)
}

func writeError(w http.ResponseWriter, status int, message string) {
	writeJSON(w, status, map[string]string{"error": message})
}

func writeInternalError(w http.ResponseWriter, err error) {
	writeError(w, http.StatusInternalServerError, err.Error())
}
