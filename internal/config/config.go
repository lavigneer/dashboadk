package config

import (
	"log/slog"

	autorefresh "github.com/lavigneer/browser-autorefresh"
	"k8s.io/client-go/kubernetes"
)

type Application struct {
	Logger       *slog.Logger
	ClientSet    *kubernetes.Clientset
	PageReloader *autorefresh.PageReloader
	Env          string
}

type appContextKey string

var AppContextClass = appContextKey("appContext")
