package app

import (
	"net/http"
)

type RouteInfo struct {
	Method string
	Path   string
}

type App struct {
	Mux    *http.ServeMux
	Routes []RouteInfo
}

func (a *App) AddRoute(method, path string, handler http.Handler) {
	a.Routes = append(a.Routes, RouteInfo{Method: method, Path: path})
	a.Mux.Handle(method+" "+path, handler)
}
