package server

import (
	"net/http"

	"dashboardk/cmd/web"
	"dashboardk/internal/config"
	"dashboardk/internal/handlers"
	"dashboardk/internal/services"
	"dashboardk/pkg/reload"
)

func (s *Server) RegisterRoutes(app *config.Application) http.Handler {
	mux := http.NewServeMux()
	fileServer := http.FileServer(http.FS(web.Files))
	mux.Handle("/assets/", fileServer)

	pageReloader := reload.New("/__dev/")
	pageReloader.RegisterRoutes(mux)

	nodeService := services.NewNodeService(app)
	nodeHandler := handlers.NewNodeHandler(app, nodeService)
	nodeHandler.RegisterRoutes(mux)

	namespaceService := services.NewNamespaceService(app)
	namespaceHandler := handlers.NewNamespaceHandler(app, namespaceService)
	namespaceHandler.RegisterRoutes(mux)

	return s.logMiddleware(s.appEnvMiddleware((mux)))
}
