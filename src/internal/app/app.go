package app

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"

	"github.com/MatiasKopp/prosig-code-challenge/posts"
	"github.com/caarlos0/env/v11"
	"github.com/go-chi/chi/v5"
	_ "github.com/mattn/go-sqlite3"
)

// Config App configuration structure.
type Config struct {
	Port       string `env:"APP_PORT"`
	DBLocation string `env:"DB_LOCATION"`
}

// App Represents productive app.
type App struct {
	Config Config
	Router chi.Router

	// Handlers
	PostsHTTPAdapter posts.HTTPAdapter
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

	a.Router.Route("/api", func(api chi.Router) {
		api.Get("/posts", a.PostsHTTPAdapter.GetAllPosts)
	})
}

// bootstrap Bootstraps handlers
func (a *App) bootstrap() {
	db, err := sql.Open("sqlite3", a.Config.DBLocation)
	if err != nil {
		log.Fatal(err)
	}

	repository, err := posts.NewRepository(db)
	if err != nil {
		panic("error creating repository")
	}

	service, err := posts.NewService(repository)
	if err != nil {
		panic("error creating service")
	}

	httpAdapter, err := posts.NewHTTPAdapter(service)
	if err != nil {
		panic("error creating service")
	}

	a.PostsHTTPAdapter = httpAdapter
}
