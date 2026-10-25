package main

import (
	"bufio"
	"context"
	_ "embed"
	"log"
	"net/http"
	"os"
	"strings"
	"time"

	"vibecode/database"
	"vibecode/handlers"
	"vibecode/middleware"
)

//go:embed database/schema.sql
var schemaSQL string

func main() {
	// 1. Muat file .env jika ada (tanpa menimpa environment variable sistem yang sudah ada)
	loadEnv(".env")

	// 2. Baca konfigurasi wajib
	databaseURL := os.Getenv("DATABASE_URL")
	if databaseURL == "" {
		log.Fatal("Error: Environment variable DATABASE_URL belum diatur. Silakan periksa file .env Anda.")
	}

	port := os.Getenv("APP_PORT")
	if port == "" {
		port = "8080"
	}

	// 3. Hubungkan ke database PostgreSQL
	pool, err := database.Connect(databaseURL)
	if err != nil {
		log.Fatalf("Gagal terhubung ke database: %v", err)
	}
	defer pool.Close()
	log.Println("Berhasil terhubung ke database PostgreSQL.")

	// 4. Inisialisasi tabel database otomatis
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	if err := database.Migrate(ctx, pool, schemaSQL); err != nil {
		log.Fatalf("Gagal inisialisasi schema database: %v", err)
	}
	log.Println("Schema database berhasil diverifikasi/diperbarui.")

	// 5. Inisialisasi handler
	authHandler := handlers.NewAuthHandler(pool)
	notesHandler := handlers.NewNotesHandler(pool)
	authMiddleware := middleware.AuthMiddleware(pool)

	mux := http.NewServeMux()

	// 6. Daftarkan rute API Publik
	mux.HandleFunc("POST /api/register", authHandler.Register)
	mux.HandleFunc("POST /api/login", authHandler.Login)
	mux.HandleFunc("POST /api/logout", authHandler.Logout)

	// 7. Daftarkan rute API Terproteksi (wajib login)
	mux.Handle("GET /api/me", authMiddleware(http.HandlerFunc(authHandler.Me)))
	mux.Handle("GET /api/notes", authMiddleware(http.HandlerFunc(notesHandler.List)))
	mux.Handle("POST /api/notes", authMiddleware(http.HandlerFunc(notesHandler.Create)))
	mux.Handle("GET /api/notes/{id}", authMiddleware(http.HandlerFunc(notesHandler.GetByID)))
	mux.Handle("PUT /api/notes/{id}", authMiddleware(http.HandlerFunc(notesHandler.Update)))
	mux.Handle("DELETE /api/notes/{id}", authMiddleware(http.HandlerFunc(notesHandler.Delete)))

	// 8. Daftarkan penyajian file frontend statis
	fileServer := http.FileServer(http.Dir("web"))
	mux.Handle("/", fileServer)

	// Route redirect ramah pengguna
	mux.HandleFunc("GET /login", func(w http.ResponseWriter, r *http.Request) {
		http.Redirect(w, r, "/login.html", http.StatusSeeOther)
	})
	mux.HandleFunc("GET /register", func(w http.ResponseWriter, r *http.Request) {
		http.Redirect(w, r, "/register.html", http.StatusSeeOther)
	})

	serverAddr := ":" + port
	log.Printf("Server Personal Notes berjalan di http://localhost%s\n", serverAddr)
	if err := http.ListenAndServe(serverAddr, mux); err != nil {
		log.Fatalf("Server berhenti dengan error: %v", err)
	}
}

// loadEnv membaca file .env secara manual tanpa dependensi eksternal
func loadEnv(filepath string) {
	file, err := os.Open(filepath)
	if err != nil {
		// Jika file .env tidak ditemukan, lewati (menggunakan environment variable sistem)
		return
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if line == "" || strings.HasPrefix(line, "#") {
			continue
		}

		parts := strings.SplitN(line, "=", 2)
		if len(parts) == 2 {
			key := strings.TrimSpace(parts[0])
			val := strings.TrimSpace(parts[1])
			// Hanya set jika belum diatur di sistem
			if os.Getenv(key) == "" {
				_ = os.Setenv(key, val)
			}
		}
	}
}
