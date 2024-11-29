FROM golang:alpine AS builder

WORKDIR /bookingservice
COPY ./go.mod ./go.sum ./
RUN go mod download

COPY bookingservice/ ./bookingservice/
COPY hotel-lib/ ./hotel-lib/
COPY protos/ ./protos/
COPY .env ./
RUN go build -o booking-service ./bookingservice/cmd/main.go
CMD ["./booking-service"]