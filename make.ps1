$env:OLDGOOS=$env:GOOS
$env:OLDGOARCH=$env:GOARCH

$env:GOARCH="amd64"
$env:GOOS="linux"
go build -ldflags="-w -s" -o distr/DealHunter-linux-x64

#$env:GOARCH="386"
#$env:GOOS="linux"
#go build -o distr/WncPlugin-linux-x86
#
#$env:GOARCH="amd64"
#$env:GOOS="windows"
#go build -o distr/WncPlugin-windows-x64.exe
#
#$env:GOARCH="amd64"
#$env:GOOS="darwin"
#go build -o distr/WncPlugin-macos-x64

$env:GOOS=$env:OLDGOOS
$env:GOARCH=$env:OLDGOARCH