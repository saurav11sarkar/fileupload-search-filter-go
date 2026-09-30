package user

import (
	"context"
	"errors"
	"mime/multipart"
	"net/http"

	"github.com/jackc/pgx/v5"
	"github.com/saurav11sarkar/001practic/internal/utils"
)

type Service struct {
	repo     *Repository
	uploader *utils.CloudinaryService
}

func NewService(repo *Repository, uploader *utils.CloudinaryService) *Service {
	return &Service{repo: repo, uploader: uploader}
}

func (s *Service) GetProfile(ctx context.Context, userID string) (*User, error) {
	user, err := s.repo.GetProfile(ctx, userID)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, utils.NewAppError(http.StatusNotFound, "NOT_FOUND", "User not found")
		}
		return nil, err
	}
	return user, nil
}

func (s *Service) GetAllUser(ctx context.Context, q utils.Query) ([]User, error) {
	users, err := s.repo.GetAllUser(ctx, q)
	if err != nil {
		return nil, utils.NewAppError(http.StatusInternalServerError, "INTERNAL_ERROR", "Failed to fetch users")
	}
	return users, nil
}

func (s *Service) GetSingleUser(ctx context.Context, userID string) (*User, error) {
	user, err := s.repo.GetSingleUser(ctx, userID)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, utils.NewAppError(http.StatusNotFound, "NOT_FOUND", "User not found")
		}
		return nil, err
	}
	return user, nil
}

func (s *Service) UpdateProfile(ctx context.Context, userID string, profile multipart.File) error {

	imageUrl, err := s.uploader.UploadFile(ctx, profile, "user_profiles")
	if err != nil {
		return utils.NewAppError(http.StatusInternalServerError, "INTERNAL_ERROR", "Failed to upload profile image")
	}

	err = s.repo.UpdateProfile(ctx, userID, imageUrl)
	if err != nil {
		return utils.NewAppError(http.StatusInternalServerError, "INTERNAL_ERROR", "Failed to update user profile")
	}

	return nil
}

func (s *Service) UpdateUser(ctx context.Context, userID string, updateUserDto UpdateUserRequest) (*User, error) {
	user, err := s.repo.UpdateUser(ctx, userID, updateUserDto)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, utils.NewAppError(http.StatusNotFound, "NOT_FOUND", "User not found")
		}
		return nil, err
	}
	return user, nil
}


func (s *Service) DeleteUser(ctx context.Context, userID string) error {
	err := s.repo.DeleteUser(ctx, userID)
	if err != nil {
		return utils.NewAppError(http.StatusInternalServerError, "INTERNAL_ERROR", "Failed to delete user")
	}
	return nil
}
