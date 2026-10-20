.PHONY: all build test dev plan apply destroy seed clean check

all: build test

build:
	@mkdir -p bin
	GOOS=linux GOARCH=arm64 CGO_ENABLED=0 go build -trimpath -ldflags="-s -w" -o bin/presign ./cmd/presign
	GOOS=linux GOARCH=arm64 CGO_ENABLED=0 go build -trimpath -ldflags="-s -w" -o bin/intake ./cmd/intake
	GOOS=linux GOARCH=arm64 CGO_ENABLED=0 go build -trimpath -ldflags="-s -w" -o bin/analysis ./cmd/analysis
	GOOS=linux GOARCH=arm64 CGO_ENABLED=0 go build -trimpath -ldflags="-s -w" -o bin/caregiver_list ./cmd/caregiver_list
	GOOS=linux GOARCH=arm64 CGO_ENABLED=0 go build -trimpath -ldflags="-s -w" -o bin/caregiver_status ./cmd/caregiver_status
	GOOS=linux GOARCH=arm64 CGO_ENABLED=0 go build -trimpath -ldflags="-s -w" -o bin/digest ./cmd/digest
	@if [ -d frontend ]; then cd frontend && npm run build; fi

test:
	go test -v ./...
	@if [ -d frontend ]; then cd frontend && npm test; fi

dev:
	go run ./cmd/dev

plan:
	terraform -chdir=infra plan

apply:
	terraform -chdir=infra apply -auto-approve

destroy:
	terraform -chdir=infra destroy -auto-approve

seed:
	go run ./cmd/seed

check:
	gofmt -s -w .
	go vet ./...
	@if [ -d frontend ]; then cd frontend && npx svelte-check; fi

clean:
	rm -rf bin dist frontend/build
