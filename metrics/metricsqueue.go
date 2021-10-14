package metrics

import (
	"encoding/json"
	"fmt"

	log "github.com/sirupsen/logrus"
)

type IoTMetric struct {
	TimeGeneratedUtc string
	Name             string
	Value            float64
	Labels           map[string]interface{}
}

type MetricsQueue struct {
	elements []string
}

const (
	MAX_LENGTH = 100

	testJson = `[{
		"TimeGeneratedUtc": "2021-10-13T20:08:54.525854Z",
		"Name": "edgeAgent_total_disk_read_bytes",
		"Value": 188447,
		"Labels": {
		  "iothub": "IoTStarter-iothub.azure-devices.net",
		  "edge_device": "edgeDeviceDEV",
		  "instance_number": "edd0a513-7047-4dea-90a3-d8d82fb0a9dd",
		  "module_name": "edgeHub",
		  "ms_telemetry": "False"
		}
	  },
	  {
		"TimeGeneratedUtc": "2021-10-13T20:08:54.525854Z",
		"Name": "edgeAgent_total_disk_read_bytes",
		"Value": 10563799,
		"Labels": {
		  "iothub": "IoTStarter-iothub.azure-devices.net",
		  "edge_device": "edgeDeviceDEV",
		  "instance_number": "edd0a513-7047-4dea-90a3-d8d82fb0a9dd",
		  "module_name": "prometheus",
		  "ms_telemetry": "False"
		}
	  }
	]`
)

func Unmarshal() {
	var iotmetrics []IoTMetric

	json.Unmarshal([]byte(testJson), &iotmetrics)

	fmt.Println(iotmetrics[1].Name)
	fmt.Println(iotmetrics[1].TimeGeneratedUtc)
	fmt.Println(iotmetrics[1].Value)

	for key, value := range iotmetrics[1].Labels {
		fmt.Println(key, value)
	}

}

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
	if len(p.elements) > MAX_LENGTH {
		p.shrink()
	}
}

func (p *MetricsQueue) Dequeue() string {
	result := ""
	if !p.IsEmpty() {
		result = p.elements[0]
		p.elements[0] = ""
		p.shrink()
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
		//TODO Figure out parsing problem
		if err != nil {
			log.Infof("Could not parse the metrics %v", err)
		}
		iotmetrics = append(iotmetrics, metrics...)
	}

	return iotmetrics

}
