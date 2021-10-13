package main

import (
	"os"
)

type Config struct {
	subscriptionID       string
	resourceGroupName    string
	storageAccountName   string
	storageContainerName string
}

func loadConfigs() *Config {
	config := new(Config)
	config.subscriptionID = os.Getenv("AZURE_SUBSCRIPTION_ID")
	config.resourceGroupName = os.Getenv("EVENT_HUB_CHECPOINT_RESOURCE_GROUP")
	config.storageAccountName = os.Getenv("EVENT_HUB_CHECPOINT_STORAGE_ACCOUNT_NAME")
	config.storageContainerName = os.Getenv("EVENT_HUB_CHECPOINT_STORAGE_CONTAINER_NAME")
	return config
}
