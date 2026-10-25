package handlers

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"vibecode/middleware"
)

// TestNoteValidation menguji validasi input pembuatan catatan
func TestNoteValidation(t *testing.T) {
	handler := &NotesHandler{pool: nil}

	tests := []struct {
		name       string
		payload    NoteRequest
		wantStatus int
	}{
		{
			name: "Judul kosong",
			payload: NoteRequest{
				Title:   "",
				Content: "Isi catatan",
			},
			wantStatus: http.StatusBadRequest,
		},
		{
			name: "Judul melebihi 200 karakter",
			payload: NoteRequest{
				Title:   strings.Repeat("a", 201),
				Content: "Isi catatan",
			},
			wantStatus: http.StatusBadRequest,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			body, _ := json.Marshal(tt.payload)
			req := httptest.NewRequest(http.MethodPost, "/api/notes", bytes.NewReader(body))
			// Simulasikan request yang sudah terautentikasi (memiliki userID di context)
			ctx := context.WithValue(req.Context(), middleware.UserIDKey, 1)
			req = req.WithContext(ctx)

			rec := httptest.NewRecorder()
			handler.Create(rec, req)

			if rec.Code != tt.wantStatus {
				t.Errorf("Create() status = %v, diharapkan %v", rec.Code, tt.wantStatus)
			}
		})
	}
}

