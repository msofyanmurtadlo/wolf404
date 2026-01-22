# 🐺 Wolf404 Framework

**Created by ishowpen**

**Framework Web Fullstack menggunakan Bahasa Wolf404 (Jawa) - Seperti Laravel!**

Ini adalah kerangka kerja (framework) lengkap yang siap digunakan untuk membangun aplikasi web modern dengan kearifan lokal.

## 📂 Struktur Proyek

```
wolf404/
├── app/                    # Logika Aplikasi
│   ├── Controllers/        # Controller HTTP
│   ├── Middleware/         # Middleware
│   └── Models/             # Model Database
├── bootstrap/              # Bootstrap Framework
├── config/                 # Konfigurasi
├── database/               # Migrasi & Seed
├── routes/                 # Deklarasi Route API & Web
├── system/                 # Core System Framework
├── wolf404-vscode/         # VS Code Extension (Ikon & Syntax)
├── compiler/               # Kode sumber interpreter Wolf404 (Go)
├── server.wlf              # Entry point aplikasi
├── wlf.exe                 # Binary interpreter Wolf404 (Windows)
└── README.md               # Dokumentasi ini
```

## 🎨 VS Code Icons & Syntax

Agar file `.wlf` memiliki logo Wolf404 dan fitur syntax highlighting di editor:

### Instalasi Otomatis (Direkomendasikan)

Gunakan perintah PowerShell berikut untuk menginstal ekstensi VS Code secara otomatis:

```powershell
Copy-Item -Recurse -Force "wolf404-vscode" "$env:USERPROFILE\.vscode\extensions\wolf404-language-1.0.1"
```

### Instalasi Manual

1. Salin folder `wolf404-vscode` ke direktori ekstensi VS Code: `%USERPROFILE%\.vscode\extensions\wolf404-language-1.0.1`
2. Restart VS Code.
3. **Aktifkan Tema Ikon**:
   - Tekan `Ctrl+Shift+P`
   - Ketik "File Icon Theme"
   - Pilih "**A Wolf404 Icons**"

Setelah itu, semua file dengan ekstensi `.wlf` akan menampilkan logo Wolf404! 🐺

## 🌟 Fitur Utama

- **Arsitektur MVC**: Terbagi menjadi Controller (`app/Controllers`), View (`resources/views`), dan Model (`app/Models`).
- **Real Database Engine**: Menggunakan SQLite asli (via `database/sql`) bukan simulasi. Mendukung penyimpanan data persisten.
- **Keamanan SQL**: Mendukung prepared statements (`?` placeholders) untuk mencegah SQL Injection secara otomatis.
- **CSRF Protection**: Melindungi route POST dari serangan Cross-Site Request Forgery menggunakan token dinamis (`@csrf`).
- **Dynamic Routing**: Mendukung parameter dinamis dalam route (contoh: `/api/users/{id}`).
- **Template Engine (Blade-style)**: Mendukung injeksi variabel (`{{ variabel }}`) dan direktif keamanan (`@csrf`).
- **XSS Protection**: Auto-escaping otomatis pada semua output template.
- **Sintaks Bilingual**: Anda dapat menulis kode dalam bahasa **Jawa** (`garap`, `ketok`) atau **Inggris** (`hunt`, `howl`).
- **CLI Seperti Artisan**: Gunakan perintah `wlf gas server` untuk menjalankan aplikasi Anda.

## 🚀 Cara Menjalankan (How to Run)

Anda sudah memiliki binary `wlf.exe` yang sudah dikompilasi. Anda bisa langsung menjalankannya!

### 1. Menjalankan Server

Buka terminal (PowerShell/CMD) di folder proyek ini, lalu ketik:

```powershell
.\wlf.exe gas server
```

Server akan berjalan di: `http://localhost:8080`

**Catatan:** Seperti Laravel (`php artisan serve`), jika port sudah digunakan, server akan otomatis memberikan peringatan dan berhenti dengan pesan yang jelas. Anda tidak perlu mematikan proses secara manual.

### 2. Uji Coba API & Tampilan

**Halaman Utama (Tailwind CSS):**
Buka browser dan akses: [http://localhost:8080](http://localhost:8080)

**API Login (Contoh):**

```bash
curl -X POST http://localhost:8080/api/auth/login
```

## 🛠️ Pengembangan (Development)

### Membuat Model & Migrasi (Gaya Laravel)

Anda dapat membuat Model dan Migrasi sekaligus dengan satu perintah:

```powershell
.\wlf.exe gawe:model NamaModel
```

Perintah ini akan secara otomatis membuat:

- `app/Models/NamaModel.wlf` - File Model dengan metode CRUD dasar.
- `database/migrations/TIMESTAMP_create_table.wlf` - File Migrasi untuk tabel terkait.

### Menambah Route

Edit file `routes/api.wlf` atau `routes/web.wlf`.

```wolf404
$router.get("/api/baru", garap($req)
    balekno http_json("Halo dari route baru!")
)
```

### Menambah Model Baru

Buat file baru di direktori `app/Models/`.

### Menambah Controller Baru

Buat file baru di direktori `app/Controllers/`.

## ⚙️ Setup Environment (Opsional)

Jika Anda ingin membangun ulang interpreter (misalnya setelah mengubah kode di folder `compiler/`):

1. Pastikan Go sudah terinstal di sistem Anda.
2. Jalankan perintah berikut:
   ```bash
   go build -o wlf.exe compiler/main.go
   ```

---

**Wolf404 Framework** - _Coding pusing? Tidak masalah, yang penting jalan!_ 🐺
