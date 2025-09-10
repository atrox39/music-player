# Music Player - Golang WASM
## Simple music player managment with Golang using WASM

### Info

- Go version: go1.24.2 windows/amd64

### Set Environment for WASM
```powershell
$env:GOOS="js"
$env:GOARCH="wasm"
```

### Compile
```bash
go build -o main.wasm main.go
```

### Use `build.ps1` for build the `main.wasm`
