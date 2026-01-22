# Arsitektur & Internal

Wolf404 menggunakan _Tree-Walking Interpreter_ yang ditulis dalam bahasa Go.

## Alur Kerja (Pipeline)

1.  **Source Code** (`.wlf`) -> Input teks.
2.  **Lexer** (`compiler/lexer`): Memecah input menjadi token-token. Menangani **Status Indentasi** dengan menyisipkan token `INDENT` dan `DEDENT` berdasarkan analisa spasi.
3.  **Parser** (`compiler/parser`): Menggunakan teknik **Pratt Parsing** (Top-Down Operator Precedence) untuk membangun Abstract Syntax Tree (AST).
4.  **Evaluator** (`compiler/evaluator`): Menjelajahi AST secara rekursif dan mengeksekusi logika.
    - Mengelola **Environment** (Lingkup Variabel) untuk penyimpanan nilai.
    - Menangani **Object System** (`compiler/object`) untuk nilai runtime.

## Komponen Kunci

### Lexer

- Mendukung pelacakan indentasi gaya Python (Stack-based indent level).
- Mendukung lexing gaya variabel PHP (`$IDENT`).

### Evaluator

- **Tree-Walking**: Mengeksekusi node AST secara langsung.
- **Konkurensi**: Pernyataan `prowl` dipetakan 1:1 ke mekanisme _goroutine_ milik Go (`go func()`).

## Mengembangkan Bahasa (Extending)

Untuk menambahkan fungsi bawaan baru:

1.  Edit file `compiler/evaluator/builtins.go`.
2.  Tambahkan entri baru ke peta `builtins`.
3.  Lakukan re-build (kompilasi ulang).
