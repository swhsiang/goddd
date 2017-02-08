package booking

import (
	"time"

	"github.com/go-kit/kit/metrics"

	"github.com/marcusolsson/goddd/cargo"
	"github.com/marcusolsson/goddd/location"

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

func (s *instrumentingService) BookNewCargo(origin, destination location.UNLocode, deadline time.Time) (cargo.TrackingID, error) {
	defer func(begin time.Time) {
		s.requestCount.With("method", "book").Add(1)
		s.requestLatency.With("method", "book").Observe(time.Since(begin).Seconds())
		s.Store.SaveMetrics()
	}(time.Now())

	return s.Service.BookNewCargo(origin, destination, deadline)
}

func (s *instrumentingService) LoadCargo(id cargo.TrackingID) (c Cargo, err error) {
	defer func(begin time.Time) {
		s.requestCount.With("method", "load").Add(1)
		s.requestLatency.With("method", "load").Observe(time.Since(begin).Seconds())
		s.Store.SaveMetrics()
	}(time.Now())

	return s.Service.LoadCargo(id)
}

func (s *instrumentingService) RequestPossibleRoutesForCargo(id cargo.TrackingID) []cargo.Itinerary {
	defer func(begin time.Time) {
		s.requestCount.With("method", "request_routes").Add(1)
		s.requestLatency.With("method", "request_routes").Observe(time.Since(begin).Seconds())
		s.Store.SaveMetrics()
	}(time.Now())

	return s.Service.RequestPossibleRoutesForCargo(id)
}

func (s *instrumentingService) AssignCargoToRoute(id cargo.TrackingID, itinerary cargo.Itinerary) (err error) {
	defer func(begin time.Time) {
		s.requestCount.With("method", "assign_to_route").Add(1)
		s.requestLatency.With("method", "assign_to_route").Observe(time.Since(begin).Seconds())
		s.Store.SaveMetrics()
	}(time.Now())

	return s.Service.AssignCargoToRoute(id, itinerary)
}

func (s *instrumentingService) ChangeDestination(id cargo.TrackingID, l location.UNLocode) (err error) {
	defer func(begin time.Time) {
		s.requestCount.With("method", "change_destination").Add(1)
		s.requestLatency.With("method", "change_destination").Observe(time.Since(begin).Seconds())
		s.Store.SaveMetrics()
	}(time.Now())

	return s.Service.ChangeDestination(id, l)
}

func (s *instrumentingService) Cargos() []Cargo {
	defer func(begin time.Time) {
		s.requestCount.With("method", "list_cargos").Add(1)
		s.requestLatency.With("method", "list_cargos").Observe(time.Since(begin).Seconds())
		s.Store.SaveMetrics()
	}(time.Now())

	return s.Service.Cargos()
}

func (s *instrumentingService) Locations() []Location {
	defer func(begin time.Time) {
		s.requestCount.With("method", "list_locations").Add(1)
		s.requestLatency.With("method", "list_locations").Observe(time.Since(begin).Seconds())
		s.Store.SaveMetrics()
	}(time.Now())

	return s.Service.Locations()
}
