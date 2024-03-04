.PHONY: migrate
migrate:
	psql -U postgres -d postgres -h localhost -a -f ./migrations/create_table.sql
	psql -U postgres -d postgres -h localhost -a -f ./migrations/insert_data.sql

CLICKHOUSE_CLIENT = clickhouse-client
CLICKHOUSE_HOST = localhost
CLICKHOUSE_DATABASE = your_database

.PHONY: migrate_clickhouse
migrate_clickhouse:$(CLICKHOUSE_CLIENT) --host $(CLICKHOUSE_HOST) --database $(CLICKHOUSE_DATABASE) < ./migrations/clickhouse.sql

serverStart:
	go run ../hezzlizer/cmd/server/main.go

clientStart:
	go run ../hezzlizer/cmd/client/main.go

