package obsmetrics

import (
	"context"
	"fmt"
	"net/http"
	"runtime"
	"time"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

// Metrics holds all Prometheus metrics
type Metrics struct {
	// System metrics
	goroutines  prometheus.Gauge
	memoryAlloc prometheus.Gauge
	memoryHeap  prometheus.Gauge
	cpuUsage    prometheus.Gauge

	// Application metrics
	httpRequestsTotal    *prometheus.CounterVec
	httpRequestDuration  *prometheus.HistogramVec
	httpRequestsInFlight prometheus.Gauge

	// IoT specific metrics
	iotMessagesReceived    prometheus.Counter
	iotMessagesProcessed   prometheus.Counter
	iotMessagesFailed      prometheus.Counter
	iotBatchProcessingTime *prometheus.HistogramVec
	iotBatchSize           *prometheus.HistogramVec
	iotDeviceCount         prometheus.Gauge
	iotSensorReadings      prometheus.Counter

	// Kafka metrics
	kafkaMessagesConsumed prometheus.Counter
	kafkaConsumerLag      prometheus.Gauge
	kafkaProducerMessages prometheus.Counter
	kafkaProducerErrors   prometheus.Counter

	// Database metrics
	dbConnectionsActive prometheus.Gauge
	dbConnectionsIdle   prometheus.Gauge
	dbQueryDuration     *prometheus.HistogramVec
	dbTransactionsTotal prometheus.Counter
	dbErrorsTotal       prometheus.Counter

	// Business metrics
	anomaliesDetected  prometheus.Counter
	alertsGenerated    prometheus.Counter
	deviceOnlineStatus *prometheus.GaugeVec
}

// NewMetrics creates and registers all Prometheus metrics
func NewMetrics(namespace string) *Metrics {
	m := &Metrics{
		// System metrics
		goroutines: promauto.NewGauge(prometheus.GaugeOpts{
			Namespace: namespace,
			Name:      "goroutines_total",
			Help:      "Number of goroutines currently running",
		}),
		memoryAlloc: promauto.NewGauge(prometheus.GaugeOpts{
			Namespace: namespace,
			Name:      "memory_alloc_bytes",
			Help:      "Current memory allocation in bytes",
		}),
		memoryHeap: promauto.NewGauge(prometheus.GaugeOpts{
			Namespace: namespace,
			Name:      "memory_heap_bytes",
			Help:      "Current heap memory usage in bytes",
		}),
		cpuUsage: promauto.NewGauge(prometheus.GaugeOpts{
			Namespace: namespace,
			Name:      "cpu_usage_percent",
			Help:      "Current CPU usage percentage",
		}),

		// HTTP metrics
		httpRequestsTotal: promauto.NewCounterVec(prometheus.CounterOpts{
			Namespace: namespace,
			Name:      "http_requests_total",
			Help:      "Total number of HTTP requests",
		}, []string{"method", "endpoint", "status"}),
		httpRequestDuration: promauto.NewHistogramVec(prometheus.HistogramOpts{
			Namespace: namespace,
			Name:      "http_request_duration_seconds",
			Help:      "HTTP request duration in seconds",
			Buckets:   prometheus.DefBuckets,
		}, []string{"method", "endpoint"}),
		httpRequestsInFlight: promauto.NewGauge(prometheus.GaugeOpts{
			Namespace: namespace,
			Name:      "http_requests_in_flight",
			Help:      "Number of HTTP requests currently being processed",
		}),

		// IoT metrics
		iotMessagesReceived: promauto.NewCounter(prometheus.CounterOpts{
			Namespace: namespace,
			Name:      "iot_messages_received_total",
			Help:      "Total number of IoT messages received",
		}),
		iotMessagesProcessed: promauto.NewCounter(prometheus.CounterOpts{
			Namespace: namespace,
			Name:      "iot_messages_processed_total",
			Help:      "Total number of IoT messages processed successfully",
		}),
		iotMessagesFailed: promauto.NewCounter(prometheus.CounterOpts{
			Namespace: namespace,
			Name:      "iot_messages_failed_total",
			Help:      "Total number of IoT messages that failed processing",
		}),
		iotBatchProcessingTime: promauto.NewHistogramVec(prometheus.HistogramOpts{
			Namespace: namespace,
			Name:      "iot_batch_processing_duration_seconds",
			Help:      "Time taken to process IoT message batches",
			Buckets:   prometheus.DefBuckets,
		}, []string{"batch_size"}),
		iotBatchSize: promauto.NewHistogramVec(prometheus.HistogramOpts{
			Namespace: namespace,
			Name:      "iot_batch_size",
			Help:      "Size of IoT message batches",
			Buckets:   []float64{1, 10, 50, 100, 500, 1000, 5000, 10000},
		}, []string{"status"}),
		iotDeviceCount: promauto.NewGauge(prometheus.GaugeOpts{
			Namespace: namespace,
			Name:      "iot_devices_total",
			Help:      "Total number of IoT devices",
		}),
		iotSensorReadings: promauto.NewCounter(prometheus.CounterOpts{
			Namespace: namespace,
			Name:      "iot_sensor_readings_total",
			Help:      "Total number of sensor readings processed",
		}),

		// Kafka metrics
		kafkaMessagesConsumed: promauto.NewCounter(prometheus.CounterOpts{
			Namespace: namespace,
			Name:      "kafka_messages_consumed_total",
			Help:      "Total number of Kafka messages consumed",
		}),
		kafkaConsumerLag: promauto.NewGauge(prometheus.GaugeOpts{
			Namespace: namespace,
			Name:      "kafka_consumer_lag",
			Help:      "Current Kafka consumer lag",
		}),
		kafkaProducerMessages: promauto.NewCounter(prometheus.CounterOpts{
			Namespace: namespace,
			Name:      "kafka_producer_messages_total",
			Help:      "Total number of Kafka messages produced",
		}),
		kafkaProducerErrors: promauto.NewCounter(prometheus.CounterOpts{
			Namespace: namespace,
			Name:      "kafka_producer_errors_total",
			Help:      "Total number of Kafka producer errors",
		}),

		// Database metrics
		dbConnectionsActive: promauto.NewGauge(prometheus.GaugeOpts{
			Namespace: namespace,
			Name:      "db_connections_active",
			Help:      "Number of active database connections",
		}),
		dbConnectionsIdle: promauto.NewGauge(prometheus.GaugeOpts{
			Namespace: namespace,
			Name:      "db_connections_idle",
			Help:      "Number of idle database connections",
		}),
		dbQueryDuration: promauto.NewHistogramVec(prometheus.HistogramOpts{
			Namespace: namespace,
			Name:      "db_query_duration_seconds",
			Help:      "Database query duration in seconds",
			Buckets:   prometheus.DefBuckets,
		}, []string{"operation", "table"}),
		dbTransactionsTotal: promauto.NewCounter(prometheus.CounterOpts{
			Namespace: namespace,
			Name:      "db_transactions_total",
			Help:      "Total number of database transactions",
		}),
		dbErrorsTotal: promauto.NewCounter(prometheus.CounterOpts{
			Namespace: namespace,
			Name:      "db_errors_total",
			Help:      "Total number of database errors",
		}),

		// Business metrics
		anomaliesDetected: promauto.NewCounter(prometheus.CounterOpts{
			Namespace: namespace,
			Name:      "anomalies_detected_total",
			Help:      "Total number of anomalies detected",
		}),
		alertsGenerated: promauto.NewCounter(prometheus.CounterOpts{
			Namespace: namespace,
			Name:      "alerts_generated_total",
			Help:      "Total number of alerts generated",
		}),
		deviceOnlineStatus: promauto.NewGaugeVec(prometheus.GaugeOpts{
			Namespace: namespace,
			Name:      "device_online_status",
			Help:      "Online status of IoT devices (1=online, 0=offline)",
		}, []string{"device_id", "device_type"}),
	}

	return m
}

// StartMetricsServer starts the Prometheus metrics HTTP server
func (m *Metrics) StartMetricsServer(ctx context.Context, addr string) error {
	http.Handle("/metrics", promhttp.Handler())

	server := &http.Server{
		Addr:    addr,
		Handler: http.DefaultServeMux,
	}

	go func() {
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			fmt.Printf("Metrics server error: %v\n", err)
		}
	}()

	// Start system metrics collection
	go m.collectSystemMetrics(ctx)

	return nil
}

