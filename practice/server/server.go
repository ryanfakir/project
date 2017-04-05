package server

import (
	"fmt"
	"net"
	"net/http"

	"github.com/rs/cors"
)

// Server struct
type Server struct {
	ln      net.Listener
	Handler *Handler
	Addr    string
}

// Open opens a socket and serves the HTTP server
func (s *Server) Open() error {
	ln, err := net.Listen("tcp", s.Addr)
	if err != nil {
		return err
	}
	s.ln = ln
	fmt.Printf("server start at %s\n", s.Addr)
	handler := cors.Default().Handler(s.Handler)
	http.Serve(s.ln, handler)
	return nil
}

// Close server
func (s *Server) Close() error {
	if s.ln != nil {
		s.ln.Close()
	}
	return nil
}

// New server instance
func New() *Server {
	return &Server{
		Handler: NewHandler(),
		Addr:    "localhost:8080",
	}
}
