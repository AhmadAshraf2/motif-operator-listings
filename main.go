package main

import (
	"fmt"
	"sync"

	"github.com/AhmadAshraf2/motif-operator-listings/apiServer"
	"github.com/AhmadAshraf2/motif-operator-listings/eventListener"
	_ "github.com/lib/pq"
	"github.com/spf13/viper"
)

func main() {
	InitConfigFile()
	var wg sync.WaitGroup
	// Function to restart a goroutine if it exits
	restartable := func(name string, fn func()) {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for {
				fmt.Printf("Starting %s...\n", name)
				fn()
				fmt.Printf("%s exited. Restarting...\n", name)
			}
		}()
	}

	// Start all processes with restartable logic
	restartable("MetaData listerner", eventListener.ListenUpdateMetaData)
	restartable("MotifListener listerner", eventListener.ListenMotifOperatorRegistered)
	restartable("apiServer", apiServer.StartApiServer)

	wg.Wait()
}

func InitConfigFile() {
	viper.AddConfigPath("./configs")
	viper.SetConfigName("config") // Register config file name (no extension)
	viper.SetConfigType("json")   // Look for specific type
	err := viper.ReadInConfig()
	if err != nil {
		fmt.Println("Error reading config file: ", err)
	}
}
