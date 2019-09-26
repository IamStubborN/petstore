package workers

import (
	"context"
	"net/http"
	"strconv"
	"time"

	"github.com/IamStubborN/petstore/config"
	"github.com/IamStubborN/petstore/workers/api"
	"github.com/go-chi/chi"
	"go.uber.org/zap"
)

type APIWorker struct {
	Port     int
	WTimeout time.Duration
	RTimeout time.Duration
	GTimeout time.Duration
}

func newAPIWorker(cfg *config.Config) Worker {
	return &APIWorker{
		Port:     cfg.API.Port,
		WTimeout: time.Duration(cfg.API.WTimeout) * time.Second,
		RTimeout: time.Duration(cfg.API.RTimeout) * time.Second,
		GTimeout: time.Duration(cfg.API.GTimeout) * time.Second,
	}
}

func (aw *APIWorker) Run(ctx context.Context) {
	srv := &http.Server{
		Addr:         ":" + strconv.Itoa(aw.Port),
		Handler:      chi.ServerBaseContext(ctx, api.NewRouter()),
		WriteTimeout: aw.WTimeout,
		ReadTimeout:  aw.RTimeout,
	}

	go func() {
		<-ctx.Done()

		ctxShutDown, cancel := context.WithTimeout(context.Background(), aw.GTimeout)
		defer cancel()

		if err := srv.Shutdown(ctxShutDown); err != nil {
			zap.L().Info(err.Error())
		}
	}()

	if err := srv.ListenAndServe(); err != nil {
		zap.L().Info(err.Error())
	}
}