// collectSystemMetrics periodically collects system-level metrics
func (m *Metrics) collectSystemMetrics(ctx context.Context) {
	ticker := time.NewTicker(15 * time.Second)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			return
		case <-ticker.C:
			// Update goroutine count
			m.goroutines.Set(float64(runtime.NumGoroutine()))

			// Update memory stats
			var mem runtime.MemStats
			runtime.ReadMemStats(&mem)
			m.memoryAlloc.Set(float64(mem.Alloc))
			m.memoryHeap.Set(float64(mem.HeapAlloc))

			// Note: CPU usage would require more complex implementation
			// For now, we'll use a placeholder
			m.cpuUsage.Set(0.0) // TODO: Implement actual CPU monitoring
		}
	}
}

// HTTP middleware methods
func (m *Metrics) HTTPMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()

		// Increment in-flight requests
		m.httpRequestsInFlight.Inc()
		defer m.httpRequestsInFlight.Dec()

		// Create a response writer wrapper to capture status code
		wrapped := &responseWriter{ResponseWriter: w, statusCode: 200}

		next.ServeHTTP(wrapped, r)

		// Record metrics
		duration := time.Since(start).Seconds()
		m.httpRequestDuration.WithLabelValues(r.Method, r.URL.Path).Observe(duration)
		m.httpRequestsTotal.WithLabelValues(r.Method, r.URL.Path, fmt.Sprintf("%d", wrapped.statusCode)).Inc()
	})
}

