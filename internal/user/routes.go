package user

import (
	"net/http"

	"github.com/saurav11sarkar/001practic/internal/middlewares"
)

func UserRoutes(mux *http.ServeMux, h *Handler) {
	mux.Handle("GET /user/profile", middlewares.Auth("admin", "user")(http.HandlerFunc(h.GetProfile)))
	mux.Handle("GET /user", http.HandlerFunc(h.GetAllUser))
	mux.Handle("GET /user/{id}", http.HandlerFunc(h.GetSingleUser))
	mux.Handle("PUT /user/profile", middlewares.Auth("admin", "user")(http.HandlerFunc(h.UpdateProfile)))
	mux.Handle("PUT /user/update-profile", middlewares.Auth("admin", "user")(http.HandlerFunc(h.UpdateUser)))
	mux.Handle("DELETE /user/{id}", middlewares.Auth("admin")(http.HandlerFunc(h.DeleteUser)))
}
