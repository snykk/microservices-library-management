# # Stage 1: Build
# FROM golang:1.23.3-alpine AS builder

# WORKDIR /app

# COPY go.mod go.sum ./

# RUN go mod download

# COPY . .

# RUN go build -o gatewayApp .

# Stage 2: Runtime
FROM alpine:latest

RUN apk add --no-cache tzdata curl

RUN cp /usr/share/zoneinfo/Asia/Jakarta /etc/localtime && echo "Asia/Jakarta" > /etc/timezone

RUN mkdir /app

# COPY --from=builder /app/gatewayApp /app/
COPY gatewayApp /app/

CMD ["/app/gatewayApp"]
