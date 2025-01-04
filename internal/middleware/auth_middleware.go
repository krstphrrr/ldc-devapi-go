package middleware

import (
	"context"
	"net/http"
	"strings"

	"go-api-app/internal/auth"
	"go-api-app/config"
)

// Custom context key
type contextKey string

const tenantKey contextKey = "tenant"

func AuthMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		tenant := "public" // Default tenant

		authHeader := r.Header.Get("Authorization")
		if authHeader != "" {
			token := strings.TrimPrefix(authHeader, "Bearer ")

			// Use Cognito for verification
			verifiedToken, err := auth.VerifyJWTWithCognito(
				token,
				config.Config.AwsCognito.UserPoolId,
				config.Config.AwsCognito.ClientId,
			)
			if err == nil {
				tenant = groupDiscrimination(verifiedToken)
			}
		}

		ctx := context.WithValue(r.Context(), tenantKey, tenant)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

// groupDiscrimination determines the tenant based on groups in the verified JWT
func groupDiscrimination(verifiedToken map[string]interface{}) string {
	groups, ok := verifiedToken["cognito:groups"].([]interface{})
	if !ok || len(groups) == 0 {
		return "public"
	}

	for _, group := range groups {
		switch group {
		case "LEGAL":
			return "legal"
		case "PUBLICATION":
			return "publication"
		}
	}

	return "public"
}
