set GOARCH=amd64
set GOOS=linux
go build -ldflags="-w -s" -o distr/DealHunter-linux-x64

set GOARCH=amd64
set GOOS=windows
go build -ldflags="-w -s" -o distr/DealHunter-windows-x64.exe
