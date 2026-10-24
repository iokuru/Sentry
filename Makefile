export GOOS := linux
export GOARCH := arm64
export CGO_ENABLED := 0

.PHONY: all build test dev plan apply destroy seed clean check

all: build test

build:
	go build -trimpath -ldflags="-s -w" -o bin/presign ./cmd/presign
	go build -trimpath -ldflags="-s -w" -o bin/intake ./cmd/intake
	go build -trimpath -ldflags="-s -w" -o bin/analysis ./cmd/analysis
	go build -trimpath -ldflags="-s -w" -o bin/caregiver_list ./cmd/caregiver_list
	go build -trimpath -ldflags="-s -w" -o bin/caregiver_status ./cmd/caregiver_status
	go build -trimpath -ldflags="-s -w" -o bin/digest ./cmd/digest
	cd frontend && npm run build

test:
	go test -v ./...
	cd frontend && npm test

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
	cd frontend && npx svelte-check

clean:
	rm -rf bin dist frontend/build
