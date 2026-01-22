# 🐺 Bahasa Pemrograman Wolf404 (Edisi Jowo Kasar)

**Wolf404** adalah bahasa pemrograman modern, terintepretasi, yang dibangun untuk efisiensi, kemudahan baca, dan konkurensi dengan kearifan lokal Jowo Kasar. Bahasa ini merupakan "kimera" yang menggabungkan DNA terbaik dari Python, PHP, dan Go:

- 🐍 **Gampang Diwoco**: Sintaks blok nggo indentasi (Koyo Python).
- 🐘 **Web-Native**: Variabel nggo tondo `$` & tipe dinamis (Koyo PHP).
- 🐹 **Playon**: Fitur background job / thread enteng nggo kata kunci `playon` (Koyo Go).

## 🚀 Memulai Cepat

### 1. Kompilasi (Build)

Wolf404 ditulis nganggo Go. Sampeyan kudu nginstall Go nggo mbuild kompilere.

```bash
cd compiler
go build -o ../wlf.exe main.go   # Nggo Windows
# go build -o ../wlf main.go     # Nggo Linux/Mac
cd ..
```

### 2. Jalanke Kode Pertama

Gawe file jenenge `halo.wlf`:

```w404
ketok("Halo, Cah Serigala!")

$jumlah = 0
menowo $jumlah < 5
    ketok("Hitungan ke " + $jumlah)
```

Jalanke:

```bash
./wlf gas halo.wlf
```

## 🌟 Fitur Unggulan

### Konkurensi (`playon`)

Jalanke tugas abot neng background saknaliko tanpa mblokir program utomo.

```w404
playon kirim_email("user@contoh.com")
ketok("Email lagi dikirim neng background...")
```

### Sintaks Garap-garapan

Fungsi iku warga kelas siji (_First-class citizens_).

```w404
$tambah = garap($a, $b)
    balekno $a + $b

ketok($tambah(10, 20))
```

### Perulangan (`baleni`)

Sintaks looping sederhana mirip `while`.

```w404
$i = 0
baleni $i < 5
    ketok("Hitungan " + $i)
    $i = $i + 1
```

### OOP (`gerombolan`)

Mendukung Class lan Object Instance (Gerombolan).

```w404
gerombolan Hero
    garap serang()
        ketok("Serangan bertubi-tubi sak enak udele!")

$pahlawan = Hero()
$pahlawan.serang()
```

### Struktur Data Modern

Dukungan bawaan nggo Array lan Hash Map (Objek mirip JSON).

```w404
$profil = {
    "nama": "Wolfie",
    "level": 99,
    "skill": ["gigit", "ketok", "playon"]
}

ketok($profil["nama"])
```

## 📚 Dokumentasi

Dokumentasi lengkap kasedhiya neng direktori [docs/](docs/):

- [**Panduan Instalasi**](docs/INSTALLATION.md) - Cara build nggo Windows, Linux, lan macOS.
- [**Sintaks Boso**](docs/SYNTAX.md) - Sinau variabel, loop, logika if, lan fungsi.
- [**Pustaka Standar**](docs/STDLIB.md) - Referensi fungsi gowoan.
- [**Arsitektur**](docs/ARCHITECTURE.md) - Cara kerjo internal Interpreter/Compiler.

## 🤝 Kontribusi

Wolf404 iku proyek sumber terbuka (open-source). Ayo gabung neng kawanan!

## 📄 Lisensi

Lisensi MIT.
