package main

import (
	"OnlineMarket/back/Router"
	"net/http"
)

var (
	router *Router.Router
	controllers interface{}
)

func main() {
	router = new(Router.Router)
	router.Manage()
	_ = http.ListenAndServe(":80", nil)
}
