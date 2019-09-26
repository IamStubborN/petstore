package handler

import (
	"context"
	"encoding/json"
	"net/http"
	"regexp"

	"go.uber.org/zap"
)

type response struct {
	Code    int    `json:"code"`
	Type    string `json:"type"`
	Message string `json:"message"`
}

func respond(w http.ResponseWriter, err error, code int, message string) {
	res := response{
		Code:    code,
		Type:    http.StatusText(code),
		Message: message,
	}

	if err != nil {
		zap.L().Error("respond error", zap.Error(err))
	}

	w.WriteHeader(code)

	if err := json.NewEncoder(w).Encode(&res); err != nil {
		zap.L().Error("can't write response to ResponseWriter", zap.Error(err))
	}
}

func checkErrors(f func() error) {
	if err := f(); err != nil {
		zap.L().Error("error in defer", zap.Error(err))
	}
}

func isNotValid(str string) bool {
	reg := regexp.MustCompile(`^[0-9a-zA-Z]{5,15}$`)
	return !reg.MatchString(str)
}

// Create context 80% of WriteTimeout
func genContext(r *http.Request) (context.Context, context.CancelFunc) {
	writeTimeout := r.Context().Value(http.ServerContextKey).(*http.Server).WriteTimeout
	return context.WithTimeout(context.Background(), writeTimeout*80/100)
}
