package server

import (
	"fmt"
	"net/http"
	"os"
	"strconv"
	"time"

	"dashboardk/internal/config"

	_ "github.com/joho/godotenv/autoload"
)

type Server struct {
	app    *config.Application
	port   int
	appEnv string
}

func NewServer(app *config.Application) *http.Server {
	appEnv := os.Getenv("APP_ENV")
	port, _ := strconv.Atoi(os.Getenv("PORT"))
	NewServer := &Server{
		port:   port,
		app:    app,
		appEnv: appEnv,
	}

	// Declare Server config
	server := &http.Server{
		Addr:         fmt.Sprintf(":%d", NewServer.port),
		Handler:      NewServer.RegisterRoutes(app),
		IdleTimeout:  time.Minute,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 30 * time.Second,
	}

	return server
}
