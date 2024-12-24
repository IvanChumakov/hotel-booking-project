package metrics

import "github.com/prometheus/client_golang/prometheus"

type Metrics struct {
	requestAllHotels          *prometheus.CounterVec
	requestAddHotels          *prometheus.CounterVec
	requestAddBooking         *prometheus.CounterVec
	requestGetBookings        *prometheus.CounterVec
	requestGetBookingsByName  *prometheus.CounterVec
	requestGetFreeRoomsByDate *prometheus.CounterVec
}

func NewMetrics() *Metrics {
	requestAllHotels := prometheus.NewCounterVec(prometheus.CounterOpts{
		Name: "request_all_hotels",
		Help: "Кол-во запросов для получения всех отелей",
	}, []string{})
	prometheus.MustRegister(requestAllHotels)

	requestAddHotels := prometheus.NewCounterVec(prometheus.CounterOpts{
		Name: "request_add_hotels",
		Help: "Кол-во запросов на добавление отелей",
	}, []string{})
	prometheus.MustRegister(requestAddHotels)

	requestAddBooking := prometheus.NewCounterVec(prometheus.CounterOpts{
		Name: "request_add_booking",
		Help: "Кол-во запросов на бронирование отеля",
	}, []string{})
	prometheus.MustRegister(requestAddBooking)

	requestGetBookings := prometheus.NewCounterVec(prometheus.CounterOpts{
		Name: "request_get_bookings",
		Help: "Кол-во запросов на получение бронирований",
	}, []string{})
	prometheus.MustRegister(requestGetBookings)

	requestGetBookingsByName := prometheus.NewCounterVec(prometheus.CounterOpts{
		Name: "request_get_bookings_by_name",
		Help: "Кол-во запросов на получение бронирований по имени отеля",
	}, []string{})

	requestGetFreeRoomsByDate := prometheus.NewCounterVec(prometheus.CounterOpts{
		Name: "request_get_free_rooms_by_date",
		Help: "Кол-во запросов на получение свободных комнат по дате и отелю",
	}, []string{})

	return &Metrics{
		requestAllHotels:          requestAllHotels,
		requestAddHotels:          requestAddHotels,
		requestAddBooking:         requestAddBooking,
		requestGetBookings:        requestGetBookings,
		requestGetBookingsByName:  requestGetBookingsByName,
		requestGetFreeRoomsByDate: requestGetFreeRoomsByDate,
	}
}

func (m *Metrics) IncRequestAllHotels() {
	m.requestAllHotels.WithLabelValues().Inc()
}

func (m *Metrics) IncRequestAddHotels() {
	m.requestAddHotels.WithLabelValues().Inc()
}

func (m *Metrics) IncRequestAddBooking() {
	m.requestAddBooking.WithLabelValues().Inc()
}

func (m *Metrics) IncRequestGetBookings() {
	m.requestGetBookings.WithLabelValues().Inc()
}

func (m *Metrics) IncRequestGetBookingsByName() {
	m.requestGetBookingsByName.WithLabelValues().Inc()
}

func (m *Metrics) IncRequestGetFreeRoomsByDate() {
	m.requestGetFreeRoomsByDate.WithLabelValues().Inc()
}
