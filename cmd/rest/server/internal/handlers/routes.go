package handlers

import (
	"net/http"

	"github.com/sankarvj/grpc/internal/platform/web"
)

// API constructs an http.Handler with all application routes defined.
func API() http.Handler {
	app := web.NewApp()
	app.Handle("GET", "/check", status)
	return app
}
