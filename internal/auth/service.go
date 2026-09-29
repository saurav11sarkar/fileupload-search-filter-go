package auth

import (
	"context"
	"crypto/rand"
	"errors"
	"fmt"
	"math/big"
	"net/http"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/saurav11sarkar/001practic/internal/config"
	"github.com/saurav11sarkar/001practic/internal/email"
	"github.com/saurav11sarkar/001practic/internal/utils"
)

type Service struct {
	repo  *Repository
	cfg   config.Config
	email *email.Email
}

func NewService(repo *Repository, cfg config.Config, email *email.Email) *Service {
	return &Service{repo: repo, cfg: cfg, email: email}
}

func (s *Service) Register(ctx context.Context, create RegisterRequest) (*UserResponse, error) {
	_, err := s.repo.GetEmail(ctx, create.Email)
	if err == nil {
		return nil, utils.NewAppError(http.StatusConflict, "EMAIL_ALREADY_EXISTS", "Email already exists")
	}
	if !errors.Is(err, pgx.ErrNoRows) {
		return nil, err
	}

	hashedPassword, err := utils.HashPassword(create.Password)
	if err != nil {
		return nil, err
	}
	create.Password = hashedPassword

	user, err := s.repo.Register(ctx, create)
	if err != nil {
		return nil, err
	}
	return user, nil
}

func (s *Service) Login(ctx context.Context, login LoginRequest) (*LoginResponse, error) {
	user, err := s.repo.GetEmail(ctx, login.Email)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, utils.NewAppError(http.StatusUnauthorized, "INVALID_CREDENTIALS", "Invalid credentials")
		}
		return nil, err
	}

	if err := utils.ComparePassword(login.Password, user.Password); err != nil {
		return nil, utils.NewAppError(http.StatusUnauthorized, "INVALID_CREDENTIALS", "Invalid credentials")
	}

	accessToken, err := utils.CreateToken(user.ID, user.Role, "access", s.cfg.Auth.JWTAccessSecret, time.Duration(s.cfg.Auth.AccessTokenMinutes)*time.Minute)
	if err != nil {
		return nil, utils.NewAppError(http.StatusInternalServerError, "TOKEN_GENERATION_FAILED", "Failed to generate access token")
	}
	refreshToken, err := utils.CreateToken(user.ID, user.Role, "refresh", s.cfg.Auth.JWTRefreshSecret, time.Duration(s.cfg.Auth.RefreshTokenDays)*24*time.Hour)
	if err != nil {
		return nil, utils.NewAppError(http.StatusInternalServerError, "TOKEN_GENERATION_FAILED", "Failed to generate refresh token")
	}

	err = s.repo.SaveRefreshToken(ctx, SaveRefreshTokenRequest{
		UserID:       user.ID,
		RefreshToken: refreshToken,
		ExpiresAt:    time.Now().Add(time.Duration(s.cfg.Auth.RefreshTokenDays) * 24 * time.Hour),
	})
	if err != nil {
		return nil, err
	}

	return &LoginResponse{
		User:         *user,
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
	}, nil
}

func (s *Service) RefreshTokens(ctx context.Context, token string) (map[string]any, error) {
	claims, err := utils.ParseToken(token, s.cfg.Auth.JWTRefreshSecret, "refresh")
	if err != nil {
		return nil, utils.NewAppError(http.StatusUnauthorized, "INVALID_REFRESH_TOKEN", "Invalid refresh token")
	}
	records, err := s.repo.FindRefreshToken(ctx, claims.UserID)
	if err != nil {
		return nil, utils.NewAppError(http.StatusUnauthorized, "INVALID_REFRESH_TOKEN", "Invalid refresh token")
	}

	var tokenID string
	for _, record := range records {
		if record.RefreshToken == token {
			tokenID = record.ID
			break
		}
	}
	if tokenID == "" {
		return nil, utils.NewAppError(http.StatusUnauthorized, "INVALID_REFRESH_TOKEN", "Invalid refresh token")
	}

	if err := s.repo.DeleteRefreshToken(ctx, tokenID); err != nil {
		return nil, utils.NewAppError(http.StatusInternalServerError, "TOKEN_GENERATION_FAILED", "Failed to generate access token")
	}

	accessToken, err := utils.CreateToken(
		claims.UserID,
		claims.Role,
		"access",
		s.cfg.Auth.JWTAccessSecret,
		time.Duration(s.cfg.Auth.AccessTokenMinutes)*time.Minute,
	)
	if err != nil {
		return nil, utils.NewAppError(
			http.StatusInternalServerError,
			"TOKEN_GENERATION_FAILED",
			"Failed to generate access token",
		)
	}

	newRefreshToken, err := utils.CreateToken(
		claims.UserID,
		claims.Role,
		"refresh",
		s.cfg.Auth.JWTRefreshSecret,
		time.Duration(s.cfg.Auth.RefreshTokenDays)*24*time.Hour,
	)
	if err != nil {
		return nil, utils.NewAppError(
			http.StatusInternalServerError,
			"TOKEN_GENERATION_FAILED",
			"Failed to generate refresh token",
		)
	}

	if err := s.repo.SaveRefreshToken(ctx, SaveRefreshTokenRequest{
		UserID:       claims.UserID,
		RefreshToken: newRefreshToken,
		ExpiresAt:    time.Now().Add(time.Duration(s.cfg.Auth.RefreshTokenDays) * 24 * time.Hour),
	}); err != nil {
		return nil, err
	}

	return map[string]any{
		"access_token":  accessToken,
		"refresh_token": newRefreshToken,
		"user": map[string]any{
			"id":   claims.UserID,
			"role": claims.Role,
		},
	}, nil
}

func (s *Service) Logout(ctx context.Context, token string) {
	claims, err := utils.ParseToken(token, s.cfg.Auth.JWTRefreshSecret, "refresh")
	if err != nil {
		return
	}
	_ = s.repo.DeleteUserrefreshToken(ctx, claims.UserID)
}

func (s *Service) ForgotPassword(ctx context.Context, email string) error {
	user, err := s.repo.GetEmail(ctx, email)
	if err != nil {
		return utils.NewAppError(http.StatusNotFound, "NOT_FOUND", "Email not found")
	}

	n, err := rand.Int(rand.Reader, big.NewInt(1000000))
	if err != nil {
		return err
	}
	code := n.Int64()
	err = s.repo.SaveResetCode(ctx, user.ID, fmt.Sprintf("%d", code), time.Now().Add(time.Minute*10))
	if err != nil {
		return err
	}

	err = s.email.SendEmail(email, "Reset Password", fmt.Sprintf("Your reset code is %d", code))
	if err != nil {
		return err
	}
	return nil
}

func (s *Service) ResetPassword(
	ctx context.Context,
	email string,
	code string,
	password string,
) error {

	userID, err := s.repo.FindUserResetCode(ctx, email, code)
	if err != nil {
		return utils.NewAppError(
			http.StatusBadRequest,
			"INVALID_RESET_CODE",
			"Invalid or expired reset code",
		)
	}

	hashedPassword, err := utils.HashPassword(password)
	if err != nil {
		return err
	}

	if err := s.repo.UpdatePassword(ctx, userID, hashedPassword); err != nil {
		return err
	}

	return nil
}
