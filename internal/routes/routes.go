package routes

import (
	"database/sql"
	"encoding/json"
	"net/http"

	"go-api-app/config"
	"go-api-app/internal/middleware"
	"go-api-app/internal/services"
	"go-api-app/internal/version"
)

func RegisterRoutes(mux *http.ServeMux, db *sql.DB) {
	for endpoint := range config.EndpointToTableMap {
		// Handling endpoints that don't require middleware 

		// if endpoint == "/tblProject" {
		// 	mux.Handle(endpoint, http.HandlerFunc(services.tblProjectHandler))

		// } else {
		// 	mux.Handle(endpoint, middleware.AuthMiddleware(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// 		services.GenericDynamicDataHandler(w, r, db)
		// 	})))
		// }
		// mux.Handle(endpoint, middleware.AuthMiddleware(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// 	services.GenericDynamicDataHandler(w, r, db)
		// })))
		mux.Handle(endpoint, middleware.AuthMiddleware(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			// Differentiate between GET and POST requests in the handler
			switch r.Method {
			case http.MethodGet:
				services.GenericDynamicDataHandler(w, r, db) // Handle GET
			case http.MethodPost:
				services.GenericDynamicDataPostHandler(w, r, db) // Handle POST
			default:
				http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
			}
		})))
	}
	mux.HandleFunc("/version", versionHandler)
}
func versionHandler(w http.ResponseWriter, r *http.Request) {
    versionInfo := map[string]string{
        "appName":   version.AppName,
        "version":   version.Version,
        "buildDate": version.BuildDate,
    }
    json.NewEncoder(w).Encode(versionInfo)
}