# --- Stage 1: Test (Optional) ---
FROM golang:1.24-alpine AS test
WORKDIR /app
COPY src/go.mod src/go.sum ./
RUN go mod download
COPY src/ .
RUN go test ./test/utils/...

# --- Stage 2: Build Go Application ---
FROM golang:1.24-alpine AS go-builder
ENV GOMODCACHE=/go/pkg/mod
WORKDIR /app
COPY src/go.mod src/go.sum ./
RUN go mod download
COPY src/ .
RUN GOOS=linux GOARCH=amd64 CGO_ENABLED=0 go build -ldflags="-s -w" -o main ./cmd/main.go && chmod +x main

# --- Stage 3: Minimal Final Image ---
FROM alpine:3.19
RUN adduser -D -h /home/appuser appuser
USER appuser
WORKDIR /home/appuser
COPY --from=go-builder /app/main .
EXPOSE 8080
ENTRYPOINT ["/home/appuser/main"]
