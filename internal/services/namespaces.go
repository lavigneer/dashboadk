package services

import (
	"context"

	"dashboardk/internal/config"

	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func NewNamespaceService(app *config.Application) *NamespaceService {
	return &NamespaceService{app: app}
}

type NamespaceService struct {
	app *config.Application
}

func (n *NamespaceService) GetAll(ctx context.Context) (*v1.NamespaceList, error) {
	return n.app.ClientSet.CoreV1().Namespaces().List(ctx, metav1.ListOptions{})
}
