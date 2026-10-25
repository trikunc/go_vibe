// Logika Autentikasi Frontend (Login & Register)
document.addEventListener("DOMContentLoaded", () => {
  const loginForm = document.getElementById("loginForm");
  const registerForm = document.getElementById("registerForm");
  const alertBox = document.getElementById("alertBox");

  function showAlert(message, type = "error") {
    if (!alertBox) return;
    alertBox.textContent = message;
    alertBox.className = `alert alert-${type}`;
    alertBox.style.display = "block";
  }

  function hideAlert() {
    if (!alertBox) return;
    alertBox.style.display = "none";
  }

  // Handle Form Login
  if (loginForm) {
    loginForm.addEventListener("submit", async (e) => {
      e.preventDefault();
      hideAlert();

      const email = document.getElementById("email").value.trim();
      const password = document.getElementById("password").value;
      const submitBtn = document.getElementById("submitBtn");

      if (!email || !password) {
        showAlert("Email dan password wajib diisi");
        return;
      }

      submitBtn.disabled = true;
      submitBtn.textContent = "Memproses...";

      try {
        const res = await fetch("/api/login", {
          method: "POST",
          headers: { "Content-Type": "application/json" },
          body: JSON.stringify({ email, password }),
        });

        const data = await res.json();

        if (!res.ok) {
          showAlert(data.error || "Gagal masuk, silakan coba lagi.");
          submitBtn.disabled = false;
          submitBtn.textContent = "Masuk";
          return;
        }

        // Login berhasil, arahkan ke halaman utama
        window.location.href = "/";
      } catch (err) {
        showAlert("Tidak dapat terhubung ke server. Periksa koneksi Anda.");
        submitBtn.disabled = false;
        submitBtn.textContent = "Masuk";
      }
    });
  }

  // Handle Form Register
  if (registerForm) {
    registerForm.addEventListener("submit", async (e) => {
      e.preventDefault();
      hideAlert();

      const name = document.getElementById("name").value.trim();
      const email = document.getElementById("email").value.trim();
      const password = document.getElementById("password").value;
      const submitBtn = document.getElementById("submitBtn");

      if (!name || !email || !password) {
        showAlert("Semua kolom wajib diisi");
        return;
      }

      if (password.length < 6) {
        showAlert("Password minimal 6 karakter");
        return;
      }

      submitBtn.disabled = true;
      submitBtn.textContent = "Mendaftarkan...";

      try {
        const res = await fetch("/api/register", {
          method: "POST",
          headers: { "Content-Type": "application/json" },
          body: JSON.stringify({ name, email, password }),
        });

        const data = await res.json();

        if (!res.ok) {
          showAlert(data.error || "Gagal mendaftar akun.");
          submitBtn.disabled = false;
          submitBtn.textContent = "Daftar";
          return;
        }

        showAlert("Pendaftaran berhasil! Mengalihkan ke halaman login...", "success");
        setTimeout(() => {
          window.location.href = "/login.html";
        }, 1500);
      } catch (err) {
        showAlert("Tidak dapat terhubung ke server. Periksa koneksi Anda.");
        submitBtn.disabled = false;
        submitBtn.textContent = "Daftar";
      }
    });
  }
});

