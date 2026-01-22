# Sintaks Bahasa Wolf404

## Variabel

Semua variabel dimulai dengan tanda `$`, mirip dengan PHP.

```w404
$nama = "Alpha"
$umur = 10
$aktif = true
$kosong = nil
```

## Struktur Data

### Array (Larik)

Daftar berurutan dimulai dari indeks 0.

```w404
$angka = [10, 20, 30]
howl($angka[0]) // Mencetak 10
```

### Hash Maps (Objek)

Penyimpanan Key-Value (mirip JSON).

```w404
$profil = {
    "username": "wolf123",
    "peran": "admin"
}
howl($profil["username"])
```

## Fungsi (`hunt`)

Fungsi didefinisikan menggunakan kata kunci `hunt`. Fungsi bisa disimpan ke dalam variabel.

```w404
// Fungsi anonim disimpan ke variabel
$sapa = hunt($nama)
    howl("Halo " + $nama)

$sapa("Luna")

// Fungsi dengan nilai kembalian ('bring')
$tambah = hunt($a, $b)
    bring $a + $b

$hasil = $tambah(5, 5)
```

## Alur Kontrol (Control Flow)

### Sniff (If-Else)

Menggunakan kata kunci `sniff`. Blok kode ditentukan oleh **indentasi** (gaya Python).

```w404
$kekuatan = 9000

sniff $kekuatan > 8000
    howl("It's over 9000!")   // Blok ter-indentasi
missing
    howl("Kelemahan terdeteksi.") // Blok Else (missing)
```

## Konkurensi (`prowl`)

Jalankan fungsi di latar belakang (_background_) menggunakan `prowl`.

```w404
$tugas_cepat = hunt()
    howl("Tugas cepat selesai")

$tugas_berat = hunt()
    // simulasi kerja berat...
    howl("Tugas berat selesai")

prowl $tugas_berat() // Jalan di background (async)
$tugas_cepat()       // Jalan di foreground (sync)
```
