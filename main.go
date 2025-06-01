package main

import (
	"os"
	"os/signal"
	"syscall"

	"github.com/foto-leistenschneider/admin-panel/internal/app"
)

func main() {
	app.Start()
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM, syscall.SIGHUP, syscall.SIGQUIT)
	<-c
	app.Stop()
}
