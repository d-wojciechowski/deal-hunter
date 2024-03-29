$env:PRODUCT_NAME="DealHunter"
Write-Output("Building $env:PRODUCT_NAME`n")

Remove-Item –path "distr\*"

Write-Output("Current GOOS $env:GOOS : Current GOARCH :$env:GOARCH.`n")
$env:OLDGOOS=$env:GOOS
$env:OLDGOARCH=$env:GOARCH

$env:GOARCH="amd64"
$env:GOOS="linux"
Write-Output("Linux x64 build| GOOS $env:GOOS : GOARCH :$env:GOARCH.")
go build -ldflags="-w -s" -o distr/$env:PRODUCT_NAME-linux-x64
Write-Output("Linux build done.`n")

$env:GOARCH="amd64"
$env:GOOS="windows"
Write-Output("Windows x64 build| GOOS $env:GOOS : GOARCH :$env:GOARCH.")
go build -o distr/$env:PRODUCT_NAME-windows-x64.exe
Write-Output("Windows x64 build done.`n")

$env:GOOS=$env:OLDGOOS
$env:GOARCH=$env:OLDGOARCH

upx -9 distr/DealHunter-linux-x64
upx -9 distr/DealHunter-windows-x64.exe

docker build -t dwojciechowski/dealhunter:latest .
docker save --output deal-hunter.tar.gz dwojciechowski/dealhunter:latest
