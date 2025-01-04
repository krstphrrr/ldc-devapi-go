package routes

import (
	"database/sql"
	"net/http"
	"go-api-app/internal/services"
	"go-api-app/config"
)

func RegisterRoutes(mux *http.ServeMux, db *sql.DB) {
	for endpoint := range config.EndpointToTableMap {
		mux.HandleFunc(endpoint, func(w http.ResponseWriter, r *http.Request) {
			services.GenericDynamicDataHandler(w, r, db)
		})
	}
}
