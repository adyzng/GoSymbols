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

var resRoutes = []Route{
	{
		Name:    "Index",
		Method:  []string{"GET"},
		Pattern: "/",
		Handler: IndexHandle,
	},
}

var apiRoutes = []Route{
	{
		Name:    "GetBranchList",
		Method:  []string{"GET"},
		Pattern: "/branches",
		Handler: v1.RestBranchList,
	},
	{
		Name:    "GetBuildList",
		Method:  []string{"GET"},
		Pattern: "/branches/{name}",
		Handler: v1.RestBuildList,
	},
	{
		Name:    "GetSymbolList",
		Method:  []string{"GET"},
		Pattern: "/branches/{name}/{bid}",
		Handler: v1.RestSymbolList,
	},
	{
		Name:    "DownloadSymbol",
		Method:  []string{"GET"},
		Pattern: "/symbol/{branch}/{hash}/{name}",
		Handler: v1.DownloadSymbol,
	},
	{
		Name:    "ModifyBranch",
		Method:  []string{"POST"},
		Pattern: "/branches/modify",
		Handler: v1.ModifyBranch,
	},
	{
		Name:    "ValidateBranch",
		Method:  []string{"POST"},
		Pattern: "/branches/check",
		Handler: v1.ValidateBranch,
	},
	{
		Name:    "DeleteBranch",
		Method:  []string{"DELETE"},
		Pattern: "/branches/{name}",
		Handler: v1.DeleteBranch,
	},
	{
		Name:    "FetchTodayMessage",
		Method:  []string{"GET"},
		Pattern: "/messages",
		Handler: v1.FetchTodayMsg,
	},
	{
		Name:    "Login",
		Method:  []string{"GET"},
		Pattern: "/auth/login",
		Handler: v1.AuthLogin,
	},
	{
		Name:    "Authorize",
		Method:  []string{"POST"},
		Pattern: "/auth/authorize",
		Handler: v1.Authorize,
	},
	{
		Name:    "Logout",
		Method:  []string{"GET"},
		Pattern: "/auth/logout",
		Handler: v1.AuthLogout,
	},
	{
		Name:    "UserProfile",
		Method:  []string{"GET"},
		Pattern: "/user/profile",
		Handler: v1.GetUserProfile,
	},
	{
		Name:    "UserPhoto",
		Method:  []string{"GET"},
		Pattern: "/user/photo",
		Handler: v1.GetUserPhoto,
	},
}

// NewRouter return the registered router
//
func NewRouter() *mux.Router {
	router := mux.NewRouter()
	router.StrictSlash(true)

	// static files handler
	router.
		PathPrefix("/static/").
		Handler(StaticHandler(config.WebRoot))
		//Handler(http.StripPrefix("/static/", StaticHandler(config.WebRoot)))

	// normal handler
	for _, route := range resRoutes {
		logHandler := LogHandler(route.Handler, route.Name)
		router.
			Methods(route.Method...).
			Path(route.Pattern).
			Handler(logHandler).
			Name(route.Name)
	}

	// restful api handler
	for _, route := range apiRoutes {
		logHandler := LogHandler(route.Handler, route.Name)
		router.PathPrefix("/api/").
			Methods(route.Method...).
			Path(route.Pattern).
			Handler(logHandler).
			Name(route.Name)
	}

	return router
}
