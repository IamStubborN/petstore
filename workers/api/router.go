package api

import (
	"github.com/IamStubborN/petstore/workers/api/handler"
	"github.com/IamStubborN/petstore/workers/api/mware"
	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
	"github.com/go-chi/cors"
)

func NewRouter() *chi.Mux {
	router := chi.NewRouter()
	router.Use(cors.New(cors.Options{
		AllowedOrigins: []string{"*"},
		AllowedMethods: []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders: []string{"Accept", "Authorization", "Content-Type"},
	}).Handler)

	router.Use(middleware.Recoverer)
	router.Use(mware.RequestLogger)
	router.Use(mware.ResponseDefaultHeaders)

	router.Mount("/debug", middleware.Profiler())

	router.Route("/api/v2", func(r chi.Router) {
		r.Route("/pet", handler.PetHandlers)
		r.Route("/store", handler.StoreHandlers)
		r.Route("/user", handler.UserHandlers)
	})

	return router
}
