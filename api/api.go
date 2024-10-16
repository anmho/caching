package api

import (
	"database/sql"
	"net/http"
)

func New(db *sql.DB) *http.ServeMux {
	mux := http.NewServeMux()

	registerRoutes(mux, db)

	return mux
}
