FROM golang:1.21-alpine

WORKDIR /app
COPY . .
RUN go mod init mock_external_service && go mod tidy
RUN go build -o mockservice .

EXPOSE 8050
CMD ["./mockservice"]
