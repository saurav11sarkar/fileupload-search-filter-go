package middlewares

import (
	"context"
	"net/http"
	"strings"

	"github.com/saurav11sarkar/001practic/internal/config"
	"github.com/saurav11sarkar/001practic/internal/utils"
)

func Auth(roles ...string) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			cfg, err := config.MustLoad()
			if err != nil {
				http.Error(w, "Authorization header is missing or invalid", http.StatusUnauthorized)
				return
			}
			parts := strings.SplitN(r.Header.Get("Authorization"), " ", 2)

			if len(parts) != 2 || parts[0] != "Bearer" {
				utils.HanldeError(w, utils.NewAppError(http.StatusUnauthorized, "UNAUTHORIZED", "Authorization header is missing or invalid"))
				return
			}

			token := parts[1]
			claims, err := utils.ParseToken(token, cfg.Auth.JWTAccessSecret, "access")

			if err != nil {
				utils.HanldeError(w, utils.NewAppError(http.StatusUnauthorized, "UNAUTHORIZED", "Invalid access token"))
				return
			}

			if len(roles) > 0 {
				hashRole := false
				for _, role := range roles {
					if role == claims.Role {
						hashRole = true
						break
					}
				}

				if !hashRole {
					utils.HanldeError(w, utils.NewAppError(http.StatusForbidden, "FORBIDDEN", "Insufficient role"))
					return
				}
			}
			next.ServeHTTP(w, r.WithContext(context.WithValue(r.Context(), "user", claims)))
		})
	}
}
