package middleware

import (
	"context"
	"log"
	"net/http"
	"strings"

	"go-api-app/internal/auth"
	"go-api-app/config"
)

// Custom context keys
type contextKey string

const (
	TenantKey contextKey = "tenant"
	UserIDKey contextKey = "user_id"
)

// AuthMiddleware verifies the JWT and extracts the user ID
func AuthMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		log.Println("AuthMiddleware: Processing request.")

		tenant := "public" // Default tenant
		userID := "anonymous" // Default user
		authHeader := r.Header.Get("Authorization")

		if authHeader != "" {
			log.Println("AuthMiddleware: Authorization header found.")
			token := strings.TrimPrefix(authHeader, "Bearer ")

			// Verify the JWT with Cognito
			verifiedToken, err := auth.VerifyJWTWithCognito(
				token,
				config.Config.AwsCognito.UserPoolId,
				config.Config.AwsCognito.ClientId,
			)

			if err != nil {
				log.Printf("AuthMiddleware: JWT verification failed. Error: %v\n", err)
			} else {
				log.Println("AuthMiddleware: JWT verified successfully.")
				tenant = groupDiscrimination(verifiedToken)

				// Extract the Cognito User ID (sub)
				if sub, ok := verifiedToken["sub"].(string); ok {
					userID = sub
				}
				log.Printf("AuthMiddleware: User ID: %s, Tenant: %s\n", userID, tenant)
			}
		}

		// Attach tenant & user_id to context
		ctx := context.WithValue(r.Context(), TenantKey, tenant)
		ctx = context.WithValue(ctx, UserIDKey, userID)

		next.ServeHTTP(w, r.WithContext(ctx))
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