APP_GO_FILES := $(shell find . -name '*.go')

all: build windows.zip linux.zip

windows.zip: build/windows_amd64/clientfolders.exe LICENSE
	zip -r windows build/windows_amd64/clientfolders.exe LICENSE

linux.zip: build/linux_amd64/clientfolders LICENSE
	zip -r linux build/linux_amd64/clientfolders LICENSE

build: build/linux_amd64/clientfolders build/windows_amd64/clientfolders.exe

linux: build/linux_amd64/clientfolders
build/linux_amd64/clientfolders: $(APP_GO_FILES)
	go build -o ./build/linux_amd64/ ./cmd/...

windows: build/windows_amd64/clientfolders.exe
build/windows_amd64/clientfolders.exe: $(APP_GO_FILES)
	GOOS=windows GOARCH=amd64 go build -o ./build/windows_amd64/ ./cmd/...