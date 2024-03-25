build:
	@go build -o bin/main main.go

run: build
	@go run main.go

test:
	@go test -v ./...

seed:
	@go run scripts/seed.go