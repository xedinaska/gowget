all: build

lint:
	gofmt -w .
	go vet ./...
	golint $$(ls -d */ | grep -v vendor)
	gocyclo -over 10  $$(ls -d */ | grep -v vendor)

build:
	go build -o ./bin/gowget ./cmd/main.go
