// Copyright (c) Microsoft Corporation.
// Licensed under the MIT License.

package eventhub

import (
	"context"

	eventhub "github.com/Azure/azure-event-hubs-go/v3"
	"github.com/microsoft/azure-iothub-exporter/metrics"
	log "github.com/sirupsen/logrus"
)

type EventHubListener struct {
	metricsqueue             *metrics.MetricsQueue
	eventHubSonnectionString string
	hub                      *eventhub.Hub
}

func NewEventHubListener(eventHubSonnectionString string, metricsqueue *metrics.MetricsQueue) *EventHubListener {
	eventHubListener := new(EventHubListener)
	eventHubListener.eventHubSonnectionString = eventHubSonnectionString
	eventHubListener.metricsqueue = metricsqueue
	return eventHubListener
}

func (p *EventHubListener) newMessageHangdler(c context.Context, event *eventhub.Event) error {
	log.Info("New message has arrived from Event Hub")
	log.Debugln("===========message start=========")
	log.Debug(string(event.Data))
	log.Debugln("===========message end=========")
	p.metricsqueue.Enqueue(string(event.Data))
	return nil
}

//https://github.com/Azure/azure-event-hubs-go
func (p *EventHubListener) Run() {
	hub, err := eventhub.NewHubFromConnectionString(p.eventHubSonnectionString)
	if err != nil {
		log.Fatal(err)
		return
	}

	p.hub = hub
	ctx := context.Background()

	runtimeInfo, err := hub.GetRuntimeInformation(ctx)
	if err != nil {
		log.Fatal(err)
		return
	}

	for _, partitionID := range runtimeInfo.PartitionIDs {
		_, err := hub.Receive(ctx, partitionID, p.newMessageHangdler, eventhub.ReceiveWithLatestOffset())
		if err != nil {
			log.Fatal(err)
			return
		}
	}
	log.Info("Event Hub listener has started")
}

func (p *EventHubListener) Close() {
	err := p.hub.Close(context.Background())
	if err != nil {
		log.Fatal(err)
	}
	log.Info("Hub is closed")
}
