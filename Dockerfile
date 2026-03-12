FROM node:22-alpine AS frontend
WORKDIR /app/frontend
COPY frontend/package*.json ./
RUN npm ci
COPY frontend/ .
RUN npm run build

FROM golang:1.25-alpine AS builder
WORKDIR /app
# Cache buster: 20260312a
COPY go.* ./
RUN go mod download
COPY . .
COPY --from=frontend /app/frontend/dist ./frontend/dist
RUN CGO_ENABLED=0 go build -o main ./cmd/loomhub

FROM alpine:latest
WORKDIR /app
ARG CACHE_BUST=1
COPY --from=builder /app/main /app/main
EXPOSE 8080
CMD ["/app/main", "serve", "--port", "8080"]
