.PHONY: dev dev-backend dev-frontend dev-infra build test test-backend test-frontend migrate migrate-down lint

# all 
dev: dev-infra dev-backend dev-frontend

# postgres + mailhog
dev-infra:
	docker compose up -d

# hot reload
dev-backend:
	cd backend && air

dev-frontend:
	cd frontend && pnpm dev

build-backend:
	cd backend && go build -o server ./cmd/server

build-frontend:
	cd frontend && pnpm build

test: test-backend test-frontend

test-backend:
	cd backend && go test ./...

test-frontend:
	cd frontend && pnpm test

migrate:
	cd backend && migrate -path migrations -database "$${DATABASE_URL}" up

migrate-down:
	cd backend && migrate -path migrations -database "$${DATABASE_URL}" down 1

migrate-create:
	cd backend && migrate create -ext sql -dir migrations -seq $(name)

lint-backend:
	cd backend && go vet ./... && staticcheck ./...

lint-frontend:
	cd frontend && pnpm lint
