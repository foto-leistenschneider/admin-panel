package server

import (
	"fmt"
	"io"
	"net/http"

	"github.com/foto-leistenschneider/admin-panel/internal/runners"
	"github.com/foto-leistenschneider/admin-panel/pkg/protos"
	"google.golang.org/protobuf/proto"
)

func runnerPingHandler(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()
	body, err := io.ReadAll(r.Body)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	ping := protos.Ping{}
	if err := proto.Unmarshal(body, &ping); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	if ping.Name == "" {
		http.Error(w, "ping name is empty", http.StatusBadRequest)
		return
	}

	newJobs, err := runners.Ping(&ping)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	newJobsBytes, err := proto.Marshal(newJobs)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/octet-stream")
	_, _ = w.Write(newJobsBytes)
}

func runnerJobsHandler(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()
	switch r.Method {
	case "POST":
		addJobHandler(w, r)
	default:
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
	}
}

func addJobHandler(w http.ResponseWriter, r *http.Request) {
	if !isAuthenticated(r) {
		http.Error(w, "not authenticated", http.StatusUnauthorized)
		return
	}

	runnerName := r.PathValue("runner")
	if runnerName == "" {
		http.Error(w, "runner name is empty", http.StatusBadRequest)
		return
	}

	runner, ok := runners.Register[runnerName]
	if !ok {
		http.Error(w, "runner not found", http.StatusNotFound)
		return
	}

	if err := r.ParseForm(); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	if err := runner.AddJob(r.Form.Get("scope"), r.Form.Get("command")); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	http.Redirect(w, r, fmt.Sprintf("/runners/%s", runnerName), http.StatusFound)
}
