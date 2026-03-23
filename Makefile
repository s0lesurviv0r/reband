.PHONY: build

build:
	go build -o build/reband main.go

test:
	go test -v ./...

clean:
	rm -rf build

fmt:
	gofmt -w **/*.go
