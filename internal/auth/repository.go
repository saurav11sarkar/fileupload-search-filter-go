package auth

import (
	"context"
	"errors"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/saurav11sarkar/001practic/internal/utils"
)

type Repository struct {
	db *pgxpool.Pool
}

func NewRepository(db *pgxpool.Pool) *Repository {
	return &Repository{db: db}
}

func (r *Repository) Register(ctx context.Context, dto RegisterRequest) (*UserResponse, error) {
	var user UserResponse
	err := r.db.QueryRow(
		ctx,
		`INSERT INTO users (
			id,
			name,
			email,
			password,
			role,
			status
		) VALUES ($1, $2, $3, $4, $5, $6)
		RETURNING id, name, email, role, status, created_at, updated_at`,
		utils.NewID(),
		dto.Name,
		dto.Email,
		dto.Password,
		"user",
		"active",
	).Scan(&user.ID, &user.Name, &user.Email, &user.Role, &user.Status, &user.CreatedAt, &user.UpdatedAt)

	if err != nil {
		return nil, err
	}

	return &user, nil
}

func (r *Repository) GetEmail(ctx context.Context, email string) (*UserResponse, error) {
	var user UserResponse
	err := r.db.QueryRow(
		ctx,
		`SELECT id,name,email,password,role,status FROM users WHERE email = $1`,
		email,
	).Scan(&user.ID, &user.Name, &user.Email, &user.Password, &user.Role, &user.Status)
	if err != nil {
		return nil, err
	}
	return &user, nil
}

func (r *Repository) SaveRefreshToken(
	ctx context.Context,
	dto SaveRefreshTokenRequest,
) error {

	_, err := r.db.Exec(
		ctx,
		`INSERT INTO refresh_tokens (
			id,
			user_id,
			refresh_token,
			expires_at
		) VALUES ($1, $2, $3, $4)`,
		utils.NewID(),
		dto.UserID,
		dto.RefreshToken,
		dto.ExpiresAt,
	)

	return err
}

func (r *Repository) FindRefreshToken(
	ctx context.Context,
	userID string,
) ([]RefreshTokenRecode, error) {
	rows, err := r.db.Query(
		ctx,
		`SELECT id, refresh_token
		 FROM refresh_tokens
		 WHERE user_id = $1
		 AND expires_at > now()`,
		userID,
	)
	if err != nil {
		return nil, err
	}

	defer rows.Close()

	var items []RefreshTokenRecode

	for rows.Next() {
		var item RefreshTokenRecode
		if err := rows.Scan(&item.ID, &item.RefreshToken); err != nil {
			return nil, err
		}
		items = append(items, item)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	if len(items) == 0 {
		return nil, nil
	}

	return items, nil
}

func (r *Repository) DeleteRefreshToken(ctx context.Context, id string) error {
	_, err := r.db.Exec(ctx, `DELETE FROM refresh_tokens WHERE id = $1`, id)
	return err
}

func (r *Repository) DeleteUserrefreshToken(ctx context.Context, userID string) error {
	_, err := r.db.Exec(ctx, `DELETE FROM refresh_tokens WHERE user_id = $1`, userID)
	return err
}

func (r *Repository) SaveResetCode(
	ctx context.Context,
	userID string,
	code string,
	expiresAt time.Time,
) error {
	_, err := r.db.Exec(
		ctx,
		`UPDATE users
		 SET reset_code = $1,
		     reset_code_expires_at = $2,
		     updated_at = now()
		 WHERE id = $3`,
		code,
		expiresAt,
		userID,
	)

	return err
}

func (r *Repository) FindUserResetCode(
	ctx context.Context,
	email string,
	code string,
) (string, error) {

	var userID string

	err := r.db.QueryRow(
		ctx,
		`SELECT id
		 FROM users
		 WHERE email = $1
		   AND reset_code = $2
		   AND reset_code_expires_at > now()`,
		email,
		code,
	).Scan(&userID)

	return userID, err
}

func (r *Repository) UpdatePassword(
	ctx context.Context,
	userID string,
	hashedPassword string,
) error {

	result, err := r.db.Exec(
		ctx,
		`UPDATE users
         SET password = $1,
             reset_code = NULL,
             reset_code_expires_at = NULL,
             updated_at = now()
         WHERE id = $2`,
		hashedPassword,
		userID,
	)

	if err != nil {
		return err
	}

	if result.RowsAffected() == 0 {
		return errors.New("user not found")
	}

	return nil
}
