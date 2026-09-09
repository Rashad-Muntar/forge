package httpserver

import (
	"net/http"
	"strconv"
	"time"

	"github.com/go-chi/chi/v5"
)

type Server struct {
	httpServer *http.Server
}

func New(port int, handler http.Handler) *Server {
	router := chi.NewRouter()

	router.Mount("/", handler)

	return &Server{
		httpServer: &http.Server{
			Addr:              ":" + strconv.Itoa(port),
			Handler:           router,
			ReadHeaderTimeout: 5 * time.Second,
			ReadTimeout:       15 * time.Second,
			WriteTimeout:      15 * time.Second,
			IdleTimeout:       60 * time.Second,
		},
	}
}

func (s *Server) Start() error {
	return s.httpServer.ListenAndServe()
}

func (s *Server) Shutdown() error {
	return s.httpServer.Close()
}