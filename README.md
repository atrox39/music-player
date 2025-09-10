# Music Player - Golang WASM
## Simple music player managment with Golang using WASM

### Set Environment for WASM
```powershell
$env:GOOS="js"
$env:GOARCH="wasm"
```

### Compile
```bash
go build -o main.wasm main.go
```
