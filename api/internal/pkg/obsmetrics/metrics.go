package obsmetrics

import (
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
)

var (
	// Batch processing metrics
	BatchDuration = promauto.NewHistogram(
		prometheus.HistogramOpts{
			Name:    "iot_batch_processing_duration_seconds",
			Help:    "Time taken to process a batch of messages",
			Buckets: []float64{.005, .01, .025, .05, .1, .25, .5, 1},
		},
	)

	// Database metrics
	DBOperationDuration = promauto.NewHistogramVec(
		prometheus.HistogramOpts{
			Name:    "iot_db_operation_duration_seconds",
			Help:    "Duration of database operations",
			Buckets: []float64{.005, .01, .025, .05, .1, .25, .5, 1},
		},
		[]string{"operation"}, // insert, query
	)

	ProcessingDuration = promauto.NewHistogram(prometheus.HistogramOpts{
		Name:    "iot_processing_duration_seconds",
		Help:    "Time taken to process messages",
		Buckets: []float64{.005, .01, .025, .05, .1, .25, .5, 1},
	})

	// Message processing metrics
	MessagesProcessed = promauto.NewCounterVec(
		prometheus.CounterOpts{
			Name: "iot_messages_processed_total",
			Help: "Total number of IoT messages processed",
		},
		[]string{"device_id", "status"}, // labels for device ID and success/error
	)

	// Batch processing metrics
	BatchSize = promauto.NewGauge(
		prometheus.GaugeOpts{
			Name: "iot_current_batch_size",
			Help: "Current batch size of IoT messages",
		},
	)

	ProcessingLatency = promauto.NewHistogramVec(
		prometheus.HistogramOpts{
			Name:    "iot_message_processing_seconds",
			Help:    "Time taken to process IoT messages",
			Buckets: []float64{.001, .005, .01, .025, .05, .1, .25, .5, 1},
		},
		[]string{"operation"}, // insert, query
	)

	// Sensor metrics
	SensorReadings = promauto.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "iot_sensor_reading",
			Help: "Current sensor readings",
		},
		[]string{"device_id", "type"}, // temperature, humidity, co2
	)
)
