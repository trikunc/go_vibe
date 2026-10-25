package database

import (
	"context"
	"fmt"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
)

// Connect membuat koneksi pool ke database PostgreSQL menggunakan pgx
func Connect(databaseURL string) (*pgxpool.Pool, error) {
	if databaseURL == "" {
		return nil, fmt.Errorf("DATABASE_URL tidak boleh kosong")
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	config, err := pgxpool.ParseConfig(databaseURL)
	if err != nil {
		return nil, fmt.Errorf("format DATABASE_URL tidak valid: %w", err)
	}

	pool, err := pgxpool.NewWithConfig(ctx, config)
	if err != nil {
		return nil, fmt.Errorf("gagal membuat koneksi database: %w", err)
	}

	// Memverifikasi apakah database dapat dihubungi
	if err := pool.Ping(ctx); err != nil {
		pool.Close()
		return nil, fmt.Errorf("tidak dapat menghubungi database PostgreSQL: %w", err)
	}

	return pool, nil
}

// Migrate membuat tabel-tabel yang dibutuhkan jika belum ada
func Migrate(ctx context.Context, pool *pgxpool.Pool, schemaSQL string) error {
	_, err := pool.Exec(ctx, schemaSQL)
	if err != nil {
		return fmt.Errorf("gagal menjalankan migrasi schema: %w", err)
	}
	return nil
}

