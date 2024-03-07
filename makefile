db-init:
	psql -c 'CREATE DATABASE "integrated-library-management-service"' -U $(user)
migrationup:
	migrate -path db/migration -database "postgres://SYS:aegleadmin@localhost:5432/integrated-library-management-service?sslmode=disable" -verbose up
	# migrate -path db/migration -database "postgres://$(DB_SERVICE_USER):$(DB_PASSWORD)@$(DB_HOST):$(DB_PORT)/order-management-service?sslmode=disable" -verbose up
migrationdown:
	migrate -path db/migration -database "postgres://$(user):$(password)@$(host):$(port)/order-management-service?sslmode=disable" -verbose down