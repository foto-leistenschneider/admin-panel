package server

import (
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strings"

	"github.com/foto-leistenschneider/admin-panel/internal/config"
	"github.com/foto-leistenschneider/admin-panel/internal/runners"
	"github.com/foto-leistenschneider/admin-panel/pkg/protos"
)

func backupHandler(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()
	fileName := filepath.Join(config.BackupDir, r.PathValue("filename"))

	if runner, ok := runners.Register[r.PathValue("runner")]; !ok {
		http.Error(w, "runner not found", http.StatusNotFound)
		return
	} else {
		if job, ok := runner.Jobs[jobIdFromFileName(fileName)]; !ok {
			http.Error(w, "job not found", http.StatusNotFound)
			return
		} else if job.Status >= protos.JobStatus_Done {
			http.Error(w, "job not running", http.StatusAlreadyReported)
			return
		} else if job.Scope != protos.JobScope_Backup {
			http.Error(w, "job is not a backup job", http.StatusConflict)
			return
		}
	}

	f, err := os.Create(fileName)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer f.Close()

	_, err = io.Copy(f, r.Body)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)
}

func jobIdFromFileName(fileName string) string {
	return strings.Split(filepath.Base(fileName), ".")[0]
}
