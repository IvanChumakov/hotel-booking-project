FROM golang:alpine AS builder

WORKDIR /hotelservice
COPY ./go.mod ./go.sum ./
RUN go mod download

COPY hotelservice/ ./hotelservice/
COPY hotel-lib/ ./hotel-lib/
COPY protos/ ./protos/
COPY .env ./
RUN go build -o hotel-service ./hotelservice/cmd/main.go
CMD ["./hotel-service"]