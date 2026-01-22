# 🐺 Wolf404 Framework

**Created by ishowpen**

**Fullstack Web Framework nganggo Bahasa Wolf404 (Jowo) - Kaya Laravel!**

Ini adalah kerangka kerja (framework) lengkap siap pakai.

## 📂 Struktur Project

```
wolf404/
├── app/                    # Application Logic
│   ├── Controllers/        # HTTP Controllers
│   ├── Middleware/         # Middleware
│   └── Models/            # Database Models
├── bootstrap/              # Framework Bootstrap
├── config/                 # Configuration
├── database/               # Migrations & Seeds
├── routes/                 # API Routes declaration
├── system/                 # Core System Framework
├── wolf404-vscode/         # VS Code Extension (Icons & Syntax)
├── compiler/               # Sumber kode interpreter Wolf404 (Go)
├── server.wlf              # Entry point aplikasi
├── wlf.exe                 # Binary interpreter Wolf404 (Windows)
└── README.md               # Dokumentasi iki
```

## 🎨 VS Code Icons & Syntax

Agar `.wlf` file memiliki logo Wolf404 dan syntax highlighting di editor:

1. Copy folder `wolf404-vscode` ke folder extensions VS Code:
   `%USERPROFILE%\.vscode\extensions\`
2. Restart VS Code.
3. Pilih "Wolf404 Icons" di File > Preferences > File Icon Theme.

## 🌟 Features

- **MVC Architecture**: Controllers (`app/Controllers`), Views (`resources/views`), and Models (`app/Models`).
- **Tailwind CSS Ready**: Included via CDN in default views for rapid UI development.
- **Bilingual Syntax**: Write code in **Javanese** (`garap`, `ketok`) or **English** (`hunt`, `howl`).
- **Artisan-like CLI**: Use `wlf gas server` to start your application.
- **SQLite Database**: Built-in support for persistent storage.

## 🚀 Cara Jalanke (How to Run)

Sampeyan wis duwe binary `wlf.exe` sing wis dicompile. Langsung gas wae!

### 1. Jalanke Server

Bukak terminal (PowerShell/CMD) neng folder project iki, terus ketik:

```powershell
.\wlf.exe gas server
```

Server bakal mlaku neng: `http://localhost:8080`

**Catatan:** Kaya Laravel (`php artisan serve`), yen port wis digunakan, server bakal otomatis ngasih tau lan metu kanthi pesan sing jelas. Ora perlu manual kill process.

### 2. Test API & View

**Halaman Utama (Tailwind CSS):**
Buka browser: [http://localhost:8080](http://localhost:8080)

**Login API (Contoh):**

```bash
curl -X POST http://localhost:8080/api/auth/login
```

## 🛠️ Pengembangan (Development)

### Gawe Model & Migration (Laravel-style)

Generate Model dan Migration sekaligus:

```powershell
.\wlf.exe gawe:model Product
```

Iki bakal gawe:

- `app/Models/Product.wlf` - Model file kanthi CRUD methods
- `database/migrations/TIMESTAMP_create_products_table.wlf` - Migration file

### Nambah Route

Edit file `routes/api.wlf` utowo `routes/web.wlf`.

```wolf404
$router.get("/api/anyar", garap($req)
    balekno http_json("Halo saka route anyar!")
)
```

### Nambah Model

Gawe file anyar neng `app/Models/`.

### Nambah Controller

Gawe file anyar neng `app/Controllers/`.

## ⚙️ Setup Environment (Opsional)

Yen pengen mbuild ulang interpreter (contone bar ngedit folder `compiler/`):

1. Pastikan wis install Go.
2. Jalanke perintah:
   ```bash
   go build -o wlf.exe compiler/main.go
   ```

---

**Wolf404 Framework** - _Coding mumet? Ora masalah, sing penting mlaku!_ 🐺
