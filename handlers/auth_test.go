package handlers

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"golang.org/x/crypto/bcrypt"
)

// TestPasswordHashing menguji keamanan proses hash password menggunakan bcrypt
func TestPasswordHashing(t *testing.T) {
	rawPassword := "rahasia123"

	hashed, err := bcrypt.GenerateFromPassword([]byte(rawPassword), bcrypt.DefaultCost)
	if err != nil {
		t.Fatalf("Gagal menghasilkan hash password: %v", err)
	}

	// Pastikan hash cocok dengan password asli
	if err := bcrypt.CompareHashAndPassword(hashed, []byte(rawPassword)); err != nil {
		t.Errorf("Password asli seharusnya cocok dengan hash")
	}

	// Pastikan password salah ditolak
	if err := bcrypt.CompareHashAndPassword(hashed, []byte("passwordsalah")); err == nil {
		t.Errorf("Password salah seharusnya tidak cocok dengan hash")
	}
}

// TestRegisterValidation menguji validasi input pendaftaran pengguna baru
func TestRegisterValidation(t *testing.T) {
	handler := &AuthHandler{pool: nil} // pool nil karena kita hanya menguji validasi request

	tests := []struct {
		name       string
		payload    RegisterRequest
		wantStatus int
	}{
		{
			name: "Nama kosong",
			payload: RegisterRequest{
				Name:     "",
				Email:    "test@example.com",
				Password: "password123",
			},
			wantStatus: http.StatusBadRequest,
		},
		{
			name: "Email tidak valid",
			payload: RegisterRequest{
				Name:     "Budi",
				Email:    "budi-bukan-email",
				Password: "password123",
			},
			wantStatus: http.StatusBadRequest,
		},
		{
			name: "Password kurang dari 6 karakter",
			payload: RegisterRequest{
				Name:     "Budi",
				Email:    "budi@example.com",
				Password: "123",
			},
			wantStatus: http.StatusBadRequest,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			body, _ := json.Marshal(tt.payload)
			req := httptest.NewRequest(http.MethodPost, "/api/register", bytes.NewReader(body))
			rec := httptest.NewRecorder()

			handler.Register(rec, req)

			if rec.Code != tt.wantStatus {
				t.Errorf("Register() status = %v, diharapkan %v", rec.Code, tt.wantStatus)
			}
		})
	}
}

