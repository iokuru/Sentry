param(
    [string]$Command = "dev"
)

$ErrorActionPreference = "Stop"

switch ($Command) {
    "build" {
        go run ./cmd/build -frontend
    }
    "build-lambdas" {
        go run ./cmd/build
    }
    "test" {
        go test ./...
        Push-Location frontend
        npm test
        Pop-Location
    }
    "dev" {
        go run ./cmd/devserver
    }
    "seed" {
        go run ./cmd/seed
    }
    "smoke" {
        go run ./cmd/smoke
    }
    "smoke-live" {
        go run ./cmd/smoke -live
    }
    "check" {
        gofmt -s -w .
        go vet ./...
        Push-Location frontend
        npx svelte-check
        Pop-Location
    }
    "telegram-hook" {
        go run ./cmd/telegram_set_webhook
    }
    "deploy" {
        go run ./cmd/costguard
        go run ./cmd/deploy
    }
    "preflight" {
        go run ./cmd/preflight
    }
    "go-live" {
        go run ./cmd/golive
    }
    "status" {
        go run ./cmd/status
    }
    "costguard" {
        go run ./cmd/costguard
    }
    "clean" {
        go run ./cmd/build -clean
    }
    default {
        Write-Error "Unknown command: $Command. Available: build, build-lambdas, test, dev, seed, smoke, check, clean"
    }
}
