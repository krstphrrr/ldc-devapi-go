package middleware

import (
	"context"
	"log"
	"net/http"
	"strings"

	"go-api-app/internal/auth"
	"go-api-app/config"
)

// Custom context key
type contextKey string

const TenantKey contextKey = "tenant"

func AuthMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		log.Println("AuthMiddleware: Started processing request.")

		tenant := "public" // Default tenant
		authHeader := r.Header.Get("Authorization")

		if authHeader != "" {
			log.Println("AuthMiddleware: Authorization header found.")
			token := strings.TrimPrefix(authHeader, "Bearer ")

			// Use Cognito for verification
			log.Println("AuthMiddleware: Verifying JWT with Cognito.")
			verifiedToken, err := auth.VerifyJWTWithCognito(
				token,
				config.Config.AwsCognito.UserPoolId,
				config.Config.AwsCognito.ClientId,
			)

			if err != nil {
				log.Printf("AuthMiddleware: JWT verification failed. Error: %v\n", err)
			} else {
				log.Println("AuthMiddleware: JWT verified successfully.")
				log.Printf("AuthMiddleware: Verified token claims: %v\n", verifiedToken)
				tenant = groupDiscrimination(verifiedToken)
				log.Printf("AuthMiddleware: Tenant determined: %s\n", tenant)
			}
		} else {
			log.Println("AuthMiddleware: No Authorization header found. Defaulting to 'public' tenant.")
		}

		// Attach tenant to the context
		ctx := context.WithValue(r.Context(), TenantKey, tenant)
		log.Println("AuthMiddleware: Tenant added to context.")
		next.ServeHTTP(w, r.WithContext(ctx))
		log.Println("AuthMiddleware: Request passed to the next handler.")
	})
}

// groupDiscrimination determines the tenant based on groups in the verified JWT
func groupDiscrimination(verifiedToken map[string]interface{}) string {
	log.Println("groupDiscrimination: Determining tenant based on token groups.")
	groups, ok := verifiedToken["cognito:groups"].([]interface{})
	if !ok || len(groups) == 0 {
		log.Println("groupDiscrimination: No groups found in token. Defaulting to 'public' tenant.")
		return "public"
	}

	for _, group := range groups {
		log.Printf("groupDiscrimination: Processing group: %v\n", group)
		switch group {
		case "LEGAL":
			return "legal"
		case "PUBLICATION":
			return "publication"
		}
	}

	log.Println("groupDiscrimination: No matching group found. Defaulting to 'public' tenant.")
	return "public"
}