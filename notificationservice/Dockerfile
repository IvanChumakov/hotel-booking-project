FROM golang:alpine AS builder

WORKDIR /notificationservice
COPY ./go.mod ./go.sum ./
RUN go mod download

COPY notificationservice/ ./notificationservice/
COPY hotel-lib/ ./hotel-lib/
COPY protos/ ./protos/
COPY .env ./
RUN go build -o notification-service ./notificationservice/cmd/main.go
CMD ["./notification-service"]