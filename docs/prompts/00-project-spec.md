# Project Specification — Vibe Coding

> **Tugas 1:** Isi spesifikasi aplikasi sebelum mulai melakukan vibe coding.
>
> Tentukan kebutuhan aplikasi dan pilih teknologi yang akan digunakan. Jangan mulai membuat kode sebelum spesifikasi selesai.

## 1. Informasi Aplikasi

**Nama aplikasi:**  
Personal Notes

**Deskripsi aplikasi:**  
Aplikasi web sederhana untuk menyimpan, melihat, menambahkan, mengedit, dan menghapus catatan pribadi. Setiap pengguna harus login untuk mengakses catatannya sendiri.

**Target pengguna:**  
Pengguna yang ingin menyimpan catatan pribadi melalui aplikasi web.

**Masalah yang ingin diselesaikan:**  
Pengguna membutuhkan tempat untuk menyimpan dan mengelola catatan secara digital dengan aman, tanpa tercampur dengan catatan pengguna lain.

## 2. Kebutuhan Aplikasi

Fitur utama aplikasi:

1. Register akun menggunakan nama, email, dan password.
2. Login dan logout.
3. Melihat daftar catatan milik pengguna yang sedang login.
4. Membuat catatan baru dengan judul dan isi.
5. Mengedit catatan yang sudah dibuat.
6. Menghapus catatan setelah konfirmasi.
7. Memastikan setiap pengguna hanya dapat melihat dan mengelola catatannya sendiri.

**Halaman aplikasi:**

- Login
- Register
- Main Page — My Notes

Pada Main Page, pengguna dapat melihat daftar catatan, membuat catatan, mengedit, menghapus, dan logout.

## 3. Pilihan Teknologi

### A. Backend Framework

- [x] 1. Go Standard Library (`net/http`)
- [ ] 2. Gin
- [ ] 3. Lainnya

**Pilihan saya:** 1

### B. ORM (Object-Relational Mapping)

- [x] 1. Tidak menggunakan ORM (SQL manual)
- [ ] 2. GORM
- [ ] 3. Lainnya

**Pilihan saya:** 1

**Ketentuan:** Gunakan SQL manual dengan library `pgx` untuk koneksi dan query PostgreSQL. Jangan menambahkan ORM.

### C. Database

- [x] 1. PostgreSQL
- [ ] 2. MySQL
- [ ] 3. SQLite

**Pilihan saya:** 1

### D. Frontend

- [x] 1. HTML, CSS, JavaScript murni
- [ ] 2. React
- [ ] 3. Vue

**Pilihan saya:** 1

**Ketentuan:** Frontend menggunakan HTML, CSS, dan JavaScript tanpa framework frontend.

### E. Sistem Autentikasi

- [ ] 1. Tidak membutuhkan autentikasi
- [x] 2. Session
- [ ] 3. JWT

**Pilihan saya:** 2

**Ketentuan:** Gunakan session berbasis server untuk mengelola status login pengguna.

### F. Penyimpanan Sesi atau Token

- [x] 1. Cookie
- [ ] 2. localStorage
- [ ] 3. Tidak diperlukan

**Pilihan saya:** 1

**Ketentuan:** Simpan ID session dalam cookie dengan pengaturan keamanan yang sesuai. Jangan menyimpan password atau session ID di `localStorage`.

### G. API

- [x] 1. REST API
- [ ] 2. Tidak membutuhkan API terpisah

**Pilihan saya:** 1

**Ketentuan:** Backend menyediakan REST API untuk autentikasi dan pengelolaan catatan. Frontend mengakses API menggunakan `fetch()`.

### H. Version Control

- [x] 1. Git
- [ ] 2. Tidak menggunakan Git

**Pilihan saya:** 1

## 4. Detail Autentikasi

**Registrasi akun:** Ya

**Login:** Ya

**Logout:** Ya

**Role pengguna:** User saja. Tidak ada role Admin.

**Aturan akses:**

1. Pengguna harus login untuk mengakses catatan.
2. Pengguna yang belum login hanya dapat mengakses halaman Login dan Register.
3. Setiap catatan harus memiliki `user_id` sebagai pemiliknya.
4. Backend harus memeriksa kepemilikan catatan pada setiap operasi melihat detail, mengedit, dan menghapus.
5. Pengguna tidak boleh melihat atau mengubah catatan pengguna lain, termasuk dengan mengubah ID pada URL atau request API.
6. Password harus disimpan dalam bentuk hash menggunakan algoritma yang sesuai, bukan sebagai teks biasa.
7. Setelah logout, session harus dinonaktifkan dan cookie session dibersihkan.

