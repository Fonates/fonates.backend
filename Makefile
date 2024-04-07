.PHONY: build dev prod
build:
	go build -o bin/server ./cmd/server/main.go
dev: 
	TONPROOF_PAYLOAD_SIGNATURE_KEY="secret" go run ./cmd/server/main.go
prod:
	go mod tidy && NODE_ENV=production go build -o bin/server ./cmd/server/main.go && NODE_ENV=production ./bin/server