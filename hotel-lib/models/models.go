package models

import (
	"encoding/json"
	"github.com/google/uuid"
	"time"
)

type PaymentInfo struct {
	Price   int     `json:"price"`
	Booking Booking `json:"booking"`
}

type Booking struct {
	Id         uuid.UUID  `json:"id"`
	HotelName  string     `json:"hotel_name"`
	RoomNumber int        `json:"room_number"`
	From       CustomDate `json:"from"`
	To         CustomDate `json:"to"`
}

type CustomDate struct {
	time.Time
}

type Hotels struct {
	Id    uuid.UUID `json:"id"`
	Name  string    `json:"name"`
	Rooms []Room    `json:"room"`
}

type Room struct {
	Id         uuid.UUID `json:"id,omitempty"`
	Price      int       `json:"price"`
	HotelId    uuid.UUID `json:"hotel_id,omitempty"`
	RoomNumber int       `json:"room_number"`
}

func (c *CustomDate) UnmarshalJSON(b []byte) (err error) {
	c.Time, err = time.Parse(`"2006-01-02"`, string(b))
	if err != nil {
		return err
	}
	return nil
}

func (c CustomDate) MarshalJSON() ([]byte, error) {
	formattedDate := c.Time.Format("2006-01-02")
	return json.Marshal(formattedDate)
}
