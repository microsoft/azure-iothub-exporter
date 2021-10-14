// Copyright (c) Microsoft Corporation.
// Licensed under the MIT License.

package metrics

import (
	"encoding/json"

	log "github.com/sirupsen/logrus"
)

type IoTMetric struct {
	TimeGeneratedUtc string
	Name             string
	Value            float64
	Labels           map[string]interface{}
}

//TODO: Perhaps we can use Prometheus collectors instead of the Queue
//      1. Receive metric -> Log with Prometheus to shared registry
//      2. Scrapping request -> promhttp.HandlerFor(shared registry, ...

type MetricsQueue struct {
	elements []string
}

const (
	MAX_LENGTH = 1000
)

func NewMetricsQueue() *MetricsQueue {
	metricsQueue := new(MetricsQueue)
	metricsQueue.elements = []string{}

	return metricsQueue
}

func (p *MetricsQueue) shrink() {
	p.elements = p.elements[1:]
}

func (p *MetricsQueue) Enqueue(metric string) {
	p.elements = append(p.elements, metric)
	log.Debugf("Message has been added to the queue. Queue size is %v", len(p.elements))
	if len(p.elements) > MAX_LENGTH {
		p.shrink()
		log.Debugf("Message queue has been shrunk. Queue size is %v", len(p.elements))
	}

}

func (p *MetricsQueue) Dequeue() string {
	result := ""
	if !p.IsEmpty() {
		result = p.elements[0]
		p.elements[0] = ""
		p.shrink()
		log.Debugf("Message has been scrapped from the queue. Queue size is %v", len(p.elements))
	}
	return result
}

func (p *MetricsQueue) IsEmpty() bool {
	return len(p.elements) == 0
}

func (p *MetricsQueue) PopMetrics() []IoTMetric {
	iotmetrics := []IoTMetric{}

	for !p.IsEmpty() {
		metric := p.Dequeue()
		metrics := []IoTMetric{}
		err := json.Unmarshal([]byte(metric), &metrics)
		//TODO Handle NaN conversion issue
		if err != nil {
			log.Warnf("Could not parse the metrics %v", err)
			log.Debug("===parse error metrics====")
			log.Debug(metric)
			log.Debug("==========================")
		}
		iotmetrics = append(iotmetrics, metrics...)
	}

	log.Infof("Scrapped %v metrics", len(iotmetrics))

	return iotmetrics

}
