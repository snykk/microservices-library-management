# # Stage 1: Build
# FROM golang:1.23.3-alpine AS builder

# WORKDIR /app

# COPY go.mod go.sum ./

# RUN go mod download

# COPY . .

# RUN go build -o authorApp .

# Stage 2: Runtime
FROM alpine:latest

RUN GRPC_HEALTH_PROBE_VERSION=v0.4.13 && \
    wget -qO/bin/grpc_health_probe https://github.com/grpc-ecosystem/grpc-health-probe/releases/download/${GRPC_HEALTH_PROBE_VERSION}/grpc_health_probe-linux-amd64 && \
    chmod +x /bin/grpc_health_probe
    
RUN mkdir /app

# COPY --from=builder /app/authorApp /app/
COPY authorApp /app/

CMD ["/app/authorApp"]
