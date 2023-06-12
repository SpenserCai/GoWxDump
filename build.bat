go mod tidy
set CGO_ENABLED=1
set GOARCH=386
go build -ldflags "-w -s"
copy GoWxDump.exe Release\GoWxDump.exe