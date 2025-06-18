package server

import (
	"io"
	"net/http"
	"strconv"

	"github.com/foto-leistenschneider/admin-panel/internal/db"
	"github.com/foto-leistenschneider/admin-panel/internal/tasks"
	"github.com/foto-leistenschneider/admin-panel/pkg/protos"
)

func tasksHandler(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()
	switch r.Method {
	case "POST":
		addTaskHandler(w, r)
	case "DELETE":
		deleteTaskHandler(w, r)
	default:
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
	}
}

func addTaskHandler(w http.ResponseWriter, r *http.Request) {
	if !isAuthenticated(r) {
		http.Error(w, "not authenticated", http.StatusUnauthorized)
		return
	}

	if err := r.ParseForm(); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	description := r.Form.Get("description")
	if description == "" {
		http.Error(w, "description is empty", http.StatusBadRequest)
		return
	}

	schedule := r.Form.Get("schedule")
	if schedule == "" {
		http.Error(w, "schedule is empty", http.StatusBadRequest)
		return
	}

	selector := r.Form.Get("selector")

	scope := r.Form.Get("scope")
	if scope == "" {
		http.Error(w, "scope is empty", http.StatusBadRequest)
		return
	}

	command := r.Form.Get("command")
	if command == "" && scope != protos.JobScope_Backup.String() {
		http.Error(w, "command is empty", http.StatusBadRequest)
		return
	}

	if t, err := db.Q.CreateTask(r.Context(), description, schedule, selector, command, scope); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	} else {
		tasks.Add(t)
	}
	http.Redirect(w, r, "/", http.StatusFound)
}

func deleteTaskHandler(w http.ResponseWriter, r *http.Request) {
	if !isAuthenticated(r) {
		http.Error(w, "not authenticated", http.StatusUnauthorized)
		return
	}

	var id int64
	if idBytes, err := io.ReadAll(r.Body); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	} else {
		id, err = strconv.ParseInt(string(idBytes), 10, 64)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
	}

	if err := db.Q.DeleteTask(r.Context(), id); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	tasks.Remove(id)
}
