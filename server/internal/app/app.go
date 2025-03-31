package app

type App struct {
}

func NewApp() (*App, error) {
	return &App{}, nil
}

func (app *App) Run() int {
	return 1
}
