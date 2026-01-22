# Pustaka Standar (Standard Library)

Wolf404 dilengkapi dengan sekumpulan fungsi bawaan yang tersedia secara global.

## Fungsi Inti

### `howl(arg1, ...)`

Mencetak nilai ke standard output (konsol/terminal). Menerima banyak argumen sekaligus.

```w404
howl("Halo Dunia")
howl("Nilainya adalah:", $x)
```

### `len(arg)`

Mengembalikan panjang (length) dari string atau array.

```w404
$arr = [1, 2, 3]
howl(len($arr)) // Output: 3

$teks = "Wolf"
howl(len($teks)) // Output: 4
```

---

_Catatan: Modul standar lainnya seperti `http`, `fs` (file system), dan `json` sedang dalam tahap pengembangan._
