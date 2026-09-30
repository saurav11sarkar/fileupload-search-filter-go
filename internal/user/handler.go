package user

import (
	"net/http"

	"github.com/google/uuid"
	"github.com/saurav11sarkar/001practic/internal/utils"
)

type Handler struct {
	service *Service
}

func NewHandler(service *Service) *Handler {
	return &Handler{service: service}
}

func (h *Handler) GetProfile(w http.ResponseWriter, r *http.Request) {
	claims, ok := r.Context().Value("user").(*utils.Claim)
	if !ok || claims == nil {
		utils.HanldeError(w, utils.NewAppError(http.StatusUnauthorized, "UNAUTHORIZED", "Unauthorized access"))
		return
	}
	user, err := h.service.GetProfile(r.Context(), claims.UserID)
	if err != nil {
		utils.HanldeError(w, err)
		return
	}
	utils.JSON(w, http.StatusOK, "User profile retrieved successfully", user)
}

func (h *Handler) GetAllUser(w http.ResponseWriter, r *http.Request) {
	q := utils.ParseQuery(r.URL.Query())
	users, err := h.service.GetAllUser(r.Context(), q)
	if err != nil {
		utils.HanldeError(w, err)
		return
	}
	responseData := map[string]interface{}{
		"user": users,
		"meta": map[string]interface{}{
			"page":      q.Page,
			"limit":     q.Limit,
			"total":     len(users),
			"totalPage": (len(users) + q.Limit - 1) / q.Limit,
		},
	}

	utils.JSON(w, http.StatusOK, "Users retrieved successfully", responseData)
}

func (h *Handler) GetSingleUser(w http.ResponseWriter, r *http.Request) {
	idPath := r.PathValue("id")

	if err := uuid.Validate(idPath); err != nil {
		utils.HanldeError(w, utils.NewAppError(http.StatusBadRequest, "BAD_REQUEST", "Invalid user ID"))
		return
	}

	user, err := h.service.GetSingleUser(r.Context(), idPath)
	if err != nil {
		utils.HanldeError(w, err)
		return
	}

	utils.JSON(w, http.StatusOK, "User retrieved successfully", user)
}

func (h *Handler) UpdateProfile(w http.ResponseWriter, r *http.Request) {
	claims, ok := r.Context().Value("user").(*utils.Claim)
	if !ok || claims == nil {
		utils.HanldeError(w, utils.NewAppError(http.StatusUnauthorized, "UNAUTHORIZED", "Unauthorized access"))
		return
	}

	r.Body = http.MaxBytesReader(w, r.Body, 5<<20)
	if err := r.ParseMultipartForm(5 << 20); err != nil {
		utils.HanldeError(w, utils.NewAppError(http.StatusBadRequest, "FILE_TOO_LARGE", "File exceeds maximum size of 5MB or invalid form data"))
		return
	}

	file, fileHandler, err := r.FormFile("profile")
	if err != nil {
		utils.HanldeError(w, utils.NewAppError(http.StatusBadRequest, "BAD_REQUEST", "Profile image file is required in 'profile_image' field"))
		return
	}
	defer file.Close()

	contentType := fileHandler.Header.Get("Content-Type")
	if contentType != "image/jpeg" && contentType != "image/png" {
		utils.HanldeError(w, utils.NewAppError(http.StatusBadRequest, "BAD_REQUEST", "Profile image must be in JPEG or PNG format"))
		return
	}

	err = h.service.UpdateProfile(r.Context(), claims.UserID, file)
	if err != nil {
		utils.HanldeError(w, err)
		return
	}
	utils.JSON(w, http.StatusOK, "User profile updated successfully", nil)
}

func (h *Handler) UpdateUser(w http.ResponseWriter, r *http.Request) {
	cliams, ok := r.Context().Value("user").(*utils.Claim)
	if !ok || cliams == nil {
		utils.HanldeError(w, utils.NewAppError(http.StatusUnauthorized, "UNAUTHORIZED", "Unauthorized access"))
		return
	}

	var req UpdateUserRequest
	if err := utils.Decode(r, &req); err != nil {
		utils.Error(w, http.StatusBadRequest, "Failed to decode request body")
		return
	}

	if err := utils.ValidationStruct(req); err != nil {
		utils.HanldeError(w, err)
		return
	}

	user, err := h.service.UpdateUser(r.Context(), cliams.UserID, req)
	if err != nil {
		utils.HanldeError(w, err)
		return
	}
	utils.JSON(w, http.StatusOK, "User updated successfully", user)
}

func (h *Handler) DeleteUser(w http.ResponseWriter, r *http.Request) {
	idPath := r.PathValue("id")

	if err := uuid.Validate(idPath); err != nil {
		utils.HanldeError(w, utils.NewAppError(http.StatusBadRequest, "BAD_REQUEST", "Invalid user ID"))
		return
	}
	h.service.DeleteUser(r.Context(), idPath)
	utils.JSON(w, http.StatusOK, "User deleted successfully", nil)
}
