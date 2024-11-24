package storage

import (
	"context"
	"fmt"
	"podbor/internal/config"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"go.uber.org/zap"
)

type Storage struct {
	pool   *pgxpool.Pool
	logger *zap.Logger
}

func New(cfg *config.Config, logger *zap.Logger) (*Storage, error) {
	dbCfg := cfg.Database
	dsn := fmt.Sprintf(
		"host=%s port=%d user=%s password=%s dbname=%s sslmode=%s",
		dbCfg.Host,
		dbCfg.Port,
		dbCfg.User,
		dbCfg.Password,
		dbCfg.DBName,
		dbCfg.SSLMode,
	)

	pool, err := pgxpool.New(context.Background(), dsn)
	if err != nil {
		return nil, err
	}

	storage := &Storage{
		pool:   pool,
		logger: logger,
	}

	return storage, nil
}

func (s *Storage) Close() {
	if s.pool != nil {
		s.pool.Close()
	}
}

func (s *Storage) GetUserByID(ctx context.Context, userID string) (*User, error) {
	var user User
	query := `SELECT id, username, email, created_at FROM users WHERE id = $1`

	err := s.pool.QueryRow(ctx, query, userID).Scan(
		&user.ID,
		&user.Username,
		&user.Email,
		&user.CreatedAt,
	)
	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, nil
		}
		s.logger.Sugar().Error("Не удалось получить пользователя", zap.Error(err))
		return nil, err
	}

	return &user, nil
}
