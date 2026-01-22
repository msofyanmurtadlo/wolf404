# Sintaks Bahasa Wolf404

**Created by ishowpen**

Bahasa iki ndhukung loro mode: **Jowo** lan **Inggris**. Sampeyan iso milih sesuai selera.

| Jowo       | English | Fungsi               |
| ---------- | ------- | -------------------- |
| garap      | hunt    | Define Function      |
| gerombolan | mold    | Define Class (Class) |
| menowo     | sniff   | If Condition         |
| yenora     | missing | Else/Else If         |
| baleni     | track   | Loop (While)         |
| balekno    | bring   | Return               |
| ketok      | howl    | Print/Log            |
| undang     | summon  | Import/Include       |
| playon     | prowl   | Go Routine           |
| bener      | true    | Boolean True         |
| salah      | false   | Boolean False        |
| kopong     | nil     | Null/Nil             |

## Variabel

Kabeh variabel diwiwiti nganggo tondo `$`, koyo PHP.

```w404
$jeneng = "Alpha"
$umur = 10
$aktif = bener
$kosong = kopong  // utowo: nil
```

## Struktur Data

### Array (Larik)

```w404
$nomer = [10, 20, 30]
ketok($nomer[0])
// howl($nomer[0])
```

### Hash Maps (Objek)

```w404
$profil = {
    "user": "wolf123",
    "status": "admin"
}
ketok($profil["user"])
```

## Fungsi (`garap` / `hunt`)

Gawe fungsi nganggo tembung `garap` utowo `hunt`. Hasil ditokne nganggo `balekno` utowo `bring`.

```w404
$sapa = garap($name)
    ketok("Halo " + $name)

// English Style
hunt sum($a, $b)
    bring $a + $b
```

## Alur Kontrol

### Menowo (If-Else)

Nganggo `menowo` (`sniff`) lan `yenora` (`missing`). Blok kode diatur nganggo **indentasi**.

```w404
$tenogo = 9000

menowo $tenogo > 8000
    ketok("Sakti banget!")
yenora
    ketok("Cemen.")
```

### Baleni (Loop)

Nganggo tembung `baleni` utowo `track`.

```w404
$i = 0
track $i < 5
    howl("Angka: " + $i)
    $i = $i + 1
```

## Pemrograman Berorientasi Objek (`gerombolan` / `mold`)

```w404
gerombolan Kucing
    garap init($jeneng)
        $this.jeneng = $jeneng

    garap suara()
        ketok($this.jeneng + " muni: Meong!")

$tom = Kucing("Tom")
$tom.suara()
```

## Konkurensi (`playon` / `prowl`)

Jalanke fungsi neng background nganggo `playon` utowo `prowl`.

```w404
playon tugas_abot()
// prowl heavy_task()
ketok("Tugas jalan neng mburi...")
```
