package routes

import (
	"github.com/saurav11sarkar/001practic/internal/auth"
	"github.com/saurav11sarkar/001practic/internal/user"
)

type Dependency struct {
	Auth *auth.Handler
	User *user.Handler
}
