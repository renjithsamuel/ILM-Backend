db-init:
	psql -c 'CREATE DATABASE "order-management-service"' -U $(user)
migrationup:
	migrate -path db/migration -database "postgres://$(DB_SERVICE_USER):$(DB_PASSWORD)@$(DB_HOST):$(DB_PORT)/order-management-service?sslmode=disable" -verbose up
migrationdown:
	migrate -path db/migration -database "postgres://$(user):$(password)@$(host):$(port)/order-management-service?sslmode=disable" -verbose down