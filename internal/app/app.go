package app

import (
	"github.com/foto-leistenschneider/admin-panel/internal/db"
	"github.com/foto-leistenschneider/admin-panel/internal/server"
	"github.com/foto-leistenschneider/admin-panel/internal/tasks"
)

func Start() {
	server.Start()
}

func Stop() {
	server.Stop()
	_ = db.Q.Close()
	tasks.Close()
}
