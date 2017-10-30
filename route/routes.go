package route

import (
	"net/http"

	"github.com/adyzng/GoSymbols/config"
	"github.com/adyzng/GoSymbols/restful/v1"
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
		Pattern: "/api/branchs",
		Handler: v1.RestBranchList,
	},
	{
		Name:    "RestBuildList",
		Method:  []string{"GET"},
		Pattern: "/api/branchs/{name}",
		Handler: v1.RestBuildList,
	},
	{
		Name:    "RestSymbolList",
		Method:  []string{"GET"},
		Pattern: "/api/branchs/{name}/{bid}",
		Handler: v1.RestSymbolList,
	},
	{
		Name:    "DownloadSymbol",
		Method:  []string{"GET"},
		Pattern: "/api/symbol/{branch}/{hash}/{name}",
		Handler: v1.DownloadSymbol,
	},
	{
		Name:    "ModifyBranch",
		Method:  []string{"POST"},
		Pattern: "/api/branch",
		Handler: v1.ModifyBranch,
	},
	{
		Name:    "ValidateBranch",
		Method:  []string{"POST"},
		Pattern: "/api/branch/check",
		Handler: v1.ValidateBranch,
	},
	{
		Name:    "DeleteBranch",
		Method:  []string{"DELETE"},
		Pattern: "/api/branch/{name}",
		Handler: v1.DeleteBranch,
	},
}

// NewRouter return the registered router
//
func NewRouter() *mux.Router {
	router := mux.NewRouter()
	router.StrictSlash(true)

	// static files handler
	router.
		PathPrefix("/assets/").
		Handler(http.StripPrefix("/assets/", StaticHandler(config.Public)))

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
