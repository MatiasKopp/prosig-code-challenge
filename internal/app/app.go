package app

import (
	"fmt"
	"net/http"

	"github.com/caarlos0/env/v11"
	"github.com/go-chi/chi/v5"
)

// Config App configuration structure.
type Config struct {
	Port string `env:"APP_PORT"`
}

// App Represents productive app.
type App struct {
	Config Config
	Router chi.Router

	// Handlers
}

// New Returns new productive app implementation
func New() *App {
	// parse
	var cfg Config
	err := env.Parse(&cfg)
	if err != nil {
		panic(fmt.Errorf("error parsing env config: %s", err))
	}

	app := &App{
		Config: cfg,
		Router: chi.NewRouter(),
	}
	app.bootstrap()
	app.mapRoutes()

	return app
}

// HealthCheck Simple health check.
func HealthCheck(w http.ResponseWriter, _ *http.Request) {
	w.Write([]byte("pong"))
}

// Start Starts listening at configured port.
func (a *App) Start() {
	fmt.Printf("App listening at port %s...", a.Config.Port)
	http.ListenAndServe(a.Config.Port, a.Router)
}

// mapRoutes Maps routes to handlers
func (a *App) mapRoutes() {
	a.Router.Get("/ping", HealthCheck)
}

// bootstrap Bootstraps handlers
func (a *App) bootstrap() {

}
