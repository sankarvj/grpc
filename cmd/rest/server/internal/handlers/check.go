package handlers

import (
	"context"
	"log"
	"net/http"

	"github.com/sankarvj/grpc/internal/platform/web"
)

var counter int64

type statusResponse struct {
	message string
	code    int
}

func status(ctx context.Context, w http.ResponseWriter, r *http.Request, params map[string]string) error {
	counter++
	log.Printf("Received: mcr for %d time", counter)
	statusRes := statusResponse{
		message: "Health is good for mcr",
		code:    200,
	}
	return web.Respond(ctx, w, statusRes, http.StatusOK)
}
