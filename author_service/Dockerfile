# # Stage 1: Build
# FROM golang:1.23.3-alpine AS builder

# WORKDIR /app

# COPY go.mod go.sum ./

# RUN go mod download

# COPY . .

# RUN go build -o authApp .

# Stage 2: Runtime
FROM alpine:latest

RUN mkdir /app

# COPY --from=builder /app/authApp /app/
COPY authApp /app/

CMD ["/app/authApp"]
