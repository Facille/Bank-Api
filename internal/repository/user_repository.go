package repository

import (
	"context"
	"errors"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/therealadik/bank-api/internal/models"
)

// ErrUserNotFound возвращается, когда пользователь не найден
var ErrUserNotFound = errors.New("пользователь не найден")

// UserRepository интерфейс для работы с пользователями
type UserRepository interface {
	Create(ctx context.Context, user *models.User) (int64, error)
	GetByEmail(ctx context.Context, email string) (*models.User, error)
	GetByID(ctx context.Context, id int64) (*models.User, error)
}

// UserRepositoryPgx реализация репозитория пользователей с использованием pgx
type UserRepositoryPgx struct {
	pool *pgxpool.Pool
}

// NewUserRepository создает новый репозиторий пользователей
func NewUserRepository(pool *pgxpool.Pool) UserRepository {
	return &UserRepositoryPgx{pool: pool}
}

// Create создает нового пользователя
func (r *UserRepositoryPgx) Create(ctx context.Context, user *models.User) (int64, error) {
	var id int64

	err := r.pool.QueryRow(ctx,
		`INSERT INTO users (email, password_hash) 
         VALUES ($1, $2) 
         RETURNING id`,
		user.Email, user.Password).Scan(&id)

	if err != nil {
		return 0, err
	}

	return id, nil
}

// GetByEmail находит пользователя по email
func (r *UserRepositoryPgx) GetByEmail(ctx context.Context, email string) (*models.User, error) {
	user := &models.User{}

	err := r.pool.QueryRow(ctx,
		`SELECT id, email, password_hash, created_at 
         FROM users 
         WHERE email = $1`,
		email).Scan(&user.ID, &user.Email, &user.Password, &user.CreatedAt)

	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, ErrUserNotFound
		}
		return nil, err
	}

	return user, nil
}

// GetByID находит пользователя по ID
func (r *UserRepositoryPgx) GetByID(ctx context.Context, id int64) (*models.User, error) {
	user := &models.User{}

	err := r.pool.QueryRow(ctx,
		`SELECT id, email, password_hash, created_at 
         FROM users 
         WHERE id = $1`,
		id).Scan(&user.ID, &user.Email, &user.Password, &user.CreatedAt)

	if err != nil {
		return nil, err
	}

	return user, nil
}
