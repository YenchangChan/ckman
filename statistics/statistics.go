package statistics

import (
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/collectors"
)

var (
	prefix = "ckman_"

	SyncDistSchemaTotal = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: prefix + "sync_dist_schema_total",
		},
		[]string{"cluster"},
	)
	SyncLogicSchema = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: prefix + "sync_logic_schema_total",
		},
		[]string{"cluster"},
	)
)

func init() {
	prometheus.MustRegister(SyncDistSchemaTotal)
	prometheus.MustRegister(SyncLogicSchema)
	prometheus.MustRegister(collectors.NewBuildInfoCollector())
}
