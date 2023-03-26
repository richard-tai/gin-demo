package metrics

import (
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
)

var (
	ApiTotal = promauto.NewCounterVec(prometheus.CounterOpts{
		Name: "z_api_total",
		Help: "Total number of api accessed",
	}, []string{"path"})

	UserNum = promauto.NewGaugeVec(prometheus.GaugeOpts{
		Name: "z_user_num",
		Help: "Number of user",
	}, []string{"client"})
)
