package app

type Application struct {
	PEMHandler *PEMHandler
}

func NewApplication() *Application {
	return &Application{
		PEMHandler: NewPEMHandler(),
	}
}
