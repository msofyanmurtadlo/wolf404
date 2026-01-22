# Pustaka Standar (Standard Library) - Edisi Jowo Kasar

Wolf404 digowo nganggo pirang-pirang fungsi gowoan sing iso langsung dinggo.

## Fungsi Inti

### `ketok(arg1, ...)`

Ngetokne nilai neng standar output (konsol/terminal). Iso nrimo akeh argumen pisanan.

```w404
ketok("Halo Ndunyo")
ketok("Nilaine yoiku:", $x)
```

### `dowo(arg)`

Mbalekno dowone (length) string utowo array.

```w404
$arr = [1, 2, 3]
ketok(dowo($arr)) // Hasile: 3

$teks = "Wolf"
ketok(dowo($teks)) // Hasile: 4
```

### `takon(prompt)`

Njaluk input soko user.

```w404
$jeneng = takon("Sopo jenengmu? ")
ketok("Halo " + $jeneng)
```

### `moco_file(path)`

Moco isi file.

### `nulis_file(path, isi)`

Nulis isi neng file.

---

_Catetan: Modul standar liyane koyo `http`, `fs` (file system), lan `json` isih digarap._
