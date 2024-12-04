# # Stage 1: Build
# FROM golang:1.23.3-alpine AS builder

# WORKDIR /app

# COPY go.mod go.sum ./

# RUN go mod download

# COPY . .

# RUN go build -o bookApp .

# Stage 2: Runtime
FROM alpine:latest

RUN mkdir /app

# COPY --from=builder /app/bookApp /app/
COPY bookApp /app/

CMD ["/app/bookApp"]
