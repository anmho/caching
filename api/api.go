package api

import "net/http"


func New() *http.ServeMux{
	mux := http.NewServeMux()

	registerRoutes(mux)

	return mux
}

