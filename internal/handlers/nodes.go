package handlers

import (
	"context"
	"log/slog"
	"net/http"

	"dashboardk/cmd/web"
	"dashboardk/internal/config"

	v1 "k8s.io/api/core/v1"
)

type NodeService interface {
	GetAll(ctx context.Context) (*v1.NodeList, error)
}

func NewNodeHandler(app *config.Application, ns NodeService) *NodeHandler {
	return &NodeHandler{app: app, nodeService: ns}
}

type NodeHandler struct {
	app         *config.Application
	nodeService NodeService
}

func (h *NodeHandler) RegisterRoutes(mux *http.ServeMux) {
	mux.HandleFunc("GET /nodes", h.Get)
}

func (h *NodeHandler) Get(w http.ResponseWriter, r *http.Request) {
	nodes, err := h.nodeService.GetAll(r.Context())
	if err != nil {
		h.app.Logger.Error("failed to get nodes", slog.Any("error", err))
		http.Error(w, "failed to get nodes", http.StatusInternalServerError)
		return
	}
	web.NodesList(nodes).Render(r.Context(), w)
}
