# Sintaks Bahasa Wolf404 (Edisi Jowo Kasar)

## Variabel

Kabeh variabel diwiwiti nganggo tondo `$`, koyo PHP.

```w404
$jeneng = "Alpha"
$umur = 10
$aktif = bener
$kosong = kopong
```

## Struktur Data

### Array (Larik)

```w404
$nomer = [10, 20, 30]
ketok($nomer[0])
```

### Hash Maps (Objek)

```w404
$profil = {
    "user": "wolf123",
    "status": "admin"
}
ketok($profil["user"])
```

## Fungsi (`garap`)

Gawe fungsi nganggo tembung `garap`. Hasil ditokne nganggo `balekno`.

```w404
$sapa = garap($name)
    ketok("Halo " + $name)

$tambah = garap($a, $b)
    balekno $a + $b
```

## Alur Kontrol

### Menowo (If-Else)

Nganggo `menowo` lan `yenora`. Blok kode diatur nganggo **indentasi**.

```w404
$tenogo = 9000

menowo $tenogo > 8000
    ketok("Sakti banget!")
yenora
    ketok("Cemen.")
```

### Baleni (Loop)

Nganggo tembung `baleni`.

```w404
$i = 0
baleni $i < 5
    ketok("Angka: " + $i)
    $i = $i + 1
```

## Pemrograman Berorientasi Objek (`gerombolan`)

```w404
gerombolan Kucing
    garap init($jeneng)
        $this.jeneng = $jeneng

    garap suara()
        ketok($this.jeneng + " muni: Meong!")

$tom = Kucing("Tom")
$tom.suara()
```

## Konkurensi (`playon`)

Jalanke fungsi neng background nganggo `playon`.

```w404
playon tugas_abot()
ketok("Tugas jalan neng mburi...")
```
