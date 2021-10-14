// Copyright (c) Microsoft Corporation.
// Licensed under the MIT License.

package main

import (
	"flag"
	"os"
	"os/signal"
	"syscall"

	eventhub "github.com/microsoft/azure-iothub-exporter/eventhub"
	metrics "github.com/microsoft/azure-iothub-exporter/metrics"
	server "github.com/microsoft/azure-iothub-exporter/server"
	log "github.com/sirupsen/logrus"
)

func main() {
	processArguments()

	metricsQueue := metrics.NewMetricsQueue()

	eventHubConnectionString := os.Getenv("EVENT_HUB_CONNECTION_STRING")
	eventHubConnectionString = "Endpoint=sb://iothub-ns-iotstarter-14938303-44c1d9af3e.servicebus.windows.net/;SharedAccessKeyName=iothubowner;SharedAccessKey=GAfYONV+2n8Bvv3dteY90PXD61wuW8osjcNIUOJ4BZo=;EntityPath=iotstarter-iothub"
	eventHubListener := eventhub.NewEventHubListener(eventHubConnectionString, metricsQueue)
	eventHubListener.Run()

	ioTMetricsServer := server.NewIoTMetricsServer(metricsQueue)
	ioTMetricsServer.StartHttpServer()

	signalChan := make(chan os.Signal, 1)
	signal.Notify(signalChan, os.Interrupt, syscall.SIGTERM)
	<-signalChan

	eventHubListener.Close()

}

func processArguments() {
	debugLevel := flag.Bool("debug", false, "Debug log level")
	flag.Parse()
	if *debugLevel {
		log.SetLevel(log.DebugLevel)
	}
}
