.PHONY: all build build-lambdas test dev plan apply destroy seed clean check smoke smoke-live telegram-hook deploy costguard preflight go-live status teardown

all: build test

preflight:
	go run ./cmd/preflight

go-live:
	go run ./cmd/golive

status:
	go run ./cmd/status

teardown:
	go run ./cmd/teardown

build:
	go run ./cmd/build -frontend

build-lambdas:
	go run ./cmd/build

test:
	go test ./...
	cd frontend && npm test

dev:
	go run ./cmd/devserver

deploy:
	go run ./cmd/costguard
	go run ./cmd/deploy

costguard:
	go run ./cmd/costguard

telegram-hook:
	go run ./cmd/telegram_set_webhook

plan:
	terraform -chdir=infra plan

apply:
	terraform -chdir=infra apply -auto-approve

destroy:
	terraform -chdir=infra destroy -auto-approve

seed:
	go run ./cmd/seed

smoke:
	go run ./cmd/smoke

smoke-live:
	go run ./cmd/smoke -live

check:
	gofmt -s -w .
	go vet ./...
	cd frontend && npm run check

clean:
	go run ./cmd/build -clean
