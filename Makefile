.PHONY: dev test tidy proto web-install web-dev web-build

dev:
	go run ./cmd/server

test:
	go test ./...

tidy:
	go mod tidy

proto:
	protoc \
		--go_out=. --go_opt=paths=source_relative \
		--go-grpc_out=. --go-grpc_opt=paths=source_relative \
		proto/external/demo/v1/demo.proto \
		proto/public/todo/v1/todo.proto

web-install:
	cd web && npm install

web-dev:
	cd web && npm run dev

web-build:
	cd web && npm run build
