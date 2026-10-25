package handlers

import (
	"encoding/json"
	"errors"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"

	"vibecode/middleware"
)

type NotesHandler struct {
	pool *pgxpool.Pool
}

func NewNotesHandler(pool *pgxpool.Pool) *NotesHandler {
	return &NotesHandler{pool: pool}
}

type Note struct {
	ID        int       `json:"id"`
	UserID    int       `json:"user_id"`
	Title     string    `json:"title"`
	Content   string    `json:"content"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type NoteRequest struct {
	Title   string `json:"title"`
	Content string `json:"content"`
}

// List mengambil semua catatan milik pengguna yang sedang login
func (h *NotesHandler) List(w http.ResponseWriter, r *http.Request) {
	userID, ok := middleware.GetUserID(r.Context())
	if !ok {
		sendError(w, http.StatusUnauthorized, "Sesi tidak ditemukan")
		return
	}

	query := `
		SELECT id, user_id, title, content, created_at, updated_at 
		FROM notes 
		WHERE user_id = $1 
		ORDER BY updated_at DESC
	`
	rows, err := h.pool.Query(r.Context(), query, userID)
	if err != nil {
		sendError(w, http.StatusInternalServerError, "Gagal mengambil daftar catatan")
		return
	}
	defer rows.Close()

	notes := make([]Note, 0)
	for rows.Next() {
		var n Note
		if err := rows.Scan(&n.ID, &n.UserID, &n.Title, &n.Content, &n.CreatedAt, &n.UpdatedAt); err != nil {
			sendError(w, http.StatusInternalServerError, "Gagal membaca catatan")
			return
		}
		notes = append(notes, n)
	}

	sendJSON(w, http.StatusOK, notes)
}

// Create membuat catatan baru untuk pengguna yang sedang login
func (h *NotesHandler) Create(w http.ResponseWriter, r *http.Request) {
	userID, ok := middleware.GetUserID(r.Context())
	if !ok {
		sendError(w, http.StatusUnauthorized, "Sesi tidak ditemukan")
		return
	}

	var req NoteRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		sendError(w, http.StatusBadRequest, "Format JSON tidak valid")
		return
	}

	req.Title = strings.TrimSpace(req.Title)
	if req.Title == "" {
		sendError(w, http.StatusBadRequest, "Judul catatan tidak boleh kosong")
		return
	}
	if len(req.Title) > 200 {
		sendError(w, http.StatusBadRequest, "Judul catatan maksimal 200 karakter")
		return
	}

	var note Note
	query := `
		INSERT INTO notes (user_id, title, content) 
		VALUES ($1, $2, $3) 
		RETURNING id, user_id, title, content, created_at, updated_at
	`
	err := h.pool.QueryRow(r.Context(), query, userID, req.Title, req.Content).Scan(
		&note.ID, &note.UserID, &note.Title, &note.Content, &note.CreatedAt, &note.UpdatedAt,
	)
	if err != nil {
		sendError(w, http.StatusInternalServerError, "Gagal menyimpan catatan")
		return
	}

	sendJSON(w, http.StatusCreated, note)
}

// GetByID mengambil satu catatan berdasarkan ID (hanya jika dimiliki oleh pengguna tersebut)
func (h *NotesHandler) GetByID(w http.ResponseWriter, r *http.Request) {
	userID, ok := middleware.GetUserID(r.Context())
	if !ok {
		sendError(w, http.StatusUnauthorized, "Sesi tidak ditemukan")
		return
	}

	idStr := r.PathValue("id")
	noteID, err := strconv.Atoi(idStr)
	if err != nil {
		sendError(w, http.StatusBadRequest, "ID catatan tidak valid")
		return
	}

	var note Note
	query := `
		SELECT id, user_id, title, content, created_at, updated_at 
		FROM notes 
		WHERE id = $1 AND user_id = $2
	`
	err = h.pool.QueryRow(r.Context(), query, noteID, userID).Scan(
		&note.ID, &note.UserID, &note.Title, &note.Content, &note.CreatedAt, &note.UpdatedAt,
	)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			sendError(w, http.StatusNotFound, "Catatan tidak ditemukan atau Anda tidak memiliki akses")
			return
		}
		sendError(w, http.StatusInternalServerError, "Gagal memuat catatan")
		return
	}

	sendJSON(w, http.StatusOK, note)
}

// Update memperbarui judul dan isi catatan milik pengguna
func (h *NotesHandler) Update(w http.ResponseWriter, r *http.Request) {
	userID, ok := middleware.GetUserID(r.Context())
	if !ok {
		sendError(w, http.StatusUnauthorized, "Sesi tidak ditemukan")
		return
	}

	idStr := r.PathValue("id")
	noteID, err := strconv.Atoi(idStr)
	if err != nil {
		sendError(w, http.StatusBadRequest, "ID catatan tidak valid")
		return
	}

	var req NoteRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		sendError(w, http.StatusBadRequest, "Format JSON tidak valid")
		return
	}

	req.Title = strings.TrimSpace(req.Title)
	if req.Title == "" {
		sendError(w, http.StatusBadRequest, "Judul catatan tidak boleh kosong")
		return
	}
	if len(req.Title) > 200 {
		sendError(w, http.StatusBadRequest, "Judul catatan maksimal 200 karakter")
		return
	}

	var note Note
	query := `
		UPDATE notes 
		SET title = $1, content = $2, updated_at = CURRENT_TIMESTAMP 
		WHERE id = $3 AND user_id = $4 
		RETURNING id, user_id, title, content, created_at, updated_at
	`
	err = h.pool.QueryRow(r.Context(), query, req.Title, req.Content, noteID, userID).Scan(
		&note.ID, &note.UserID, &note.Title, &note.Content, &note.CreatedAt, &note.UpdatedAt,
	)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			sendError(w, http.StatusNotFound, "Catatan tidak ditemukan atau Anda tidak memiliki akses")
			return
		}
		sendError(w, http.StatusInternalServerError, "Gagal memperbarui catatan")
		return
	}

	sendJSON(w, http.StatusOK, note)
}

// Delete menghapus catatan milik pengguna
func (h *NotesHandler) Delete(w http.ResponseWriter, r *http.Request) {
	userID, ok := middleware.GetUserID(r.Context())
	if !ok {
		sendError(w, http.StatusUnauthorized, "Sesi tidak ditemukan")
		return
	}

	idStr := r.PathValue("id")
	noteID, err := strconv.Atoi(idStr)
	if err != nil {
		sendError(w, http.StatusBadRequest, "ID catatan tidak valid")
		return
	}

	query := `DELETE FROM notes WHERE id = $1 AND user_id = $2`
	tag, err := h.pool.Exec(r.Context(), query, noteID, userID)
	if err != nil {
		sendError(w, http.StatusInternalServerError, "Gagal menghapus catatan")
		return
	}

	if tag.RowsAffected() == 0 {
		sendError(w, http.StatusNotFound, "Catatan tidak ditemukan atau Anda tidak memiliki akses")
		return
	}

	sendJSON(w, http.StatusOK, map[string]string{
		"message": "Catatan berhasil dihapus",
	})
}

