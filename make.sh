export GOARCH=amd64
export GOOS=linux
go build -ldflags="-w -s"-o distr/DealHunter-linux-x64

export GOARCH=amd64
export GOOS=windows
go build -ldflags="-w -s" -o distr/DealHunter-windows-x64.exe
