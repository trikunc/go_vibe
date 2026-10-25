// Logika Pengelolaan Catatan Frontend
document.addEventListener("DOMContentLoaded", async () => {
  const userGreeting = document.getElementById("userGreeting");
  const logoutBtn = document.getElementById("logoutBtn");
  const notesGrid = document.getElementById("notesGrid");
  const emptyState = document.getElementById("emptyState");
  const alertBox = document.getElementById("alertBox");

  // Modal elements
  const noteModal = document.getElementById("noteModal");
  const modalTitle = document.getElementById("modalTitle");
  const noteForm = document.getElementById("noteForm");
  const noteIdInput = document.getElementById("noteId");
  const noteTitleInput = document.getElementById("noteTitle");
  const noteContentInput = document.getElementById("noteContent");
  const btnNewNote = document.getElementById("btnNewNote");
  const btnCancelModal = document.getElementById("btnCancelModal");
  const btnSaveNote = document.getElementById("btnSaveNote");

  let currentNotes = [];

  function showAlert(message, type = "error") {
    if (!alertBox) return;
    alertBox.textContent = message;
    alertBox.className = `alert alert-${type}`;
    alertBox.style.display = "block";
    setTimeout(() => {
      alertBox.style.display = "none";
    }, 4000);
  }

  // 1. Verifikasi Status Login & Ambil Profil
  try {
    const res = await fetch("/api/me");
    if (res.status === 401) {
      window.location.href = "/login.html";
      return;
    }
    const user = await res.json();
    userGreeting.textContent = `Halo, ${user.name}`;
  } catch (err) {
    window.location.href = "/login.html";
    return;
  }

  // 2. Muat Daftar Catatan
  async function loadNotes() {
    try {
      const res = await fetch("/api/notes");
      if (res.status === 401) {
        window.location.href = "/login.html";
        return;
      }
      currentNotes = await res.json();
      renderNotes();
    } catch (err) {
      showAlert("Gagal memuat catatan dari server.");
    }
  }

  // 3. Render Catatan ke Layar
  function renderNotes() {
    notesGrid.innerHTML = "";

    if (!currentNotes || currentNotes.length === 0) {
      emptyState.style.display = "block";
      return;
    }

    emptyState.style.display = "none";

    currentNotes.forEach((note) => {
      const card = document.createElement("div");
      card.className = "note-card";

      const updatedDate = new Date(note.updated_at).toLocaleDateString("id-ID", {
        day: "numeric",
        month: "short",
        year: "numeric",
        hour: "2-digit",
        minute: "2-digit",
      });

      card.innerHTML = `
        <div>
          <h3 class="note-title">${escapeHtml(note.title)}</h3>
          <p class="note-content">${escapeHtml(note.content || "")}</p>
        </div>
        <div class="note-footer">
          <span class="note-date">Diperbarui: ${updatedDate}</span>
          <div class="note-actions">
            <button class="btn btn-secondary btn-sm btn-edit" data-id="${note.id}">Edit</button>
            <button class="btn btn-danger btn-sm btn-delete" data-id="${note.id}">Hapus</button>
          </div>
        </div>
      `;

      notesGrid.appendChild(card);
    });

    // Tambahkan event listener pada tombol edit dan hapus
    document.querySelectorAll(".btn-edit").forEach((btn) => {
      btn.addEventListener("click", () => openEditModal(btn.dataset.id));
    });

    document.querySelectorAll(".btn-delete").forEach((btn) => {
      btn.addEventListener("click", () => handleDeleteNote(btn.dataset.id));
    });
  }

  // 4. Modal Buka & Tutup
  function openCreateModal() {
    modalTitle.textContent = "Catatan Baru";
    noteIdInput.value = "";
    noteTitleInput.value = "";
    noteContentInput.value = "";
    noteModal.classList.remove("hidden");
    noteTitleInput.focus();
  }

  function openEditModal(id) {
    const note = currentNotes.find((n) => n.id === parseInt(id, 10));
    if (!note) return;

    modalTitle.textContent = "Edit Catatan";
    noteIdInput.value = note.id;
    noteTitleInput.value = note.title;
    noteContentInput.value = note.content || "";
    noteModal.classList.remove("hidden");
    noteTitleInput.focus();
  }

  function closeModal() {
    noteModal.classList.add("hidden");
  }

  btnNewNote.addEventListener("click", openCreateModal);
  btnCancelModal.addEventListener("click", closeModal);

  // Tutup modal jika klik di latar belakang gelap
  noteModal.addEventListener("click", (e) => {
    if (e.target === noteModal) closeModal();
  });

  // 5. Simpan Catatan (Buat Baru atau Update)
  noteForm.addEventListener("submit", async (e) => {
    e.preventDefault();

    const id = noteIdInput.value;
    const title = noteTitleInput.value.trim();
    const content = noteContentInput.value.trim();

    if (!title) {
      alert("Judul catatan tidak boleh kosong");
      return;
    }

    btnSaveNote.disabled = true;
    btnSaveNote.textContent = "Menyimpan...";

    try {
      const isEdit = Boolean(id);
      const url = isEdit ? `/api/notes/${id}` : "/api/notes";
      const method = isEdit ? "PUT" : "POST";

      const res = await fetch(url, {
        method: method,
        headers: { "Content-Type": "application/json" },
        body: JSON.stringify({ title, content }),
      });

      if (!res.ok) {
        const data = await res.json();
        alert(data.error || "Gagal menyimpan catatan");
        btnSaveNote.disabled = false;
        btnSaveNote.textContent = "Simpan";
        return;
      }

      closeModal();
      btnSaveNote.disabled = false;
      btnSaveNote.textContent = "Simpan";
      await loadNotes();
      showAlert(isEdit ? "Catatan berhasil diperbarui" : "Catatan berhasil ditambahkan", "success");
    } catch (err) {
      alert("Terjadi kesalahan koneksi.");
      btnSaveNote.disabled = false;
      btnSaveNote.textContent = "Simpan";
    }
  });

  // 6. Hapus Catatan dengan Konfirmasi
  async function handleDeleteNote(id) {
    const isConfirmed = confirm("Apakah Anda yakin ingin menghapus catatan ini?");
    if (!isConfirmed) return;

    try {
      const res = await fetch(`/api/notes/${id}`, {
        method: "DELETE",
      });

      if (!res.ok) {
        const data = await res.json();
        showAlert(data.error || "Gagal menghapus catatan");
        return;
      }

      await loadNotes();
      showAlert("Catatan berhasil dihapus", "success");
    } catch (err) {
      showAlert("Gagal menghapus catatan. Periksa koneksi Anda.");
    }
  }

  // 7. Logout
  logoutBtn.addEventListener("click", async () => {
    try {
      await fetch("/api/logout", { method: "POST" });
    } finally {
      window.location.href = "/login.html";
    }
  });

  // Helper untuk sanitasi HTML tampilan guna mencegah XSS di client
  function escapeHtml(str) {
    const div = document.createElement("div");
    div.textContent = str;
    return div.innerHTML;
  }

  // Muat data awal
  loadNotes();
});

