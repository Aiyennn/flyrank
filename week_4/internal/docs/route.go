package docs

import (
	_ "embed"
	"net/http"

	"github.com/go-chi/chi/v5"
)

//go:embed openapi.json
var openapiJSON []byte

//go:embed swagger.html
var swaggerHTML []byte

type DocsRoute struct {
}

func NewDocsRoute() *DocsRoute {
	return &DocsRoute{}
}

func (dr *DocsRoute) New() chi.Router {
	r := chi.NewRouter()

	r.Get("/", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/html; charset=utf-8")
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write(swaggerHTML)
	})

	r.Get("/openapi.json", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write(openapiJSON)
	})

	return r
}
