.PHONY: build dev prod
build:
	go build -o bin/server ./cmd/server/main.go
dev: 
	go run ./cmd/server/main.go
prod:
	go mod tidy && NODE_ENV=production go build -o bin/server ./cmd/server/main.go && ./bin/server