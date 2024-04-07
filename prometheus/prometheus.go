package prometheus

import (
	"github.com/prometheus/client_golang/prometheus"
)

var (
	RequestCounter = prometheus.NewCounter(prometheus.CounterOpts{
		Name: "myapp_requests_total",
		Help: "Total number of requests processed",
	})
	TaskDuration = prometheus.NewHistogramVec(prometheus.HistogramOpts{
		Name: "myapp_task_duration_seconds",
		Help: "Duration of task processing",
	}, []string{"taskID"})
	TaskSuccess = prometheus.NewGaugeVec(prometheus.GaugeOpts{
		Name: "myapp_task_success",
		Help: "Whether task was successful or not",
	}, []string{"taskID"})
)

func InitPrometheus() {
	prometheus.MustRegister(RequestCounter)
	prometheus.MustRegister(TaskDuration)
	prometheus.MustRegister(TaskSuccess)

}
