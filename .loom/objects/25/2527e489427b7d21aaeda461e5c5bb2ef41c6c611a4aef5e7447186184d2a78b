.PHONY: build dev test clean

# Full production build
build: frontend-build embed-frontend
	go build -o bin/loomhub ./cmd/loomhub

# Go binary only (no frontend)
build-go:
	go build -o bin/loomhub ./cmd/loomhub

# Build Vue SPA
frontend-build:
	cd frontend && npm install && npm run build

# Copy frontend dist into Go-embeddable location
embed-frontend:
	rm -rf internal/server/dist
	cp -r frontend/dist internal/server/dist

# Development — run Go backend
dev:
	go run ./cmd/loomhub serve --dev

# Run all Go tests
test:
	go test ./...

# Run tests with race detector
test-race:
	go test -race ./...

# Clean
clean:
	rm -rf bin/ data/ internal/server/dist/
