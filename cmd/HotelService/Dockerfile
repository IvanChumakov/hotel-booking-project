FROM golang:alpine

WORKDIR /hotelservice
COPY . .
RUN go mod download
RUN go build -o hotelservice cmd/HotelService/main.go
CMD ["./hotelservice"]