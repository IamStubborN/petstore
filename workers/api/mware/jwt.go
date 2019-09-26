package mware

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/IamStubborN/petstore/workers/api/auth"
	"go.uber.org/zap"
)

func JWT(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		bearer := r.Header.Get("Authorization")
		prefixLen := len("Bearer ")
		if len(bearer) <= prefixLen {
			zap.L().Info("Authorization header is too short")
			if _, err := w.Write([]byte(`{"error":"Authorization header is too short"}`)); err != nil {
				zap.L().Info("can't write to response", zap.Error(err))
			}
			return
		}
		tokenRaw := bearer[prefixLen:]

		if auth.IsTokenInBlackList(tokenRaw) {
			zap.L().Info("Token in blacklist")
			if _, err := w.Write([]byte(`{"error":"Authorization token in blacklist"}`)); err != nil {
				zap.L().Info("can't write to response", zap.Error(err))
			}
			return
		}

		claims, err := auth.ParseToken(tokenRaw)
		if err != nil {
			zap.L().Info("can't parse token", zap.Error(err))
			resp := fmt.Sprintf(`{"error":"%s"}`, err.Error())
			if _, err := w.Write([]byte(resp)); err != nil {
				zap.L().Info("can't write to response", zap.Error(err))
			}
			return
		}

		if strings.Contains(claims.AllowMethods, r.Method) {
			next.ServeHTTP(w, r)
		}
	})
}
