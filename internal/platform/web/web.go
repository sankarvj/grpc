package web

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/dimfeld/httptreemux/v5"
	"go.opencensus.io/plugin/ochttp"
	"go.opencensus.io/plugin/ochttp/propagation/tracecontext"
	"go.opencensus.io/trace"
)

// ctxKey represents the type of value for the context key.
type ctxKey int

// KeyValues is how request values or stored/retrieved.
const KeyValues ctxKey = 1

// Values represent state for each request.
type Values struct {
	TraceID    string
	Now        time.Time
	StatusCode int
}

//Handler for our own framework
type Handler func(ctx context.Context, w http.ResponseWriter, r *http.Request, params map[string]string) error

//App is the entry point
type App struct {
	*httptreemux.TreeMux
	och *ochttp.Handler
}

// NewApp creates an App value that handle a set of routes for the application.
func NewApp() *App {
	app := App{
		TreeMux: httptreemux.New(),
	}

	// Create an OpenCensus HTTP Handler which wraps the router. This will start
	// the initial span and annotate it with information about the request/response.
	//
	// This is configured to use the W3C TraceContext standard to set the remote
	// parent if an client request includes the appropriate headers.
	// https://w3c.github.io/trace-context/
	app.och = &ochttp.Handler{
		Handler:     app.TreeMux,
		Propagation: &tracecontext.HTTPFormat{},
	}

	return &app
}

// Handle is our mechanism for mounting Handlers for a given HTTP verb and path
// pair, this makes for really easy, convenient routing.
func (a *App) Handle(verb, path string, handler Handler) {

	// The function to execute for each request.
	h := func(w http.ResponseWriter, r *http.Request, params map[string]string) {
		ctx, span := trace.StartSpan(r.Context(), "internal.platform.web")
		defer span.End()

		// Set the context with the required values to
		// process the request.
		v := Values{
			TraceID: span.SpanContext().TraceID.String(),
			Now:     time.Now(),
		}
		ctx = context.WithValue(ctx, KeyValues, &v)

		// Call the wrapped handler functions.
		if err := handler(ctx, w, r, params); err != nil {
			fmt.Println("Error while serving http ", err)
			return
		}
	}

	// Add this handler for the specified verb and route.
	a.TreeMux.Handle(verb, path, h)
}

// ServeHTTP implements the http.Handler interface. It overrides the ServeHTTP
// of the embedded TreeMux by using the ochttp.Handler instead. That Handler
// wraps the TreeMux handler so the routes are served.
func (a *App) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	a.och.ServeHTTP(w, r)
}
