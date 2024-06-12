build:
	@go build -o bin/complete-auth src/main.go

run: build
	@./bin/complete-auth

test:
	@go test -v ./...
