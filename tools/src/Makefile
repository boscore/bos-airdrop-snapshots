INSTANCE = bosa
APP_PATH = ./cli/bosa

default: install

pre: 
	go get -v github.com/Masterminds/glide

deps: pre
	glide install

rm-deps:
	rm -rf vendor

build:
	CGO_ENABLED=0 go build -a -installsuffix cgo ./cli/...

install:
	CGO_ENABLED=0 go install -a -installsuffix cgo ./cli/...

publish:
	mkdir -p bin
	CGO_ENABLED=0 GOOS=darwin GOARCH=amd64 go build -o bin/$(INSTANCE)-darwin $(APP_PATH)/main.go
	CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o bin/$(INSTANCE)-linux $(APP_PATH)/main.go 
	CGO_ENABLED=0 GOOS=windows GOARCH=amd64 go build -o bin/$(INSTANCE)-windows $(APP_PATH)/main.go
	chmod +x bin/*
