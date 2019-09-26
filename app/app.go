package app

import (
	"context"
	"os"
	"os/signal"
	"sync"
	"syscall"

	"github.com/IamStubborN/petstore/config"
	"github.com/IamStubborN/petstore/db"
	"github.com/IamStubborN/petstore/fileserver"
	"github.com/IamStubborN/petstore/logger"
	"github.com/IamStubborN/petstore/workers"
	"github.com/IamStubborN/petstore/workers/api/auth"

	"go.uber.org/zap"
)

type App struct {
	Logger  *zap.Logger
	Workers []workers.Worker
}

func NewApp() *App {
	app := &App{}
	cfg := initConfig()
	app.Logger = initLogger(cfg)
	app.Workers = initWorkers(cfg)
	db.InitDatabase(cfg)
	fileserver.InitMinio(cfg)
	auth.InitJWTAuth(cfg)

	return app
}

func initConfig() *config.Config {
	cfg, err := config.LoadConfig()
	if err != nil {
		zap.NewExample().Fatal("can't load config")
	}

	return cfg
}

func initLogger(cfg *config.Config) *zap.Logger {
	log, err := logger.NewLogger(cfg.Logger)
	if err != nil {
		zap.NewExample().Fatal("can't initialize logger")
	}

	return log
}

func initWorkers(cfg *config.Config) []workers.Worker {
	return workers.GenerateWorkers(cfg)
}

func (app *App) Run() {
	wg := &sync.WaitGroup{}

	ctx, cancel := context.WithCancel(context.Background())

	for _, service := range app.Workers {
		wg.Add(1)
		go func(service workers.Worker) {
			defer wg.Done()
			service.Run(ctx)
		}(service)
	}

	gracefulShutdown(cancel)
	wg.Wait()
	closeAll()
}

func gracefulShutdown(cancel context.CancelFunc) {
	c := make(chan os.Signal, 1)
	signal.Notify(c, syscall.SIGTERM, syscall.SIGINT)
	<-c
	close(c)
	cancel()
}

func closeAll() {
	if err := db.Close(); err != nil {
		zap.L().Fatal(err.Error())
	}

	zap.L().Info("database connection closed")
}