// IoT metrics methods
func (m *Metrics) RecordIOTMessageReceived() {
	m.iotMessagesReceived.Inc()
}

func (m *Metrics) RecordIOTMessageProcessed() {
	m.iotMessagesProcessed.Inc()
}

func (m *Metrics) RecordIOTMessageFailed() {
	m.iotMessagesFailed.Inc()
}

func (m *Metrics) RecordIOTBatchProcessing(batchSize int, duration time.Duration) {
	m.iotBatchProcessingTime.WithLabelValues(fmt.Sprintf("%d", batchSize)).Observe(duration.Seconds())
	m.iotBatchSize.WithLabelValues("processed").Observe(float64(batchSize))
}

func (m *Metrics) SetIOTDeviceCount(count int) {
	m.iotDeviceCount.Set(float64(count))
}

func (m *Metrics) RecordSensorReading() {
	m.iotSensorReadings.Inc()
}

// Kafka metrics methods
func (m *Metrics) RecordKafkaMessageConsumed() {
	m.kafkaMessagesConsumed.Inc()
}

func (m *Metrics) SetKafkaConsumerLag(lag int64) {
	m.kafkaConsumerLag.Set(float64(lag))
}

func (m *Metrics) RecordKafkaMessageProduced() {
	m.kafkaProducerMessages.Inc()
}

func (m *Metrics) RecordKafkaProducerError() {
	m.kafkaProducerErrors.Inc()
}

// Database metrics methods
func (m *Metrics) SetDBConnectionsActive(count int) {
	m.dbConnectionsActive.Set(float64(count))
}

func (m *Metrics) SetDBConnectionsIdle(count int) {
	m.dbConnectionsIdle.Set(float64(count))
}

func (m *Metrics) RecordDBQueryDuration(operation, table string, duration time.Duration) {
	m.dbQueryDuration.WithLabelValues(operation, table).Observe(duration.Seconds())
}

func (m *Metrics) RecordDBTransaction() {
	m.dbTransactionsTotal.Inc()
}

func (m *Metrics) RecordDBError() {
	m.dbErrorsTotal.Inc()
}

// Business metrics methods
func (m *Metrics) RecordAnomalyDetected() {
	m.anomaliesDetected.Inc()
}

func (m *Metrics) RecordAlertGenerated() {
	m.alertsGenerated.Inc()
}

func (m *Metrics) SetDeviceOnlineStatus(deviceID, deviceType string, online bool) {
	status := 0.0
	if online {
		status = 1.0
	}
	m.deviceOnlineStatus.WithLabelValues(deviceID, deviceType).Set(status)
}

// responseWriter wraps http.ResponseWriter to capture status code
type responseWriter struct {
	http.ResponseWriter
	statusCode int
}

func (rw *responseWriter) WriteHeader(code int) {
	rw.statusCode = code
	rw.ResponseWriter.WriteHeader(code)
}
