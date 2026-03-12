FROM node:22-alpine AS frontend
WORKDIR /app/frontend
COPY frontend/package*.json ./
RUN npm ci
COPY frontend/ .
RUN npm run build

FROM golang:1.25-alpine AS backend
WORKDIR /app
COPY go.* ./
RUN go mod download
COPY . .
COPY --from=frontend /app/frontend/dist ./frontend/dist
RUN CGO_ENABLED=0 go build -o loomhub ./cmd/loomhub

FROM alpine:latest
COPY --from=backend /app/loomhub /loomhub
EXPOSE 8080
CMD ["/loomhub"]
