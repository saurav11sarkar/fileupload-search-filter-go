package routes

import (
	"net/http"

	"github.com/saurav11sarkar/001practic/internal/auth"
	"github.com/saurav11sarkar/001practic/internal/config"
	"github.com/saurav11sarkar/001practic/internal/middlewares"
	"github.com/saurav11sarkar/001practic/internal/user"
)

func New(d Dependency, cfg config.Config) http.Handler {
	mux := http.NewServeMux()

	auth.AuthRouters(mux, d.Auth)
	user.UserRoutes(mux, d.User)

	return middlewares.Chain(mux, middlewares.CORS(cfg.CORSOrigin), middlewares.RequestID, middlewares.Logger, middlewares.Recover)
}
