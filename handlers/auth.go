package handlers

import (
	"crypto/rand"
	"encoding/hex"
	"encoding/json"
	"errors"
	"net/http"
	"strings"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgxpool"
	"golang.org/x/crypto/bcrypt"

	"vibecode/middleware"
)

type AuthHandler struct {
	pool *pgxpool.Pool
}

func NewAuthHandler(pool *pgxpool.Pool) *AuthHandler {
	return &AuthHandler{pool: pool}
}

type RegisterRequest struct {
	Name     string `json:"name"`
	Email    string `json:"email"`
	Password string `json:"password"`
}

type LoginRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type UserResponse struct {
	ID    int    `json:"id"`
	Name  string `json:"name"`
	Email string `json:"email"`
}

// Register menangani pendaftaran pengguna baru
func (h *AuthHandler) Register(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		sendError(w, http.StatusMethodNotAllowed, "Metode HTTP tidak diizinkan")
		return
	}

	var req RegisterRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		sendError(w, http.StatusBadRequest, "Format JSON tidak valid")
		return
	}

	req.Name = strings.TrimSpace(req.Name)
	req.Email = strings.TrimSpace(strings.ToLower(req.Email))

	// Validasi input sederhana
	if req.Name == "" {
		sendError(w, http.StatusBadRequest, "Nama tidak boleh kosong")
		return
	}
	if req.Email == "" || !strings.Contains(req.Email, "@") {
		sendError(w, http.StatusBadRequest, "Alamat email tidak valid")
		return
	}
	if len(req.Password) < 6 {
		sendError(w, http.StatusBadRequest, "Password minimal 6 karakter")
		return
	}

	// Hash password menggunakan bcrypt
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		sendError(w, http.StatusInternalServerError, "Gagal memproses password")
		return
	}

	// Simpan ke database
	var user UserResponse
	query := `INSERT INTO users (name, email, password_hash) VALUES ($1, $2, $3) RETURNING id, name, email`
	err = h.pool.QueryRow(r.Context(), query, req.Name, req.Email, string(hashedPassword)).Scan(&user.ID, &user.Name, &user.Email)
	if err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) && pgErr.Code == "23505" { // Unique violation
			sendError(w, http.StatusConflict, "Email sudah terdaftar, silakan gunakan email lain atau login")
			return
		}
		sendError(w, http.StatusInternalServerError, "Gagal mendaftarkan pengguna")
		return
	}

	sendJSON(w, http.StatusCreated, map[string]any{
		"message": "Registrasi berhasil",
		"user":    user,
	})
}

// Login memverifikasi email & password, lalu membuat sesi baru
func (h *AuthHandler) Login(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		sendError(w, http.StatusMethodNotAllowed, "Metode HTTP tidak diizinkan")
		return
	}

	var req LoginRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		sendError(w, http.StatusBadRequest, "Format JSON tidak valid")
		return
	}

	req.Email = strings.TrimSpace(strings.ToLower(req.Email))

	if req.Email == "" || req.Password == "" {
		sendError(w, http.StatusBadRequest, "Email dan password wajib diisi")
		return
	}

	var user UserResponse
	var passwordHash string

	query := `SELECT id, name, email, password_hash FROM users WHERE email = $1`
	err := h.pool.QueryRow(r.Context(), query, req.Email).Scan(&user.ID, &user.Name, &user.Email, &passwordHash)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			sendError(w, http.StatusUnauthorized, "Email atau password salah")
			return
		}
		sendError(w, http.StatusInternalServerError, "Terjadi kesalahan saat memeriksa akun")
		return
	}

	// Verifikasi hash password
	if err := bcrypt.CompareHashAndPassword([]byte(passwordHash), []byte(req.Password)); err != nil {
		sendError(w, http.StatusUnauthorized, "Email atau password salah")
		return
	}

	// Buat ID sesi acak yang aman (32 byte diubah ke string hex 64 karakter)
	tokenBytes := make([]byte, 32)
	if _, err := rand.Read(tokenBytes); err != nil {
		sendError(w, http.StatusInternalServerError, "Gagal membuat sesi login")
		return
	}
	sessionID := hex.EncodeToString(tokenBytes)
	expiresAt := time.Now().Add(24 * time.Hour) // Sesi berlaku selama 24 jam

	// Simpan sesi ke database
	sessionQuery := `INSERT INTO sessions (id, user_id, expires_at) VALUES ($1, $2, $3)`
	_, err = h.pool.Exec(r.Context(), sessionQuery, sessionID, user.ID, expiresAt)
	if err != nil {
		sendError(w, http.StatusInternalServerError, "Gagal menyimpan data sesi")
		return
	}

	// Pasang cookie HTTP ke respon klien
	http.SetCookie(w, &http.Cookie{
		Name:     "session_id",
		Value:    sessionID,
		Path:     "/",
		Expires:  expiresAt,
		HttpOnly: true,
		SameSite: http.SameSiteLaxMode,
	})

	sendJSON(w, http.StatusOK, map[string]any{
		"message": "Login berhasil",
		"user":    user,
	})
}

// Logout menghapus sesi dari database dan membersihkan cookie
func (h *AuthHandler) Logout(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		sendError(w, http.StatusMethodNotAllowed, "Metode HTTP tidak diizinkan")
		return
	}

	cookie, err := r.Cookie("session_id")
	if err == nil && cookie.Value != "" {
		// Hapus sesi dari database
		_, _ = h.pool.Exec(r.Context(), "DELETE FROM sessions WHERE id = $1", cookie.Value)
	}

	// Hapus cookie di browser dengan mengatur MaxAge negatif
	http.SetCookie(w, &http.Cookie{
		Name:     "session_id",
		Value:    "",
		Path:     "/",
		MaxAge:   -1,
		Expires:  time.Unix(0, 0),
		HttpOnly: true,
		SameSite: http.SameSiteLaxMode,
	})

	sendJSON(w, http.StatusOK, map[string]any{
		"message": "Logout berhasil",
	})
}

// Me mengembalikan data profil pengguna yang sedang login
func (h *AuthHandler) Me(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		sendError(w, http.StatusMethodNotAllowed, "Metode HTTP tidak diizinkan")
		return
	}

	userID, ok := middleware.GetUserID(r.Context())
	if !ok {
		sendError(w, http.StatusUnauthorized, "Sesi tidak ditemukan")
		return
	}

	var user UserResponse
	query := `SELECT id, name, email FROM users WHERE id = $1`
	err := h.pool.QueryRow(r.Context(), query, userID).Scan(&user.ID, &user.Name, &user.Email)
	if err != nil {
		sendError(w, http.StatusNotFound, "Pengguna tidak ditemukan")
		return
	}

	sendJSON(w, http.StatusOK, user)
}

// Helper untuk mengirim respon JSON
func sendJSON(w http.ResponseWriter, status int, data any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(data)
}

// Helper untuk mengirim pesan error JSON
func sendError(w http.ResponseWriter, status int, message string) {
	sendJSON(w, status, map[string]string{
		"error": message,
	})
}

