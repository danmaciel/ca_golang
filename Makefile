createmigration:
	migrate create -ext=sql -dir=internal/infra/migrations -seq init

# Supply MIGRATE_DATABASE_URL through the environment for local migrations.
# Never commit a database password or connection URL.
migrate:
	@ : "$${MIGRATE_DATABASE_URL:?Set MIGRATE_DATABASE_URL before running make migrate}"
	@ migrate -path=internal/infra/migrations -database "$$MIGRATE_DATABASE_URL" -verbose up

migratedown:
	@ : "$${MIGRATE_DATABASE_URL:?Set MIGRATE_DATABASE_URL before running make migratedown}"
	@ migrate -path=internal/infra/migrations -database "$$MIGRATE_DATABASE_URL" -verbose down

.PHONY: migrate createmigration migratedown
