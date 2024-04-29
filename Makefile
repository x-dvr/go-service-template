build:
	@go build -o bin/api cmd/api/main.go

run: build
	@./bin/api

test:
	@go test -v ./...

bench:
	@go test -bench=. -benchmem ./...
