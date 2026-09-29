package auth

import (
	"net/http"

	"github.com/saurav11sarkar/001practic/internal/utils"
)

type Handler struct {
	service *Service
}

func NewHandler(s *Service) *Handler {
	return &Handler{service: s}
}

func (h *Handler) Register(w http.ResponseWriter, r *http.Request) {
	var req RegisterRequest
	if err := utils.Decode(r, &req); err != nil {
		utils.Error(w, http.StatusBadRequest, "Failed to decode request body")
		return
	}

	if err := utils.ValidationStruct(req); err != nil {
		utils.HanldeError(w, err)
		return
	}

	user, err := h.service.Register(r.Context(), req)
	if err != nil {
		utils.HanldeError(w, err)
		return
	}
	utils.JSON(w, http.StatusCreated, "User registered successfully", user)
}

func (h *Handler) Login(w http.ResponseWriter, r *http.Request) {
	var req LoginRequest
	if err := utils.Decode(r, &req); err != nil {
		utils.Error(w, http.StatusBadRequest, "Failed to decode request body")
		return
	}

	if err := utils.ValidationStruct(req); err != nil {
		utils.HanldeError(w, err)
		return
	}

	user, err := h.service.Login(r.Context(), req)
	if err != nil {
		utils.HanldeError(w, err)
		return
	}
	utils.JSON(w, http.StatusOK, "User logged in successfully", user)
}

func (h *Handler) RefreshToken(w http.ResponseWriter, r *http.Request) {
	var req RefreshTokenRequest
	if err := utils.Decode(r, &req); err != nil {
		utils.Error(w, http.StatusBadRequest, "Failed to decode request body")
		return
	}

	if err := utils.ValidationStruct(req); err != nil {
		utils.HanldeError(w, err)
		return
	}

	user, err := h.service.RefreshTokens(r.Context(), req.RefreshToken)
	if err != nil {
		utils.HanldeError(w, err)
		return
	}
	utils.JSON(w, http.StatusOK, "Token refreshed successfully", user)
}

func (h *Handler) Logout(w http.ResponseWriter, r *http.Request) {
	var req RefreshTokenRequest
	if err := utils.Decode(r, &req); err != nil {
		utils.Error(w, http.StatusBadRequest, "Failed to decode request body")
		return
	}

	if err := utils.ValidationStruct(req); err != nil {
		utils.HanldeError(w, err)
		return
	}

	h.service.Logout(r.Context(), req.RefreshToken)
	utils.JSON(w, http.StatusOK, "User logged out successfully", nil)
}

func (h *Handler) ForgotPassword(w http.ResponseWriter, r *http.Request) {
	var req ForgotPasswordRequest
	if err := utils.Decode(r, &req); err != nil {
		utils.Error(w, http.StatusBadRequest, "Failed to decode request body")
		return
	}

	if err := utils.ValidationStruct(req); err != nil {
		utils.HanldeError(w, err)
		return
	}

	if err := h.service.ForgotPassword(r.Context(), req.Email); err != nil {
		utils.HanldeError(w, err)
		return
	}
	utils.JSON(w, http.StatusOK, "Password reset code sent successfully", nil)
}

func (h *Handler) ResetPassword(w http.ResponseWriter, r *http.Request) {
	var req ResetPasswordRequest

	if err := utils.Decode(r, &req); err != nil {
		utils.Error(
			w,
			http.StatusBadRequest,
			"Failed to decode request body",
		)
		return
	}

	if err := utils.ValidationStruct(req); err != nil {
		utils.HanldeError(w, err)
		return
	}

	if err := h.service.ResetPassword(
		r.Context(),
		req.Email,
		req.Code,
		req.Password,
	); err != nil {
		utils.HanldeError(w, err)
		return
	}

	utils.JSON(
		w,
		http.StatusOK,
		"Password reset successfully",
		nil,
	)
}
