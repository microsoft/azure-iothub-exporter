package main

import (
	//     aad "github.com/Azure/azure-amqp-common-go/aad"
	"context"
	"fmt"
	"os"
	"os/signal"
	"time"

	eventhub "github.com/Azure/azure-event-hubs-go/v3"

	//     azure "github.com/Azure/go-autorest/autorest/azure"

	log "github.com/sirupsen/logrus"
)

func main() {
	//https://github.com/Azure/azure-event-hubs-go

	log.Infof("init Azure")

	// Azure Event Hub connection string
	eventHubConnStr := "Endpoint=sb://iothub-ns-iotstarter-14938303-44c1d9af3e.servicebus.windows.net/;SharedAccessKeyName=iothubowner;SharedAccessKey=GAfYONV+2n8Bvv3dteY90PXD61wuW8osjcNIUOJ4BZo=;EntityPath=iotstarter-iothub"

	hub, err := eventhub.NewHubFromConnectionString(eventHubConnStr)

	if err != nil {
		fmt.Println(err)
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), 20*time.Second)
	defer cancel()

	// handler := func(c context.Context, event *eventhub.Event) error {
	// 	fmt.Println(string(event.Data))
	// 	return nil
	// }

	// listen to each partition of the Event Hub
	runtimeInfo, err := hub.GetRuntimeInformation(ctx)
	if err != nil {
		fmt.Println(err)
		return
	}

	for _, partitionID := range runtimeInfo.PartitionIDs {
		// Start receiving messages
		//
		// Receive blocks while attempting to connect to hub, then runs until listenerHandle.Close() is called
		// <- listenerHandle.Done() signals listener has stopped
		// listenerHandle.Err() provides the last error the receiver encountered
		fmt.Println(partitionID)
		// _, err = hub.Receive(ctx, partitionID, handler)
		// if err != nil {
		// 	fmt.Println(err)
		// 	return
		// }
	}

	// Wait for a signal to quit:
	signalChan := make(chan os.Signal, 1)
	signal.Notify(signalChan, os.Interrupt, os.Kill)
	<-signalChan

	err = hub.Close(context.Background())
	if err != nil {
		fmt.Println(err)
	}

	//

	// parsed, err := conn.ParsedConnectionFromStr(eventHubConnStr)
	// if err != nil {
	// 	log.Panicf("failed to parse Event Hub connection string: %s\n", err)
	// }

	// tokenProvider, err := sas.NewTokenProvider(sas.TokenProviderWithKey(parsed.KeyName, parsed.Key))
	// if err != nil {
	// 	log.Panicf("failed to configure AAD JWT provider: %s\n", err)
	// }

	// azureEnv, err := azure.EnvironmentFromName("AzurePublicCloud")
	// if err != nil {
	// 	log.Fatalf("could not get azure.Environment struct: %s\n", err)
	// }

	// config := loadConfigs()

	// // cred, err := storageLeaser.NewAADSASCredential(
	// // 	config.subscriptionID,
	// // 	config.resourceGroupName,
	// // 	config.storageAccountName,
	// // 	config.storageContainerName,
	// // 	storageLeaser.AADSASCredentialWithEnvironmentVars())
	// // if err != nil {
	// // 	log.Fatalf("could not prepare a storage credential: %s\n", err)
	// // }

	// // create a new Azure Storage Leaser / Checkpointer
	// cred, err := azblob.NewSharedKeyCredential(config.storageAccountName, config.storageContainerName)
	// if err != nil {
	// 	fmt.Println(err)
	// 	return
	// }

	// leaserCheckpointer, err := storageLeaser.NewStorageLeaserCheckpointer(
	// 	cred,
	// 	config.storageAccountName,
	// 	config.storageContainerName,
	// 	azureEnv)
	// if err != nil {
	// 	log.Fatalf("could not prepare a storage leaserCheckpointer: %s\n", err)
	// }

	// ctx, cancel := context.WithTimeout(context.Background(), 20*time.Second)
	// defer cancel()
	// // create a new EPH processor
	// processor, err := eph.New(ctx, parsed.Namespace, parsed.HubName, tokenProvider, leaserCheckpointer, leaserCheckpointer)
	// if err != nil {
	// 	fmt.Println(err)
	// 	return
	// }

	// // register a message handler -- many can be registered
	// handlerID, err := processor.RegisterHandler(ctx,
	// 	func(c context.Context, e *eventhub.Event) error {
	// 		fmt.Println(string(e.Data))
	// 		return nil
	// 	})
	// if err != nil {
	// 	fmt.Println(err)
	// 	return
	// }

	// fmt.Printf("handler id: %q is running\n", handlerID)

	// // unregister a handler to stop that handler from receiving events
	// // processor.UnregisterHandler(ctx, handleID)

	// // start handling messages from all of the partitions balancing across multiple consumers
	// err = processor.StartNonBlocking(ctx)
	// if err != nil {
	// 	fmt.Println(err)
	// 	return
	// }

	// // Wait for a signal to quit:
	// signalChan := make(chan os.Signal, 1)
	// signal.Notify(signalChan, os.Interrupt, os.Kill)
	// <-signalChan

	// err = processor.Close(context.Background())
	// if err != nil {
	// 	fmt.Println(err)
	// 	return
	// }

	// // ctx := context.Background()
	// // p, err := eph.New(
	// // 	ctx,
	// // 	"nsName",
	// // 	"hubName",
	// // 	tokenProvider,
	// // 	leaserCheckpointer,
	// // 	leaserCheckpointer)
	// // if err != nil {
	// // 	log.Fatalf("failed to create EPH: %s\n", err)
	// // }
	// // defer p.Close(context.Background())

	// // handler := func(ctx context.Context, event *eventhubs.Event) error {
	// // 	fmt.Printf("received: %s\n", string(event.Data))
	// // 	return nil
	// // }

	// // // register the handler with the EPH
	// // _, err = p.RegisterHandler(ctx, handler)
	// // if err != nil {
	// // 	log.Fatalf("failed to register handler: %s\n", err)
	// // }

}
