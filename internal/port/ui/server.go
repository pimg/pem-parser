package ui

import (
	"context"
	"embed"
	"fmt"
	"html/template"
	"log/slog"
	"net/http"
	"pem-parser/internal/app"
)

//go:embed templates/**
var tplFolder embed.FS

type Server struct {
	logger    *slog.Logger
	Http      *http.Server
	Templates *template.Template
	App       *app.Application
}

func NewServer(logger *slog.Logger, app *app.Application) (*Server, error) {
	tmpl, err := template.New("").ParseFS(tplFolder, "templates/pages/**", "templates/partials/**")
	if err != nil {
		return nil, fmt.Errorf("failed to parse templates: %w", err)
	}

	server := &Server{
		logger:    logger,
		Templates: tmpl,
		App:       app,
	}

	router := http.NewServeMux()

	router.HandleFunc("GET /", server.homeHandler)
	router.HandleFunc("POST /", server.pemParserHandler)

	httpServer := &http.Server{Addr: ":8080", Handler: router}
	server.Http = httpServer
	logger.Debug("Server configured with handlers on port 8080")
	return server, nil
}

func (s *Server) Start() error {
	return s.Http.ListenAndServe()
}

func (s *Server) Stop(ctx context.Context) error {
	return s.Http.Shutdown(ctx)
}
