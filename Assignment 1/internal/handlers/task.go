package handlers

import (
	"encoding/json"
	"net/http"
	"strconv"
)



type Task struct {
	ID    int    `json:"id"`
	Title string `json:"title"`
	Done  bool   `json:"done"`
}


var tasks = map[int]Task{
	1: {ID: 1, Title: "Write unit tests", Done: false},
	2: {ID: 2, Title: "Learn Go basics", Done: true},
}

var nextID = 3


func TasksHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	switch r.Method {

	case http.MethodDelete:
		deleteTask(w, r)


	case http.MethodGet:
		getTasks(w, r)

	case http.MethodPost:
		createTask(w, r)

	case http.MethodPatch:
		updateTask(w, r)

	default:
		w.WriteHeader(http.StatusMethodNotAllowed)
	}
}

func getTasks(w http.ResponseWriter, r *http.Request) {
	idParam := r.URL.Query().Get("id")

	if idParam != "" {
		id, err := strconv.Atoi(idParam)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			json.NewEncoder(w).Encode(map[string]string{
				"error": "invalid id",
			})
			return
		}

		task, ok := tasks[id]
		if !ok {
			w.WriteHeader(http.StatusNotFound)
			json.NewEncoder(w).Encode(map[string]string{
				"error": "task not found",
			})
			return
		}

		json.NewEncoder(w).Encode(task)
		return
	}

	list := []Task{}
	for _, t := range tasks {
		list = append(list, t)
	}

	json.NewEncoder(w).Encode(list)
}

func createTask(w http.ResponseWriter, r *http.Request) {
	var input struct {
		Title string `json:"title"`
	}

	if err := json.NewDecoder(r.Body).Decode(&input); err != nil || input.Title == "" {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{
			"error": "invalid title",
		})
		return
	}

	task := Task{
		ID:    nextID,
		Title: input.Title,
		Done:  false,
	}

	tasks[nextID] = task
	nextID++

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(task)
}

func updateTask(w http.ResponseWriter, r *http.Request) {
	idParam := r.URL.Query().Get("id")
	if idParam == "" {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{
			"error": "invalid id",
		})
		return
	}

	id, err := strconv.Atoi(idParam)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{
			"error": "invalid id",
		})
		return
	}

	task, ok := tasks[id]
	if !ok {
		w.WriteHeader(http.StatusNotFound)
		json.NewEncoder(w).Encode(map[string]string{
			"error": "task not found",
		})
		return
	}

	var input struct {
		Done *bool `json:"done"`
	}

	if err := json.NewDecoder(r.Body).Decode(&input); err != nil || input.Done == nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{
			"error": "invalid done value",
		})
		return
	}

	task.Done = *input.Done
	tasks[id] = task

	json.NewEncoder(w).Encode(map[string]bool{
		"updated": true,
	})
}

func deleteTask(w http.ResponseWriter, r *http.Request) {
	idParam := r.URL.Query().Get("id")

	if idParam == "" {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{
			"error": "invalid id",
		})
		return
	}

	id, err := strconv.Atoi(idParam)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{
			"error": "invalid id",
		})
		return
	}

	if _, ok := tasks[id]; !ok {
		w.WriteHeader(http.StatusNotFound)
		json.NewEncoder(w).Encode(map[string]string{
			"error": "task not found",
		})
		return
	}

	delete(tasks, id)

	json.NewEncoder(w).Encode(map[string]bool{
		"deleted": true,
	})
}
