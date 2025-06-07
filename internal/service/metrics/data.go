package metrics

import (
	"net/http"
	"sync"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

type Logger interface {
	Errorf(string, ...interface{})
	Infof(string, ...interface{})
}

type MetricsReg struct {
	log Logger

	ClientUpStat   *prometheus.GaugeVec
	ClientDownStat *prometheus.GaugeVec

	InboundUpStat   *prometheus.GaugeVec
	InboundDownStat *prometheus.GaugeVec

	// =============================
	Registry *prometheus.Registry
	mu       sync.RWMutex
}

func NewMetricsReg(log Logger) *MetricsReg {

	self := MetricsReg{
		Registry: prometheus.NewRegistry(),

		ClientUpStat: newGaugeVec(
			"client_traffic_up",
			"3X-UI user down stats",
			[]string{"email"},
		),
		ClientDownStat: newGaugeVec(
			"client_traffic_down",
			"3X-UI user up stats",
			[]string{"email"},
		),

		InboundUpStat: newGaugeVec(
			"inbound_traffic_up",
			"3X-UI inbound up stats",
			[]string{"tag"},
		),
		InboundDownStat: newGaugeVec(
			"inbound_traffic_down",
			"3X-UI inbound up stats",
			[]string{"tag"},
		),
	}

	self.log = log

	c := []prometheus.Collector{
		self.ClientUpStat,
		self.ClientDownStat,
		self.InboundUpStat,
		self.InboundDownStat,
	}

	for _, col := range c {
		if err := self.Registry.Register(col); err != nil {
			log.Errorf("register error: %v", err)
		}
	}
	return &self
}

type ClientExporter interface {
	DownTraffic() float64
	UpTraffic() float64
	EmailString() string
}

type InboundExporter interface {
	DownTraffic() float64
	UpTraffic() float64
	TagString() string
}

func (mre *MetricsReg) UpdateClient(client ClientExporter) {
	mre.mu.Lock()
	defer mre.mu.Unlock()

	setMetric(mre.ClientUpStat, client.UpTraffic(), client.EmailString())
	setMetric(mre.ClientDownStat, client.DownTraffic(), client.EmailString())
}

func (mre *MetricsReg) UpdateInbound(inb InboundExporter) {
	mre.mu.Lock()
	defer mre.mu.Unlock()

	setMetric(mre.InboundUpStat, inb.UpTraffic(), inb.TagString())
	setMetric(mre.InboundDownStat, inb.DownTraffic(), inb.TagString())
}

func (mre *MetricsReg) Metric() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		mre.mu.RLock()
		reg := mre.Registry
		mre.mu.RUnlock()

		if _, err := reg.Gather(); err != nil {
			mre.log.Errorf("metrics gather error: %v", err)
		}

		promhttp.HandlerFor(reg, promhttp.HandlerOpts{}).ServeHTTP(w, r)
	})
}
