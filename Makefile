APP_GO_FILES := $(shell find . -name '*.go')

build.zip: build
	zip build build
	
build: build/linux_amd64/clientfolders build/windows_amd64/clientfolders.exe

build/linux_amd64/clientfolders: $(APP_GO_FILES)
	go build -o ./build/linux_amd64/ ./cmd/...

build/windows_amd64/clientfolders.exe: $(APP_GO_FILES)
	GOOS=windows GOARCH=amd64 go build -o ./build/windows_amd64/ ./cmd/...