package middleware

import (
	"context"
	"net/http"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
)

type contextKey string

const UserIDKey contextKey = "userID"

// AuthMiddleware memeriksa apakah request memiliki cookie session yang valid
func AuthMiddleware(pool *pgxpool.Pool) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			cookie, err := r.Cookie("session_id")
			if err != nil || cookie.Value == "" {
				http.Error(w, `{"error": "Unauthorized: Sesi tidak ditemukan, silakan login"}`, http.StatusUnauthorized)
				return
			}

			var userID int
			var expiresAt time.Time

			query := `SELECT user_id, expires_at FROM sessions WHERE id = $1`
			err = pool.QueryRow(r.Context(), query, cookie.Value).Scan(&userID, &expiresAt)
			if err != nil {
				// Sesi tidak ditemukan di database
				http.Error(w, `{"error": "Unauthorized: Sesi tidak valid atau telah berakhir"}`, http.StatusUnauthorized)
				return
			}

			// Cek apakah masa berlaku sesi telah habis
			if time.Now().After(expiresAt) {
				// Hapus sesi yang kedaluwarsa dari database
				_, _ = pool.Exec(r.Context(), "DELETE FROM sessions WHERE id = $1", cookie.Value)
				http.Error(w, `{"error": "Unauthorized: Sesi telah kedaluwarsa"}`, http.StatusUnauthorized)
				return
			}

			// Simpan user_id ke dalam context request
			ctx := context.WithValue(r.Context(), UserIDKey, userID)
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}

// GetUserID mengambil user_id dari context request
func GetUserID(ctx context.Context) (int, bool) {
	val := ctx.Value(UserIDKey)
	if val == nil {
		return 0, false
	}
	id, ok := val.(int)
	return id, ok
}

