package server

import (
	"fmt"
	"log/slog"
	"mamlaka/config"
	"mamlaka/internal/pkg/database"
	"mamlaka/internal/pkg/templates"
	"net/http"
	"os"
	"strconv"
	"time"

	"github.com/matcornic/hermes/v2"

	_ "github.com/joho/godotenv/autoload"
)

type Server struct {
	port   int
	db     database.Service
	logger *slog.Logger
	config config.Config
	hermes *hermes.Hermes
}

func NewServer() *http.Server {
	conf := config.ReadConfigFromEnv()

	port, _ := strconv.Atoi(os.Getenv("PORT"))
	NewServer := &Server{
		port:   port,
		db:     database.New(conf),
		logger: slog.New(slog.NewJSONHandler(os.Stderr, nil)),
		config: conf,
		hermes: templates.InitializeHermes(),
	}

	// Declare Server config
	server := &http.Server{
		Addr:         fmt.Sprintf(":%d", NewServer.port),
		Handler:      NewServer.RegisterRoutes(),
		IdleTimeout:  time.Minute,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 30 * time.Second,
	}
	return server
}
