package app

import (
	"net/http"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/saurav11sarkar/001practic/internal/auth"
	"github.com/saurav11sarkar/001practic/internal/config"
	"github.com/saurav11sarkar/001practic/internal/email"
	"github.com/saurav11sarkar/001practic/internal/routes"
	"github.com/saurav11sarkar/001practic/internal/user"
	"github.com/saurav11sarkar/001practic/internal/utils"
)

func NewHandler(db *pgxpool.Pool, cfg config.Config) (http.Handler, error) {
	email := email.NewEmail(cfg)

	cloudenary, err := utils.NewCloudinaryService(cfg)
	if err != nil {
		return nil, err
	}

	authRepo := auth.NewRepository(db)
	authService := auth.NewService(authRepo, cfg, email)
	authHandler := auth.NewHandler(authService)

	userRepo := user.NewRepository(db)
	userService := user.NewService(userRepo,cloudenary)
	userHandler := user.NewHandler(userService)

	deps := routes.Dependency{
		Auth: authHandler,
		User: userHandler,
	}

	return routes.New(deps, cfg), nil
}
