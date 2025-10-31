proto:
	rm -rf api/pb/*.go
	rm -rf doc/swagger/*.swagger.json
	protoc --proto_path=api/proto --go_out=api/pb --go_opt=paths=source_relative --go-grpc_out=api/pb --go-grpc_opt=paths=source_relative --grpc-gateway_out=api/pb --grpc-gateway_opt paths=source_relative --openapiv2_out=doc/swagger --openapiv2_opt=allow_merge=true,merge_file_name=webb_ai api/proto/*.proto

sqlc:
	sqlc generate

test:
	go test -v -cover -short ./...

mock:
	mockgen -destination db/mock/store.go -package mockdb github.com/jasonwebb3152/webb-ai/db/sqlc Store

new_migration:
	migrate create -ext sql -dir db/migration -seq $(name)

postgres:
	docker run --name postgres18 --network bank-network -p 5432:5432 -e POSTGRES_USER=root -e POSTGRES_PASSWORD=secret -d postgres:18-alpine
	docker exec -it postgres18 createdb --username=root --owner=root simple_bank

build_docker:
	docker build -t webb-ai-api:latest .

cluster_up:
	kind create cluster --config kind/webb-ai-cluster.yaml
	kind load docker-image webb-ai-api:latest --name webb-ai-cluster
	kubectl config use-context kind-webb-ai-cluster
	kubectl apply -f kind/db-deployment.yaml
	kubectl apply -f kind/db-service.yaml
	kubectl apply -f kind/api-deployment.yaml
	kubectl apply -f kind/api-service.yaml

cluster_down:
	kind delete cluster --name webb-ai-cluster

.PHONY: proto sqlc test new_migration build_docker cluster_up cluster_down
