# Arsitektur & Internal Wolf404 v1.3

**Created by ishowpen**

Wolf404 menggunakan _Tree-Walking Interpreter_ yang ditulis dalam bahasa Go, dirancang khusus untuk membangun aplikasi web dengan filosofi "Modern Javanese MVC".

## Alur Kerja (Pipeline)

1.  **Source Code** (`.wlf`) -> Input teks.
2.  **Lexer** (`compiler/lexer`): Memecah input menjadi token-token. Menangani **Status Indentasi** dengan menyisipkan token `INDENT` dan `DEDENT`.
3.  **Parser** (`compiler/parser`): Menggunakan teknik **Pratt Parsing** untuk membangun AST.
4.  **Evaluator** (`compiler/evaluator`): Mengeksekusi AST secara rekursif.
5.  **Javanese-Blade Compiler**: Sebelum evaluasi, file view diproses oleh engine native yang menangani perwarisan (`@warisan`) dan komponen (`@leboke`).

## Fitur Utama Framework

### 1. Template Engine: Javanese-Blade

Mesin template Wolf404 memiliki fitur setara Laravel Blade:

- **Inheritance**: `@warisan("folder.layout")` untuk mewarisi struktur.
- **Components**: `@leboke("folder.view")` untuk menyertakan partials.
- **Yield/Slot**: `@panggonan("name")` dan `@bagean("name")` untuk modulasi konten.
- **Native Logic**: Mendukung loop (`@track_neng`) dan kondisi (`@yen`) langsung di HTML.

### 2. Database: Wolf-ORM

- **Base Model**: Melakukan mapping tabel ke objek secara otomatis.
- **Security**: Menggunakan _Prepared Statements_ sebagai standar keamanan dari SQL Injection.

### 3. Keamanan (Security First)

- **CSRF Protection**: Token otomatis pada form dengan direktif `@csrf`.
- **XSS Protection**: Auto-escaping pada semua output `{{ }}`.

## Struktur Folder

- `app/`: Logika aplikasi (Models, Controllers, Middleware).
- `bootstrap/`: Inisialisasi framework.
- `config/`: Pengaturan aplikasi dan database.
- `public/`: Aset statis yang dapat diakses publik (CSS, JS, Images).
- `resources/`: Sumber daya mentah (Views mentah, CSS/JS source).
- `routes/`: Definisi routing web dan API.
- `system/`: Jantung framework (Router mentah, Base Model, Helpers).
