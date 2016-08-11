package main

import "net/http"

type Route struct {
	Name        string
	Method      string
	Pattern     string
	HandlerFunc http.HandlerFunc
}

type Routes []Route

var routes = Routes{
	Route{
		"Index",
		"GET",
		"/",
		Index,
	},
	Route{
		"TodoIndex",
		"GET",
		"/todos",
		TodoIndex,
	},
	Route{
		"TodoCreate",
		"POST",
		"/todos",
		TodoCreate,
	},
	Route{
		"TodoShow",
		"GET",
		"/todos/{todoId}",
		TodoShow,
	},
	Route{
		"AccountsIndex",
		"GET",
		"/accounts",
		AccountsIndex,
	},
	Route{
		"AccountShow",
		"GET",
		"/account/{Id}",
		AccountShow,
	},
	Route{
		"AccountCreate",
		"POST",
		"/accounts",
		AccountCreate,
	},
}
