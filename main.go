package main

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/httplog"

	"github.com/osemisan/osemisan-resource-server/pkg/handlers"
	"github.com/osemisan/osemisan-resource-server/pkg/middlewares"
)

func main() {
	l := httplog.NewLogger("osemisan-resource-server", httplog.Options{
		JSON: true,
	})
	r := chi.NewRouter()

	r.Use(httplog.RequestLogger(l))
	r.Use(middlewares.VerifyToken)

	r.Get("/", handlers.RootHandler)
	r.Get("/resources", handlers.ResourcesHandler)

	http.ListenAndServe(":9002", r)
}
