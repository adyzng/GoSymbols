package main

import (
	"net/http"

	"github.com/adyzng/goPDB/restful/v1"
	"github.com/gorilla/mux"
)

// Route define the basic route
//
type Route struct {
	Name    string
	Method  []string
	Pattern string
	Handler http.HandlerFunc
}

var routes = []Route{
	{
		Name:    "Index",
		Method:  []string{"GET"},
		Pattern: "/",
		Handler: IndexHandle,
	},
	{
		Name:    "RestBranchList",
		Method:  []string{"GET"},
		Pattern: "/api/v1/branch",
		Handler: v1.RestBranchList,
	},
	{
		Name:    "RestBuildList",
		Method:  []string{"GET"},
		Pattern: "/api/v1/branch/{name}",
		Handler: v1.RestBuildList,
	},
	{
		Name:    "RestSymbolList",
		Method:  []string{"GET"},
		Pattern: "/api/v1/branch/{name}/{bid}",
		Handler: v1.RestSymbolList,
	},
	{
		Name:    "DownloadSymbol",
		Method:  []string{"GET"},
		Pattern: "/api/v1/symbol/{name}/{hash}",
		Handler: v1.DownloadSymbol,
	},
}

// NewRouter return the registered router
//
func NewRouter() *mux.Router {
	router := mux.NewRouter()
	router.StrictSlash(true)

	// static files handler
	router.
		PathPrefix("/public/").
		Handler(http.StripPrefix("/public/", StaticHandler("./public/")))

	// predefined handler
	for _, route := range routes {
		logHandler := LogHandler(route.Handler, route.Name)
		router.
			Methods(route.Method...).
			Path(route.Pattern).
			Handler(logHandler).
			Name(route.Name)
	}

	return router
}
