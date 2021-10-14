package main

import (
	"fmt"
	"net/http"

	metrics "github.com/microsoft/azure-iothub-exporter/metrics"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	log "github.com/sirupsen/logrus"
)

type IoTMetricsServer struct {
	registry     *prometheus.Registry
	metricsqueue *metrics.MetricsQueue
}

func NewIoTMetricsServer(metricsqueue *metrics.MetricsQueue) *IoTMetricsServer {
	ioTMetricsServer := new(IoTMetricsServer)
	ioTMetricsServer.metricsqueue = metricsqueue
	ioTMetricsServer.registry = prometheus.NewRegistry()
	return ioTMetricsServer
}

func (p *IoTMetricsServer) handleIotRequestPanic(w http.ResponseWriter, r *http.Request) {
	if err := recover(); err != nil {
		msg := fmt.Sprintf("%v", err)
		log.Errorf(msg)
		http.Error(w, msg, http.StatusBadRequest)
	}
}

func (p *IoTMetricsServer) queryMetrics(regsitry *prometheus.Registry) []*prometheus.GaugeVec {
	promMetrics := []*prometheus.GaugeVec{}
	ioTMetrics := p.metricsqueue.PopMetrics()

	if len(ioTMetrics) > 0 {
		var gaugeVec *prometheus.GaugeVec
		registeredMetrics := map[string]*prometheus.GaugeVec{}

		for _, ioTMetric := range ioTMetrics {

			labelNames := []string{}
			labelValues := prometheus.Labels{}

			if len(ioTMetric.Labels) > 0 {
				for key, value := range ioTMetric.Labels {
					labelNames = append(labelNames, key)
					labelValues[key] = value.(string)
				}

			}

			gaugeVec = registeredMetrics[ioTMetric.Name]

			if gaugeVec == nil {
				gaugeVec = prometheus.NewGaugeVec(prometheus.GaugeOpts{
					Name: ioTMetric.Name,
					Help: ioTMetric.Name,
				}, labelNames,
				)

				regsitry.Register(gaugeVec)
				registeredMetrics[ioTMetric.Name] = gaugeVec
			}

			gaugeVec.With(labelValues).Set(ioTMetric.Value)

			promMetrics = append(promMetrics, gaugeVec)

		}
	} else {
		log.Info("There is no new metrics")
	}

	return promMetrics
}

// fmt.Println(iotmetrics[1].Name)
// fmt.Println(iotmetrics[1].TimeGeneratedUtc)
// fmt.Println(iotmetrics[1].Value)

// for key, value := range iotmetrics[1].Labels {
// 	fmt.Println(key, value)
// }

func (p *IoTMetricsServer) handleIotRequest(w http.ResponseWriter, r *http.Request) {
	defer p.handleIotRequestPanic(w, r)

	// iotEventHubClient := NewIotEventHubClient(w, r)
	// prober.AddWorkspaces(opts.Loganalytics.Workspace...)
	// prober.Run()

	registry := prometheus.NewRegistry()

	p.queryMetrics(registry)

	// for _, metric := range metrics {
	// 	registry.Register(metric)
	// }

	h := promhttp.HandlerFor(registry, promhttp.HandlerOpts{})
	h.ServeHTTP(w, r)
}
