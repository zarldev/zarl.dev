package http

import (
	"fmt"

	"github.com/labstack/echo/v4"
	"github.com/zarldev/zarldotdev/pkg/repo"
	"github.com/zarldev/zarldotdev/pkg/transport/http/handlers"
	_ "modernc.org/sqlite"
)

type Server struct {
	Router *echo.Echo
	Config Config
}

type Config struct {
	Host       string
	Port       int
	RepoConfig repo.Config
}

var defaultRepoConfig = repo.Config{
	Connection: "./zarldotdev.db",
}

func defaultConfig() Config {
	return Config{
		Host:       "localhost",
		Port:       8080,
		RepoConfig: defaultRepoConfig,
	}
}

func (c *Config) Addr() string {
	return fmt.Sprintf("%s:%d", c.Host, c.Port)
}

func NewDefault() (*Server, error) {
	svc := &Server{
		Router: echo.New(),
		Config: defaultConfig(),
	}
	svc.RegisterRoutes()
	return svc, nil
}

func New(c Config) (*Server, error) {
	router := echo.New()
	svc := &Server{
		Router: router,
		Config: c,
	}
	err := svc.RegisterRoutes()
	if err != nil {
		return nil, err
	}

	return svc, nil
}

func (s *Server) Start() error {
	s.Router.HideBanner = true
	return s.Router.Start(s.Config.Addr())
}

func (s *Server) Stop() error {
	return s.Router.Close()
}

func (s *Server) RegisterRoutes() error {

	staticHandler := handlers.NewStaticHandler()
	staticHandler.RegisterRoutes(s.Router)

	indexHandler := handlers.NewIndexHandler()
	indexHandler.RegisterRoutes(s.Router)

	aboutHandler := handlers.NewAboutHandler()
	aboutHandler.RegisterRoutes(s.Router)

	contactHandler := handlers.NewContactHandler()
	contactHandler.RegisterRoutes(s.Router)

	articleHandler, err := handlers.NewArticleHandler(s.Config.RepoConfig)
	if err != nil {
		return err
	}
	articleHandler.RegisterRoutes(s.Router)

	adminHandeler, err := handlers.NewAdminHandler(s.Config.RepoConfig)
	if err != nil {
		return err
	}
	adminHandeler.RegisterRoutes(s.Router)

	clapsHandler, err := handlers.NewClapsHandler(s.Config.RepoConfig)
	if err != nil {
		return err
	}
	clapsHandler.RegisterRoutes(s.Router)
	return nil
}
