// Package metrics provides an abstraction around metrics collection
// in order to bundle all metrics related calls in one location
package metrics

import (
	"net/http"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

const (
	metricSecretsCreated      = "secrets_created"
	metricSecretsRead         = "secrets_read"
	metricSecretsCreateErrors = "secrets_create_errors"
	meticsSecretsReadErrors   = "secrets_read_errors"
	metricsSecretsStored      = "secrets_stored"

	labelReason = "reason"

	namespace = "ots"
)

type (
	// Collector contains all required methods to collect metrics
	// and to populate them into the Handler
	Collector struct {
		secretsCreated      prometheus.Counter
		secretsRead         prometheus.Counter
		secretsCreateErrors *prometheus.CounterVec
		secretsReadErrors   *prometheus.CounterVec
		secretsStored       prometheus.Gauge
	}
)

// Handler returns the handler to be registered at /metrics
func Handler() http.Handler { return promhttp.Handler() }

// New creates a new Collector and registers the metrics
func New() *Collector {
	return &Collector{
		secretsCreated: promauto.NewCounter(prometheus.CounterOpts{
			Namespace: namespace,
			Name:      metricSecretsCreated,
			Help:      "number of successfully created secrets",
		}),

		secretsRead: promauto.NewCounter(prometheus.CounterOpts{
			Namespace: namespace,
			Name:      metricSecretsRead,
			Help:      "number of fetched (and destroyed) secrets",
		}),

		secretsCreateErrors: promauto.NewCounterVec(prometheus.CounterOpts{
			Namespace: namespace,
			Name:      metricSecretsCreateErrors,
			Help:      "number of errors on secret creation for each reason",
		}, []string{labelReason}),

		secretsReadErrors: promauto.NewCounterVec(prometheus.CounterOpts{
			Namespace: namespace,
			Name:      meticsSecretsReadErrors,
			Help:      "number of read-errors for each reason",
		}, []string{labelReason}),

		secretsStored: promauto.NewGauge(prometheus.GaugeOpts{
			Namespace: namespace,
			Name:      metricsSecretsStored,
			Help:      "number of secrets currently held in the backend store",
		}),
	}
}

// CountSecretCreated signalizes a secret has successfully been created
func (c Collector) CountSecretCreated() { c.secretsCreated.Inc() }

// CountSecretRead signalizes a secret has successfully been read and destroyed
func (c Collector) CountSecretRead() { c.secretsRead.Inc() }

// CountSecretCreateError signalizes an error occurred during secret
// creation. The reason must not be the error.Error() but a simple
// static string describing the error.
func (c Collector) CountSecretCreateError(reason string) {
	c.secretsCreateErrors.WithLabelValues(reason).Inc()
}

// CountSecretReadError signalizes an error occurred during secret
// read. The reason must not be the error.Error() but a simple
// static string describing the error.
func (c Collector) CountSecretReadError(reason string) {
	c.secretsReadErrors.WithLabelValues(reason).Inc()
}

// UpdateSecretsCount sets the current amount of secrets stored in the
// backend storage
func (c Collector) UpdateSecretsCount(count int64) {
	c.secretsStored.Set(float64(count))
}
