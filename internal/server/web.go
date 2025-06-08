package server

import (
	"context"
	"net/http"
	"time"

	"github.com/a-h/templ"
	"github.com/foto-leistenschneider/admin-panel/internal/server/view"
	"github.com/workos/workos-go/v4/pkg/usermanagement"
)

func registerRoutes() {
	http.HandleFunc("/", indexHandler(view.Intex()))
	http.HandleFunc("/runners/{runner}", templRenderHandler(view.RunnerJobs()))
	http.HandleFunc("/api/ping", runnerPingHandler)
	http.HandleFunc("/api/runners/{runner}/jobs", runnerJobsHandler)
	http.HandleFunc("/robots.txt", robotsTxtHandler)
	http.HandleFunc("/api/login", loginHandler)
	http.HandleFunc("/api/login_callback", loginCallbackHandler)
	http.HandleFunc("/api/logout", logoutHandler)
}

func templRenderHandler(component templ.Component) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := myContextFromRequest(r)
		_ = component.Render(ctx, w)
	}
}

func indexHandler(indexComponent templ.Component) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/" {
			ctx := myContextFromRequest(r)
			_ = indexComponent.Render(ctx, w)
		} else {
			view.EmbedFSHandler(w, r)
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

type myContext struct {
	ctx     context.Context
	user    *usermanagement.User
	request *http.Request
}

func myContextFromRequest(r *http.Request) context.Context {
	if user, ok := getAuthenticatedUser(r); ok {
		return myContext{
			ctx:     r.Context(),
			user:    &user,
			request: r,
		}
	} else {
		return myContext{
			ctx:     r.Context(),
			request: r,
		}
	}
}

func (c myContext) Deadline() (deadline time.Time, ok bool) {
	return c.ctx.Deadline()
}

func (c myContext) Done() <-chan struct{} {
	return c.ctx.Done()
}

func (c myContext) Err() error {
	return c.ctx.Err()
}

func (c myContext) Value(key any) any {
	if keyString, ok := key.(string); ok {
		if keyString == "user" {
			return c.user
		}
		if value := c.request.PathValue(keyString); value != "" {
			return value
		}
	}
	return c.ctx.Value(key)
}
