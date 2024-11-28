FROM golang:alpine

WORKDIR /hotelservice
COPY . .
RUN go mod download
CMD ["go", "run", "cmd/HotelService/main.go"]