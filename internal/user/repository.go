package user

import (
	"context"
	"fmt"
	"strings"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/saurav11sarkar/001practic/internal/utils"
)

type Repository struct {
	db *pgxpool.Pool
}

func NewRepository(db *pgxpool.Pool) *Repository {
	return &Repository{db: db}
}

func (r *Repository) GetProfile(ctx context.Context, userID string) (*User, error) {
	var user User
	err := r.db.QueryRow(
		ctx,
		`SELECT id, name, email, password, role, status, COALESCE(profile, ''), created_at, updated_at FROM users WHERE id = $1`,
		userID,
	).Scan(
		&user.ID,
		&user.Name,
		&user.Email,
		&user.Password,
		&user.Role,
		&user.Status,
		&user.Profile,
		&user.CreatedAt,
		&user.UpdatedAt,
	)
	if err != nil {
		return nil, err
	}
	return &user, nil
}

func (r *Repository) GetAllUser(ctx context.Context, q utils.Query) ([]User, error) {
	query := `SELECT id, name, email, password, role, status,COALESCE(profile, '') AS profile, created_at, updated_at FROM users WHERE 1=1`
	args := pgx.NamedArgs{}

	if q.Search != "" {
		searchableColumns := []string{"name", "email", "role"}
		var searchConditions []string
		for _, col := range searchableColumns {
			searchConditions = append(searchConditions, fmt.Sprintf("%s ILIKE @search", col))
		}

		// 2. Add spaces on both sides of " OR "
		query += fmt.Sprintf(` AND (%s) `, strings.Join(searchConditions, " OR "))
		args["search"] = "%" + q.Search + "%"
	}

	allowFilters := map[string]bool{
		"role":   true,
		"status": true,
		"email":  true,
		"name":   true,
	}

	for key, val := range q.Filters {
		if allowFilters[key] {
			query += fmt.Sprintf(` AND %s = @%s`, key, key)
			args[key] = val
		}
	}

	allowedSortCols := map[string]bool{
		"id": true, "name": true, "email": true, "role": true, "status": true, "created_at": true, "updated_at": true,
	}

	sortBy := "created_at"
	if allowedSortCols[q.SortBy] {
		sortBy = q.SortBy
	}

	sortOrder := "DESC"
	if strings.ToLower(q.SortOrder) == "asc" {
		sortOrder = "ASC"
	}

	query += fmt.Sprintf(` ORDER BY %s %s LIMIT @limit OFFSET @offset`, sortBy, sortOrder)
	args["limit"] = q.Limit
	args["offset"] = q.Offset()

	rows, err := r.db.Query(ctx, query, args)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	users, err := pgx.CollectRows(rows, pgx.RowToStructByNameLax[User])
	if err != nil {
		return nil, err
	}
	if users == nil {
		return []User{}, nil
	}
	return users, nil
}

func (r *Repository) GetSingleUser(ctx context.Context, userID string) (*User, error) {
	var user User
	err := r.db.QueryRow(
		ctx,
		`SELECT id, name, email, password, role, status,COALESCE(profile, '') AS profile, created_at, updated_at FROM users WHERE id = $1`,
		userID,
	).Scan(
		&user.ID,
		&user.Name,
		&user.Email,
		&user.Password,
		&user.Role,
		&user.Status,
		&user.Profile,
		&user.CreatedAt,
		&user.UpdatedAt,
	)
	if err != nil {
		return nil, err
	}
	return &user, nil
}

func (r *Repository) UpdateProfile(ctx context.Context, userID string, profile string) error {
	_, err := r.db.Exec(ctx, `UPDATE users SET profile = $1 WHERE id = $2`, profile, userID)
	if err != nil {
		return err
	}
	return nil
}
