package reload

import (
	_ "embed"
	"fmt"
	"html/template"
	"net/http"
	"net/url"
	"time"

	"github.com/coder/websocket"
)

//go:embed "reload.js"
var Script string

type PageReloader struct {
	path string
}

func New(path string) *PageReloader {
	return &PageReloader{path: path}
}

func (p *PageReloader) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	socket, err := websocket.Accept(w, r, nil)
	if err != nil {
		_, _ = w.Write([]byte("could not open websocket"))
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	defer socket.Close(websocket.StatusGoingAway, "server closing websocket")
	ctx := r.Context()
	socketCtx := socket.CloseRead(ctx)
	for {
		socket.Ping(socketCtx)
		payload := fmt.Sprintf("server timestamp: %d", time.Now().UnixNano())
		err := socket.Write(socketCtx, websocket.MessageText, []byte(payload))
		if err != nil {
			break
		}
		time.Sleep(time.Second * 2)
	}
}

func (p *PageReloader) RegisterRoutes(mux *http.ServeMux) error {
	mux.Handle(p.path, p)
	reloadScriptPath, err := url.JoinPath(p.path, "/reload.js")
	if err != nil {
		return err
	}
	mux.HandleFunc(reloadScriptPath, func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, "const reloadPath = '%s'\n", template.JSEscapeString(p.path))
		w.Write([]byte(Script))
	})
	return nil
}
