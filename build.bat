set GOARCH=amd64

set GOOS=linux
go build -ldflags "-s -w" -trimpath -o release/TxPortMap_linux_x64 cmd/TxPortMap/TxPortMap.go

set GOOS=windows
go build -ldflags "-s -w" -trimpath -o release/TxPortMap_windows_x64.exe cmd/TxPortMap/TxPortMap.go

set GOOS=darwin
go build -ldflags "-s -w" -trimpath -o release/TxPortMap_macos_x64 cmd/TxPortMap/TxPortMap.go
