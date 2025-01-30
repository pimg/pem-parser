package app

import "log/slog"

type Application struct {
	PEMHandler *PEMHandler
}

func NewApplication(logger *slog.Logger) *Application {
	return &Application{
		PEMHandler: NewPEMHandler(logger),
	}
}
