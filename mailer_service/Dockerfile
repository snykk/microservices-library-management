# # Stage 1: Build
# FROM golang:1.23.3-alpine AS builder

# WORKDIR /app

# COPY go.mod go.sum ./

# RUN go mod download

# COPY . .

# RUN go build -o mailerApp .

# Stage 2: Runtime
FROM alpine:latest

RUN apk add --no-cache curl

RUN mkdir /app

# COPY --from=builder /app/mailerApp /app/
COPY mailerApp /app/

COPY templates /templates

CMD ["/app/mailerApp"]
