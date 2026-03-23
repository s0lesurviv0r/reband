.PHONY: build

build:
	go build -o build/channel-conv main.go

test:
	go test -v ./...

clean:
	rm -rf build

fmt:
	gofmt -w **/*.go
