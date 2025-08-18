# Build 
FROM golang:1.23.2-alpine AS builder

RUN apk update && apk add --no-cache git

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN CGO_ENABLED=0 GOOS=linux go build -o server ./cmd/api

# Use the official Debian slim image for a lean production container.
FROM alpine:3.20

WORKDIR /app

COPY --from=builder /app/server /app/server

RUN apk --no-cache add ca-certificates tzdata curl
ENV TZ=Asia/Jakarta

EXPOSE 8080

ENTRYPOINT ["/app/server"]