// Copyright (c) Microsoft Corporation.
// Licensed under the MIT License.

package server

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

const (
	PrometheusPort = "8080"
)

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
				log.Debugf("Metrics %v has been registered ro response", ioTMetric.Name)
			}

			gaugeVec.With(labelValues).Set(ioTMetric.Value)

			promMetrics = append(promMetrics, gaugeVec)

		}
	} else {
		log.Info("There is no new metrics")
	}

	log.Infof("%v merics to respond", len(promMetrics))
	return promMetrics
}

func (p *IoTMetricsServer) handleIotRequest(w http.ResponseWriter, r *http.Request) {
	defer p.handleIotRequestPanic(w, r)

	registry := prometheus.NewRegistry()
	p.queryMetrics(registry)

	h := promhttp.HandlerFor(registry, promhttp.HandlerOpts{})
	h.ServeHTTP(w, r)
}

func (p *IoTMetricsServer) StartHttpServer() {
	http.Handle("/metrics", promhttp.Handler())

	http.HandleFunc("/iotmetrics", p.handleIotRequest)

	log.Infof("Starting HTTP server to listen on: %v /n %v", "/metrics", "/iotmetrics")

	log.Fatal(http.ListenAndServe(":"+PrometheusPort, nil))
}
