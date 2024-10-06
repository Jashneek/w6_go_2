package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"github.com/gorilla/mux"
)

type Task struct {
	ID          int    `json:"id"`
	Title       string `json:"title"`
	Description string `json:"description"`
	Status      string `json:"status"` // "pending" or "completed"
}

var tasks []Task
var idCounter = 1

func main() {
	r := mux.NewRouter()

	// Set up the route handlers
	r.HandleFunc("/tasks", CreateTask).Methods("POST")
	r.HandleFunc("/tasks", GetAllTasks).Methods("GET")
	r.HandleFunc("/tasks/{id}", GetTaskByID).Methods("GET")
	r.HandleFunc("/tasks/{id}", UpdateTask).Methods("PUT")
	r.HandleFunc("/tasks/{id}", DeleteTask).Methods("DELETE")

	fmt.Println("Server is running on http://localhost:8080")
	http.ListenAndServe(":8080", r)
}

// CreateTask creates a new task
func CreateTask(w http.ResponseWriter, r *http.Request) {
	var task Task
	json.NewDecoder(r.Body).Decode(&task)

	task.ID = idCounter
	idCounter++
	task.Status = "pending"

	tasks = append(tasks, task)
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(task)
}

// GetAllTasks returns all tasks
func GetAllTasks(w http.ResponseWriter, r *http.Request) {
	json.NewEncoder(w).Encode(tasks)
}

// GetTaskByID returns a task by its ID
func GetTaskByID(w http.ResponseWriter, r *http.Request) {
	id, _ := strconv.Atoi(mux.Vars(r)["id"])
	for _, task := range tasks {
		if task.ID == id {
			json.NewEncoder(w).Encode(task)
			return
		}
	}
	http.Error(w, "Task not found", http.StatusNotFound)
}

// UpdateTask updates a task by its ID
func UpdateTask(w http.ResponseWriter, r *http.Request) {
	id, _ := strconv.Atoi(mux.Vars(r)["id"])

	for i, task := range tasks {
		if task.ID == id {
			json.NewDecoder(r.Body).Decode(&tasks[i])
			json.NewEncoder(w).Encode(tasks[i])
			return
		}
	}
	http.Error(w, "Task not found", http.StatusNotFound)
}

// DeleteTask deletes a task by its ID
func DeleteTask(w http.ResponseWriter, r *http.Request) {
	id, _ := strconv.Atoi(mux.Vars(r)["id"])

	for i, task := range tasks {
		if task.ID == id {
			tasks = append(tasks[:i], tasks[i+1:]...)
			w.WriteHeader(http.StatusNoContent)
			return
		}
	}
	http.Error(w, "Task not found", http.StatusNotFound)
}
