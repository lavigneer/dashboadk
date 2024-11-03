package services

import (
	"context"

	"dashboardk/internal/config"

	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/fields"
)

func NewNodeService(app *config.Application) *NodeService {
	return &NodeService{app: app}
}

type NodeService struct {
	app *config.Application
}

func (n *NodeService) GetAll(ctx context.Context) (*v1.NodeList, error) {
	n.app.ClientSet.CoreV1().Pods("").List(ctx, metav1.ListOptions{})
	return n.app.ClientSet.CoreV1().Nodes().List(ctx, metav1.ListOptions{})
}

func (n *NodeService) GetPodsForNode(ctx context.Context, nodeName string) (*v1.PodList, error) {
	fieldSelector := fields.OneTermEqualSelector("spec.nodeName", nodeName)
	return n.app.ClientSet.CoreV1().Pods("").List(ctx, metav1.ListOptions{FieldSelector: fieldSelector.String()})
}
