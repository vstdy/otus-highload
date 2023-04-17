LOCAL_BIN:=$(CURDIR)/bin
PG_MIGRATIONS_DIR:=$(CURDIR)/storage/psql/migrations
PG_LOCAL_CREDS:="user=user password=password host=localhost port=5555 dbname=project sslmode=disable"

.PHONY: test
test:
	go test -v ./... --count=1

.PHONY: bin-deps
bin-deps:
	$(info Installing binary dependencies...)
	GOBIN=$(LOCAL_BIN) go install github.com/pressly/goose/v3/cmd/goose@v3.10.0

.PHONY: migrations
migrations: bin-deps
	$(LOCAL_BIN)/goose -dir $(PG_MIGRATIONS_DIR) postgres $(PG_LOCAL_CREDS) up
