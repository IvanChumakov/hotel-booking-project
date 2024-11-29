FROM golang:alpine AS builder

WORKDIR /paymentservice
COPY ./go.mod ./go.sum ./
RUN go mod download

COPY paymentservice/ ./paymentservice/
COPY hotel-lib/ ./hotel-lib/
COPY protos/ ./protos/
COPY .env ./
RUN go build -o payment-service ./paymentservice/cmd/main.go
CMD ["./payment-service"]