package server

import (
	"net/http"

	"dashboardk/cmd/web"
	"dashboardk/internal/config"
	"dashboardk/internal/handlers"
	"dashboardk/internal/services"
)

func (s *Server) RegisterRoutes(app *config.Application) http.Handler {
	mux := http.NewServeMux()
	fileServer := http.FileServer(http.FS(web.Files))
	mux.Handle("/assets/", fileServer)

	mux.Handle(app.PageReloader.Path, app.PageReloader)

	nodeService := services.NewNodeService(app)
	nodeHandler := handlers.NewNodeHandler(app, nodeService)
	nodeHandler.RegisterRoutes(mux)

	namespaceService := services.NewNamespaceService(app)
	namespaceHandler := handlers.NewNamespaceHandler(app, namespaceService)
	namespaceHandler.RegisterRoutes(mux)

	return s.logMiddleware(s.appEnvMiddleware((mux)))
}
