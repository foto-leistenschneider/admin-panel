package server

import (
	"net/http"

	"github.com/a-h/templ"
	"github.com/foto-leistenschneider/admin-panel/internal/server/view"
)

func registerRoutes() {
	http.HandleFunc("/", indexHandler(view.Intex()))
	http.HandleFunc("/runners/{runner}", templRenderHandler(view.RunnerJobs()))
	http.HandleFunc("/runners/ping", runnerPingHandler)
	http.HandleFunc("/robots.txt", robotsTxtHandler)
}

func templRenderHandler(component templ.Component) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		_ = component.Render(r.Context(), w)
	}
}

func indexHandler(indexComponent templ.Component) http.HandlerFunc {
	viewFsHandler := http.FileServer(http.FS(view.FS))
	return func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/" {
			_ = indexComponent.Render(r.Context(), w)
		} else {
			viewFsHandler.ServeHTTP(w, r)
		}
	}
}

// Disallow all crawlers on all pages
var robotsTxt = []byte(`User-agent: *
Disallow: /
`)

func robotsTxtHandler(w http.ResponseWriter, _ *http.Request) {
	w.Header().Set("Content-Type", "text/plain")
	w.WriteHeader(http.StatusOK)
	_, _ = w.Write(robotsTxt)
}
