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
	Debugf(string, ...interface{})
}

type MetricsReg struct {
	log Logger

	ClientUpStat    *prometheus.GaugeVec
	ClientDownStat  *prometheus.GaugeVec
	ClientTotalStat *prometheus.GaugeVec

	InboundUpStat   *prometheus.GaugeVec
	InboundDownStat *prometheus.GaugeVec

	InboundAllUpStat   *prometheus.GaugeVec
	InboundAllDownStat *prometheus.GaugeVec

	// =============================
	Registry *prometheus.Registry

	muClient     sync.RWMutex
	muInbound    sync.RWMutex
	muInboundAll sync.RWMutex
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
		ClientTotalStat: newGaugeVec(
			"client_traffic_total",
			"3X-UI user total stats",
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

		InboundAllUpStat: newGaugeVec(
			"inbound_all_up",
			"3X-UI inbound up stats",
			[]string{"name", "proto", "email"},
		),

		InboundAllDownStat: newGaugeVec(
			"inbound_all_down",
			"3X-UI inbound up stats",
			[]string{"name", "proto", "email"},
		),
	}

	self.log = log

	c := []prometheus.Collector{
		self.ClientUpStat,
		self.ClientDownStat,
		self.ClientTotalStat,
		self.InboundUpStat,
		self.InboundDownStat,
		self.InboundAllUpStat,
		self.InboundAllDownStat,
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

type TotalExporter interface {
	TotalTraffic() float64
}

type InboundExporter interface {
	DownTraffic() float64
	UpTraffic() float64
	TagString() string
}

func (mre *MetricsReg) UpdateClient(client ClientExporter, tot TotalExporter) {
	mre.muClient.Lock()
	setMetric(mre.ClientUpStat, client.UpTraffic(), client.EmailString())
	setMetric(mre.ClientDownStat, client.DownTraffic(), client.EmailString())
	setMetric(mre.ClientTotalStat, tot.TotalTraffic(), client.EmailString())
	mre.muClient.Unlock()
}

func (mre *MetricsReg) UpdateInbound(inb InboundExporter) {
	mre.muInbound.Lock()
	setMetric(mre.InboundUpStat, inb.UpTraffic(), inb.TagString())
	setMetric(mre.InboundDownStat, inb.DownTraffic(), inb.TagString())
	mre.muInbound.Unlock()
}

type ProtoExporter interface {
	ProtocolString() string
	NameString() string
}

func (mre *MetricsReg) UpdateStats(inb ClientExporter, pr ProtoExporter) {
	mre.muInboundAll.Lock()
	setMetric(mre.InboundAllUpStat, inb.UpTraffic(), pr.NameString(), pr.ProtocolString(), inb.EmailString())
	setMetric(mre.InboundAllDownStat, inb.DownTraffic(), pr.NameString(), pr.ProtocolString(), inb.EmailString())
	mre.muInboundAll.Unlock()
}

func (mre *MetricsReg) Metric() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		mre.muClient.RLock()
		defer mre.muClient.RUnlock()

		mre.muInbound.RLock()
		defer mre.muInbound.RUnlock()

		mre.muInboundAll.RLock()
		defer mre.muInboundAll.RUnlock()

		mre.log.Debugf("export request from - %s", getReqAddr(r))

		reg := mre.Registry
		if _, err := reg.Gather(); err != nil {
			mre.log.Errorf("metrics gather error: %v", err)
		}

		promhttp.HandlerFor(reg, promhttp.HandlerOpts{}).ServeHTTP(w, r)
	})
}

func getReqAddr(r *http.Request) string {
	ip := r.Header.Get("X-Real-IP")
	if ip != "" {
		return ip
	}
	return r.RemoteAddr
}
