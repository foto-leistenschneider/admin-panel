package server

import (
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
		w.Header().Set("Content-Type", "text/plain")
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	ping := protos.Ping{}
	if err := proto.Unmarshal(body, &ping); err != nil {
		w.Header().Set("Content-Type", "text/plain")
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	if ping.Name == "" {
		w.Header().Set("Content-Type", "text/plain")
		http.Error(w, "ping name is empty", http.StatusBadRequest)
		return
	}

	newJobs, err := runners.Ping(&ping)
	if err != nil {
		w.Header().Set("Content-Type", "text/plain")
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	newJobsBytes, err := proto.Marshal(newJobs)
	if err != nil {
		w.Header().Set("Content-Type", "text/plain")
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/x-protobuf")
	_, _ = w.Write(newJobsBytes)
}
