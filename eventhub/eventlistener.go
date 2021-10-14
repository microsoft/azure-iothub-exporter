package eventhub

import (
	"context"
	"fmt"

	eventhub "github.com/Azure/azure-event-hubs-go/v3"
	"github.com/microsoft/azure-iothub-exporter/metrics"
	log "github.com/sirupsen/logrus"
)

type EventHubListener struct {
	metricsqueue             *metrics.MetricsQueue
	eventHubSonnectionString string
}

func NewEventHubListener(eventHubSonnectionString string, metricsqueue *metrics.MetricsQueue) *EventHubListener {
	eventHubListener := new(EventHubListener)
	eventHubListener.eventHubSonnectionString = eventHubSonnectionString
	eventHubListener.metricsqueue = metricsqueue
	return eventHubListener
}

func (p *EventHubListener) newMessageHangdler(c context.Context, event *eventhub.Event) error {
	fmt.Println(string(event.Data))
	p.metricsqueue.Enqueue(string(event.Data))
	return nil
}

func (p *EventHubListener) Run() {
	hub, err := eventhub.NewHubFromConnectionString(p.eventHubSonnectionString)
	if err != nil {
		log.Fatal(err)
		return
	}

	ctx := context.Background()

	runtimeInfo, err := hub.GetRuntimeInformation(ctx)
	if err != nil {
		log.Fatal(err)
		return
	}

	for _, partitionID := range runtimeInfo.PartitionIDs {
		// Start receiving messages
		//
		// Receive blocks while attempting to connect to hub, then runs until listenerHandle.Close() is called
		// <- listenerHandle.Done() signals listener has stopped
		// listenerHandle.Err() provides the last error the receiver encountered
		fmt.Println(partitionID)
		_, err := hub.Receive(ctx, partitionID, p.newMessageHangdler, eventhub.ReceiveWithLatestOffset())
		if err != nil {
			log.Error(err)
			return
		}
		//time.Sleep(2 * time.Second)

		//listenerHandler.Close(ctx)
	}

}
