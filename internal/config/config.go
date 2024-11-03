package config

import (
	"log/slog"

	"k8s.io/client-go/kubernetes"
)

type Application struct {
	Logger *slog.Logger
	ClientSet *kubernetes.Clientset
}

type appEnvContextKey string

var AppEnvContextClass = appEnvContextKey("appEnv")
