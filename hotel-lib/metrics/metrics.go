package metrics

import "github.com/prometheus/client_golang/prometheus"

type Metrics struct {
	requestAllHotels  *prometheus.CounterVec
	requestAddHotels  *prometheus.CounterVec
	requestAddBooking *prometheus.CounterVec
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

	return &Metrics{
		requestAllHotels:  requestAllHotels,
		requestAddHotels:  requestAddHotels,
		requestAddBooking: requestAddBooking,
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
