FROM golang:alpine

WORKDIR /bookingservice
COPY ./bookingservice .
COPY ./protos .
COPY ./hotel-lib .
COPY ./go.mod .
RUN go mod download
RUN go build -o bookingservice ./cmd/main.go
CMD ["./bookingservice"]