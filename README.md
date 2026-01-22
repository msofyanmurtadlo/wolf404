# 🐺 Bahasa Pemrograman Wolf404

**Wolf404** adalah bahasa pemrograman modern, terintepretasi, yang dibangun untuk efisiensi, kemudahan baca, dan konkurensi. Bahasa ini merupakan "kimera" yang menggabungkan DNA terbaik dari Python, PHP, dan Go:

- 🐍 **Mudah Dibaca**: Sintaks blok berbasis indentasi (Gaya Python).
- 🐘 **Web-Native**: Variabel dengan awalan `$` & tipe dinamis (Gaya PHP).
- 🐹 **Konkuren**: Thread ringan bawaan dengan kata kunci `prowl` (Gaya Go).

## 🚀 Memulai Cepat

### 1. Kompilasi (Build)

Wolf404 ditulis dengan Go. Anda perlu menginstal Go untuk membangun kompilernya.

```bash
cd compiler
go build -o ../wlf.exe main.go   # Untuk Windows
# go build -o ../wlf main.go     # Untuk Linux/Mac
cd ..
```

### 2. Jalankan Kode Pertama Anda

Buat file bernama `halo.wlf`:

```w404
howl("Halo, Kawanan Serigala!")

$jumlah = 0
sniff $jumlah < 5
    howl("Hitungan ke " + $jumlah)
```

Jalankan:

```bash
./wlf run halo.wlf
```

## 🌟 Fitur Unggulan

### Konkurensi (`prowl`)

Jalankan tugas berat di latar belakang secara instan tanpa memblokir program utama.

```w404
prowl kirim_email("user@contoh.com")
howl("Email sedang dikirim di background...")
```

### Sintaks Fungsional Bersih

Fungsi adalah warga kelas satu (_First-class citizens_).

```w404
$tambah = hunt($a, $b)
    bring $a + $b

howl($tambah(10, 20))
```

### Struktur Data Modern

Dukungan bawaan untuk Array dan Hash Map (Objek mirip JSON).

```w404
$profil = {
    "nama": "Wolfie",
    "level": 99,
    "skill": ["gigit", "howl", "prowl"]
}

howl($profil["nama"])
```

## 📚 Dokumentasi

Dokumentasi lengkap tersedia di direktori [docs/](docs/):

- [**Panduan Instalasi**](docs/INSTALLATION.md) - Cara build untuk Windows, Linux, dan macOS.
- [**Sintaks Bahasa**](docs/SYNTAX.md) - Pelajari variabel, loop, logika if, dan fungsi.
- [**Pustaka Standar**](docs/STDLIB.md) - Referensi fungsi bawaan.
- [**Arsitektur**](docs/ARCHITECTURE.md) - Cara kerja internal Interpreter/Compiler.

## 🤝 Kontribusi

Wolf404 adalah proyek sumber terbuka (open-source). Mari bergabung dengan kawanan!

## 📄 Lisensi

Lisensi MIT.
