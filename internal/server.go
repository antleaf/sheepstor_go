package internal

import (
	"fmt"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"net/http"
)

var Router chi.Router

func NewRouter() chi.Router {
	router := chi.NewRouter()
	router.Use(middleware.Logger)
	router.Use(middleware.StripSlashes)
	router.Use(middleware.Throttle(10))
	// TODO: figure out if it is possible to use this CORS module to add common HTTP headers to all HTTP Responses. Otherwise write a middleware handler to do this.
	router.Get("/", DefaultHandler)
	router.Post("/update/{website_id}", GitHubWebHookHandler)
	return router
}

func RunServer(port int) {
	Router = NewRouter()
	Log.Infof("Running as HTTP Process on port %d", port)
	err := http.ListenAndServe(fmt.Sprintf(":%v", port), Router)
	if err != nil {
		Log.Error(err.Error())
	}
}
