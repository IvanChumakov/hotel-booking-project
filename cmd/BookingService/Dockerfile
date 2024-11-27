FROM golang:alpine

WORKDIR /bookingservice
COPY . .
RUN go mod download
RUN go build -o bookingservice cmd/BookingService/main.go
CMD ["./bookingservice"]