## 5. Konfigurasi Environment dan Database

**Ketentuan wajib:**

1. Semua konfigurasi aplikasi yang sensitif harus menggunakan environment variable.
2. Gunakan file `.env` untuk konfigurasi saat pengembangan lokal.
3. Sediakan file `.env.example` sebagai contoh konfigurasi tanpa kredensial asli.
4. File `.env` tidak boleh dimasukkan ke Git.
5. URL atau connection string PostgreSQL harus dibaca dari environment variable bernama `DATABASE_URL`.
6. Jangan menulis URL database, username, password database, session secret, atau kredensial lainnya secara langsung di kode.
7. Aplikasi harus memberikan pesan error yang jelas jika konfigurasi wajib belum tersedia, tanpa membocorkan password atau rahasia lainnya.

Contoh `.env.example`:

```dotenv
APP_PORT=8080
DATABASE_URL=postgres://username:password@localhost:5432/personal_notes?sslmode=disable
SESSION_SECRET=replace-with-a-random-secret
```

**Catatan:** Nilai di atas hanya contoh untuk pengembangan lokal. Gunakan kredensial dan session secret milik sendiri pada file `.env`. Untuk database produksi, gunakan konfigurasi SSL sesuai persyaratan penyedia database.

## 6. Batasan Pengembangan

- [x] Struktur folder sederhana.
- [x] Kode mudah dipahami pemula.
- [x] Tidak menambahkan teknologi di luar pilihan tanpa persetujuan.
- [x] Memisahkan frontend dan backend.
- [x] Menggunakan environment variable untuk konfigurasi.
- [x] Membuat pengujian dasar.
- [x] Menjelaskan kode dan perubahan yang dibuat AI.

**Batasan tambahan:**

- Gunakan Go Standard Library untuk HTTP server dan routing.
- Gunakan `pgx` untuk koneksi dan query PostgreSQL.
- Jangan menggunakan ORM.
- Jangan menambahkan framework frontend.
- Jangan membuat fitur di luar spesifikasi tanpa persetujuan.
- Gunakan parameterized query untuk mencegah SQL injection.
- Jangan mengembalikan password hash atau session secret melalui API.
- Gunakan validasi input di backend.
- Implementasikan aplikasi secara bertahap.
- Buat perubahan sekecil mungkin dan jelaskan cara menguji hasilnya.

## 7. Pemeriksaan Sebelum Vibe Coding

- [x] Nama dan tujuan aplikasi sudah jelas.
- [x] Fitur utama sudah ditentukan.
- [x] Pilihan teknologi sudah diisi.
- [x] Kebutuhan autentikasi sudah ditentukan.
- [x] Batasan proyek sudah ditentukan.
- [x] Penggunaan `.env` dan `DATABASE_URL` sudah ditentukan.

**Status spesifikasi:** Siap untuk diperiksa dan disetujui sebelum implementasi.

---

## 8. Instruksi untuk Agy

Baca dan pahami file `docs/prompts/00-project-spec.md` secara lengkap.

Gunakan spesifikasi ini sebagai acuan utama pengembangan.

1. Periksa apakah spesifikasi sudah lengkap dan konsisten.
2. Tanyakan jika ada informasi penting yang belum jelas.
3. Usulkan struktur folder sederhana sesuai teknologi yang dipilih.
4. Jelaskan fungsi folder dan file utama.
5. Jangan langsung membuat seluruh aplikasi sekaligus.
6. Jangan membuat, mengubah, atau menghapus file sebelum saya menyetujui usulan struktur folder.
7. Implementasikan aplikasi secara bertahap setelah mendapat persetujuan.
8. Jelaskan perubahan kode dan cara menguji setiap tahap.
9. Jangan mengganti teknologi atau menambahkan dependensi tanpa menjelaskan alasannya dan meminta persetujuan saya.
10. Pastikan `DATABASE_URL` dibaca dari environment variable dan tidak pernah di-hardcode.
11. Pastikan `.env` masuk ke `.gitignore`, sementara `.env.example` boleh disimpan di Git.
12. Jangan menampilkan nilai rahasia dari `.env` pada log atau respons API.

**Tahap pertama:** periksa spesifikasi dan usulkan struktur folder saja. Jangan mulai menulis kode sampai saya menyetujui struktur tersebut.
