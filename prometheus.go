package main

import "github.com/prometheus/client_golang/prometheus"

var (
	tagNames = []string{
		"server",
		"scheme",
		"method",
		"hostname",
		"status",
		"uri",
	}

	promReqCount = prometheus.NewCounterVec(prometheus.CounterOpts{
		Name: "ngx_request_count",
		Help: "request count",
	}, tagNames)

	promReqCountByCountry = prometheus.NewCounterVec(prometheus.CounterOpts{
		Name: "ngx_request_count_by_country",
		Help: "request count by country",
	}, []string{"clientCountry"})

	promReqCountByProtocol = prometheus.NewCounterVec(prometheus.CounterOpts{
		Name: "ngx_request_count_by_protocol",
		Help: "request count by protocol",
	}, []string{"protocol"})

	promReqSize = prometheus.NewCounterVec(prometheus.CounterOpts{
		Name: "ngx_request_size_bytes",
		Help: "request size in bytes",
	}, tagNames)

	promRespSize = prometheus.NewCounterVec(prometheus.CounterOpts{
		Name: "ngx_response_size_bytes",
		Help: "response size in bytes",
	}, tagNames)

	promTime = prometheus.NewHistogramVec(prometheus.HistogramOpts{
		Name: "ngx_request_duration_seconds",
		Help: "request serving time in seconds",
	}, tagNames)
)

func init() {
	prometheus.MustRegister(
		promReqCount,
		promReqCountByCountry,
		promReqCountByProtocol,
		promReqSize,
		promRespSize,
		promTime,
	)
}

func prometheusObserve(l *logEntry) {
	tags := []string{
		l.server,
		l.scheme,
		l.method,
		l.hostname,
		l.status,
		l.uri,
	}

	promReqCount.WithLabelValues(tags...).Inc()
	promReqCountByCountry.WithLabelValues(l.clientCountry).Inc()
	promReqCountByProtocol.WithLabelValues(l.protocol).Inc()

	promRespSize.WithLabelValues(tags...).Add(float64(l.bytesSent))
	promReqSize.WithLabelValues(tags...).Add(float64(l.bytesRcvd))
	promTime.WithLabelValues(tags...).Observe(l.duration)
}
