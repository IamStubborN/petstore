package workers

import (
	"context"

	"github.com/IamStubborN/petstore/config"
)

type Worker interface {
	Run(ctx context.Context)
}

func GenerateWorkers(cfg *config.Config) []Worker {
	var workers []Worker
	for _, name := range cfg.Services {
		worker := generateWorker(name, cfg)
		if worker != nil {
			workers = append(workers, worker)
		}
	}

	return workers
}

func generateWorker(workerName string, cfg *config.Config) Worker {
	var worker Worker
	switch workerName {
	case "api":
		worker = newAPIWorker(cfg)
	case "invoice":
		worker = newInvoiceWorker(cfg)
	}

	return worker
}
