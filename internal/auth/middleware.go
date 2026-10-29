package auth

import (
	"context"
	"net/http"
	"strings"
	"uuid"
)



type APIKeyAuthenticator interface{
	Authenticate(
		ctx context.Context,
		apiKey string,
	)(uuid.UUID, error)
}

func TenantAuthMiddleware(
    authenticator APIKeyAuthenticator,
) func(http.Handler) http.Handler {
    return func(next http.Handler) http.Handler {
        return http.HandlerFunc(func(
            w http.ResponseWriter,
            r *http.Request,
        ) {
            const prefix = "Bearer "

            authHeader := r.Header.Get("Authorization")
            if !strings.HasPrefix(authHeader, prefix) {
                http.Error(
                    w,
                    "missing or invalid authorization header",
                    http.StatusUnauthorized,
                )
                return
            }

            apiKey := strings.TrimSpace(
                strings.TrimPrefix(authHeader, prefix),
            )
            if apiKey == "" {
                http.Error(
                    w,
                    "missing API key",
                    http.StatusUnauthorized,
                )
                return
            }

            tenantID, err := authenticator.Authenticate(
                r.Context(),
                apiKey,
            )
            if err != nil {
                http.Error(
                    w,
                    "unauthorized",
                    http.StatusUnauthorized,
                )
                return
            }

            ctx := WithTenantID(r.Context(), tenantID)
            next.ServeHTTP(w, r.WithContext(ctx))
        })
    }
}
