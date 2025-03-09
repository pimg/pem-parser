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

//go:embed assets/**
var assetsFolder embed.FS

type Server struct {
	logger    *slog.Logger
	Http      *http.Server
	Templates *template.Template
	App       *app.Application
}

const maxRequestBytes = 2 * 1000 * 1000 // 2MB seems large enough for PEM files

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
	router.Handle("POST /", http.MaxBytesHandler(http.HandlerFunc(server.pemParserHandler), maxRequestBytes))
	router.Handle("GET /assets/{path...}", http.FileServer(http.FS(assetsFolder)))

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
