package request

import (
	"bytes"
	"net/http"
	"net/http/httputil"
	"strings"
	"time"

	"github.com/andig/evcc/util"
	"github.com/prometheus/client_golang/prometheus"
)

type roundTripper struct {
	log  *util.Logger
	base http.RoundTripper
}

const max = 2048 * 2

var (
	reqMetric *prometheus.SummaryVec
	errMetric *prometheus.CounterVec
)

func init() {
	reqMetric = prometheus.NewSummaryVec(prometheus.SummaryOpts{
		Namespace: "evcc",
		Subsystem: "http",
		Name:      "request_duration_seconds",
		Help:      "A summary of HTTP request durations",
		Objectives: map[float64]float64{
			0.5:  0.05,  // 50th percentile with a max. absolute error of 0.05
			0.9:  0.01,  // 90th percentile with a max. absolute error of 0.01
			0.99: 0.001, // 99th percentile with a max. absolute error of 0.001
		},
	}, []string{"host"})

	if err := prometheus.Register(reqMetric); err != nil {
		panic(err)
	}

	errMetric = prometheus.NewCounterVec(prometheus.CounterOpts{
		Namespace: "evcc",
		Subsystem: "http",
		Name:      "request_errors",
		Help:      "Total count of HTTP request errors",
	}, []string{"host"})

	if err := prometheus.Register(errMetric); err != nil {
		panic(err)
	}
}

// NewTripper creates a logging roundtrip handler
func NewTripper(log *util.Logger, base http.RoundTripper) http.RoundTripper {
	tripper := &roundTripper{
		log:  log,
		base: base,
	}

	return tripper
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}

func (r *roundTripper) RoundTrip(req *http.Request) (*http.Response, error) {
	r.log.TRACE.Printf("%s %s", req.Method, req.URL.String())

	var bld strings.Builder
	if body, err := httputil.DumpRequestOut(req, true); err == nil {
		bld.WriteString("\n")
		bld.Write(bytes.TrimSpace(body[:min(max, len(body))]))
	}

	startTime := time.Now()
	resp, err := r.base.RoundTrip(req)

	if err == nil {
		if mv, err := reqMetric.GetMetricWith(prometheus.Labels{
			"host": req.URL.Hostname(),
		}); err == nil {
			mv.Observe(time.Since(startTime).Seconds())
		}

		if body, err := httputil.DumpResponse(resp, true); err == nil {
			bld.WriteString("\n\n")
			bld.Write(bytes.TrimSpace(body[:min(max, len(body))]))
		}
	} else {
		if mv, err := errMetric.GetMetricWith(prometheus.Labels{
			"host": req.URL.Hostname(),
		}); err == nil {
			mv.Add(1)
		}
	}

	if bld.Len() > 0 {
		r.log.TRACE.Println(bld.String())
	}

	return resp, err
}
