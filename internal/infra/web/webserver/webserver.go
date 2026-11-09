package webserver

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

type WebServer struct {
	Router        chi.Router
	Handlers      map[string]http.HandlerFunc
	WebServerPort string
}

func NewWebServer(serverPort string) *WebServer {
	return &WebServer{
		Router:        chi.NewRouter(),
		Handlers:      make(map[string]http.HandlerFunc),
		WebServerPort: serverPort,
	}
}

func (s *WebServer) AddHandler(path string, handler http.HandlerFunc) {
	s.Handlers[path] = handler
}

// Handler builds the routing table before the HTTP server starts. The caller
// owns the listener, timeouts and shutdown.
func (s *WebServer) Handler() http.Handler {
	r := chi.NewRouter()
	r.Use(middleware.Logger)
	for path, handler := range s.Handlers {
		r.Handle(path, handler)
	}
	return r
}

// Start remains available for the standalone exercise; callers must handle
// its error instead of silently discarding failed binds.
func (s *WebServer) Start() error {
	return http.ListenAndServe(s.WebServerPort, s.Handler())
}
