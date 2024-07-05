build:
	@go build -o bin/complete-auth cmd/app/main.go

run: build
	@./bin/complete-auth

test:
	@go test -v ./...
