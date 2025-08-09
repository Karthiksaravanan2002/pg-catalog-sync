# Build stage
FROM golang:1.21-alpine AS build
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN go build -o /pg-catalog-sync ./cmd/server

# Run stage
FROM alpine:3.18
RUN apk add --no-cache ca-certificates
COPY --from=build /pg-catalog-sync /usr/local/bin/pg-catalog-sync
EXPOSE 8080
ENTRYPOINT ["/usr/local/bin/pg-catalog-sync"]
