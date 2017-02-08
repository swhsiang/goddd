package tracking

import (
	"time"

	"github.com/go-kit/kit/metrics"

	kitinfluxdb "github.com/go-kit/kit/metrics/influx"
	stdinfluxdb "github.com/influxdata/influxdb/client/v2"
)

// Store manage metrics
type Store struct {
	*kitinfluxdb.Influx
	client *stdinfluxdb.Client
}

// SaveMetrics send metrics to influxdb.
func (store *Store) SaveMetrics() {
	// NOTE: After send metrics to influxdb, the counter will be reset
	store.Influx.WriteTo(*store.client)
}

// NewStore returns an instance of Store struct
func NewStore(influx *kitinfluxdb.Influx, client *stdinfluxdb.Client) *Store {
	return &Store{
		Influx: influx,
		client: client,
	}
}

type instrumentingService struct {
	requestCount   metrics.Counter
	requestLatency metrics.Histogram
	Service
	*Store
}

// NewInstrumentingService returns an instance of an instrumenting Service.
func NewInstrumentingService(counter metrics.Counter, latency metrics.Histogram, s Service, store *Store) Service {
	return &instrumentingService{
		requestCount:   counter,
		requestLatency: latency,
		Service:        s,
		Store:          store,
	}
}

func (s *instrumentingService) Track(id string) (Cargo, error) {
	defer func(begin time.Time) {
		s.requestCount.With("method", "track").Add(1)
		s.requestLatency.With("method", "track").Observe(time.Since(begin).Seconds())
		s.Store.SaveMetrics()
	}(time.Now())

	return s.Service.Track(id)
}
