# --- Stage 1: Test (Optional) ---
FROM golang:1.24-alpine AS test
WORKDIR /app
# Use module files from backend
COPY backend/go.mod backend/go.sum ./
RUN go mod download
COPY backend/ ./
# Adjust test command to run backend tests (optional)
RUN go test ./test/... || true

# --- Stage 2: Build Go Application ---
FROM golang:1.24-alpine AS go-builder
ENV GOMODCACHE=/go/pkg/mod
WORKDIR /app
# Copy module files and download dependencies
COPY backend/go.mod backend/go.sum ./
RUN go mod download
# Copy backend source
COPY backend/ ./
# Build the main application from backend/cmd
RUN GOOS=linux GOARCH=amd64 CGO_ENABLED=0 go build -ldflags="-s -w" -o main ./cmd && chmod +x main

# --- Stage 3: Minimal Final Image ---
FROM alpine:3.19
RUN adduser -D -h /home/appuser appuser
USER appuser
WORKDIR /home/appuser
COPY --from=go-builder /app/main .
EXPOSE 8080
ENTRYPOINT ["/home/appuser/main"]
