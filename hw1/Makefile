.PHONY: build run demo test race vet check coverage clean

build:
	mkdir -p bin
	go build -o bin/vyshkat .

run:
	go run .

demo:
	go run . --demo

test:
	go test ./...

race:
	go test -race ./...

vet:
	go vet ./...

check: test race vet build

coverage:
	go test ./... -coverprofile=coverage.out
	go tool cover -func=coverage.out

clean:
	rm -f bin/vyshkat coverage.out
