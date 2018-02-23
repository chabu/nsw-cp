NAME = nsw-cp
GOFLAGS = -ldflags "-s -w"

diversity: windows mac linux freebsd

linux:
	CGO_ENABLED=0 GOARCH=amd64 GOOS=linux go build $(GOFLAGS) -o $(NAME)_Linux_x64

mac:
	CGO_ENABLED=0 GOARCH=amd64 GOOS=darwin go build $(GOFLAGS) -o $(NAME)_macOS_x64

freebsd:
	CGO_ENABLED=0 GOARCH=amd64 GOOS=freebsd go build $(GOFLAGS) -o $(NAME)_FreeBSD_x64

windows:
	CGO_ENABLED=0 GOARCH=amd64 GOOS=windows go build $(GOFLAGS) -o $(NAME)_Windows_x64.exe
