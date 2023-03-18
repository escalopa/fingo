# General Commands
proto:
	rm -rf pb/*
	protoc  --proto_path=proto --go_out=pb --go_opt=paths=source_relative \
	--go-grpc_out=pb --go-grpc_opt=paths=source_relative proto/*.proto

cert:
	cd certs && ./gen.sh && cd ..

# AUTH Command
auth-mig:
	migrate create -ext sql -dir ./auth/internal/adapters/db/postgres/migrations -seq $(SCHEMA_NAME)

auth-sqlc:
	sqlc generate -f ./auth/sqlc.yaml

.PHONE: proto cert auth-mig auth-sqlc
