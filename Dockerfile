# Builder Stage
FROM golang:1.23.4 as builder

WORKDIR /app

# Copy files
COPY go.mod go.sum ./
RUN go mod download
COPY . .

# Inject BuildDate Metadata
RUN go build -ldflags="-X internal/version.BuildDate=$(date -u +%Y-%m-%dT%H:%M:%SZ)" -o ldcapi ./cmd
# Final Stage
FROM debian:bookworm-slim
# Install CA certificates
RUN apt-get update && apt-get install -y ca-certificates && rm -rf /var/lib/apt/lists/*

WORKDIR /app

COPY --from=builder /app/ldcapi .

RUN chmod +x ./ldcapi

EXPOSE 8080

CMD ["./ldcapi"]
# CMD ["tail", "-f", "/dev/null"]
