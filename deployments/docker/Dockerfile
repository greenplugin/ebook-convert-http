# === Stage 1: Build application ===
FROM golang:1.23 AS builder

WORKDIR /app

COPY . /app

RUN go mod download

RUN CGO_ENABLED=0 go build -a -ldflags '-extldflags "-static"' -o ebook-convert-wrapper ./cmd


FROM ubuntu:latest

# Install dependent packages and Calibre
RUN apt-get update && \
    apt-get install -y --no-install-recommends \
    calibre && \
    rm -rf /var/lib/apt/lists/*

# Copy the compiled application from the build stage
COPY --from=builder /app/ebook-convert-wrapper /usr/local/bin/ebook-convert-wrapper

# Expose a port for operation (e.g., 8080)
EXPOSE 8080

# Set the default startup command
CMD ["ebook-convert-wrapper"]
