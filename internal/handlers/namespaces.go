package handlers

import (
	"context"
	"log/slog"
	"net/http"

	"dashboardk/cmd/web"
	"dashboardk/internal/config"

	v1 "k8s.io/api/core/v1"
)

type NamespaceService interface {
	GetAll(ctx context.Context) (*v1.NamespaceList, error)
}

func NewNamespaceHandler(app *config.Application, ns NamespaceService) *NamespaceHandler {
	return &NamespaceHandler{app: app, namespaceService: ns}
}

type NamespaceHandler struct {
	app              *config.Application
	namespaceService NamespaceService
}

func (h *NamespaceHandler) RegisterRoutes(mux *http.ServeMux) {
	mux.HandleFunc("GET /namespaces", h.Get)
}

func (h *NamespaceHandler) Get(w http.ResponseWriter, r *http.Request) {
	nodes, err := h.namespaceService.GetAll(r.Context())
	if err != nil {
		h.app.Logger.Error("failed to get nodes", slog.Any("error", err))
		http.Error(w, "failed to get nodes", http.StatusInternalServerError)
		return
	}
	web.NamespacesList(nodes).Render(r.Context(), w)
}
