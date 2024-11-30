FROM golang:alpine AS builder

WORKDIR /deliveryservice
COPY ./go.mod ./go.sum ./
RUN go mod download

COPY deliveryservice/ ./deliveryservice/
COPY hotel-lib/ ./hotel-lib/
COPY protos/ ./protos/
COPY .env ./
RUN go build -o delivery-service ./deliveryservice/cmd/main.go
CMD ["./delivery-service"]