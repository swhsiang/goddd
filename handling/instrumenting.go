package handling

import (
	"time"

	"github.com/go-kit/kit/metrics"

	"github.com/marcusolsson/goddd/cargo"
	"github.com/marcusolsson/goddd/location"
	"github.com/marcusolsson/goddd/voyage"

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

func (s *instrumentingService) RegisterHandlingEvent(completed time.Time, id cargo.TrackingID, voyageNumber voyage.Number,
	loc location.UNLocode, eventType cargo.HandlingEventType) error {

	defer func(begin time.Time) {
		s.requestCount.With("method", "register_incident").Add(1)
		s.requestLatency.With("method", "register_incident").Observe(time.Since(begin).Seconds())
		s.Store.SaveMetrics()
	}(time.Now())

	return s.Service.RegisterHandlingEvent(completed, id, voyageNumber, loc, eventType)
}
