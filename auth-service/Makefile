include config/.env

.PHONY: gen-proto
gen-proto:
	protoc --go_out=api/gen/go --go_opt=paths=source_relative \
       --go-grpc_out=api/gen/go --go-grpc_opt=paths=source_relative \
       api/proto/auth.proto

.PHONY: run
run:
	CONFIG_PATH=config/.env go run cmd/auth/main.go

.PHONY: migrate
migrate:
	go run ./cmd/migrator --dsn=${STORAGE_PATH} --migrations-path=./migrations