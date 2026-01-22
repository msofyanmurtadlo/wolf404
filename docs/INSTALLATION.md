# Panduan Instalasi & Kompilasi

Wolf404 adalah _toolchain_ bahasa yang ditulis dengan Go. Artinya, konsep "Single Binary Deployment" berlaku di sini. Anda bisa membuat (build) program di satu komputer dan menjalankannya di komputer lain tanpa perlu menginstal runtime tambahan.

## Prasyarat

- [Go 1.18+](https://go.dev/dl/) terinstal di komputer pengembangan (laptop) Anda.

## Membuat (Build) dari Sumber

### Di Windows

Buka PowerShell/CMD di folder utama proyek:

```powershell
cd compiler
go build -o ../wlf.exe main.go
```

Anda akan mendapatkan file `wlf.exe` di folder utama.

### Di Linux / macOS

```bash
cd compiler
go build -o ../wlf main.go
chmod +x ../wlf
```

## Cross-Compilation (Build untuk Server)

Salah satu kekuatan super Wolf404 adalah kemampuan kompilasi silang (_cross-compilation_) berkat Go. Anda bisa membuat binary Linux langsung dari Windows.

### Build untuk Server Linux (dari Windows)

Gunakan perintah ini di PowerShell:

```powershell
$env:GOOS = "linux"
$env:GOARCH = "amd64"
go build -o ../wlf_linux compiler/main.go
# Kembalikan settingan (Reset)
$env:GOOS = "windows"
```

Upload file `wlf_linux` ke server Ubuntu/CentOS/Debian Anda dan jalankan:

```bash
chmod +x wlf_linux
./wlf_linux run script_saya.wlf
```

### Build untuk macOS (Apple Silicon M1/M2)

```powershell
$env:GOOS = "darwin"
$env:GOARCH = "arm64"
go build -o ../wlf_mac_m1 compiler/main.go
```
