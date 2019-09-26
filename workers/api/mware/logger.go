package mware

import (
	"net/http"

	"go.uber.org/zap"

	"github.com/go-chi/chi/middleware"
)

func RequestLogger(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ww := middleware.NewWrapResponseWriter(w, r.ProtoMajor)

		defer func() {
			zap.L().Info("served",
				zap.String("path", r.URL.Path),
				zap.String("method", r.Method),
				zap.Int("status", ww.Status()),
				zap.Int("size", ww.BytesWritten()))
		}()
		next.ServeHTTP(ww, r)
	})
}
