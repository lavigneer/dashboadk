package server

import (
	"context"
	"net/http"
	"time"

	"dashboardk/internal/config"
)

func (s *Server) logMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var (
			ip     = r.RemoteAddr
			method = r.Method
			uri    = r.URL.RequestURI()
		)

		requestTime := time.Now()
		next.ServeHTTP(w, r)
		s.app.Logger.Info("handled request", "ip", ip, "method", method, "uri", uri, "duration", time.Since(requestTime))
	})
}

func (s *Server) appEnvMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ctx := context.WithValue(r.Context(), config.AppContextClass, s.app)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}
