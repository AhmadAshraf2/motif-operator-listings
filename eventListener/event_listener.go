package eventListener

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/AhmadAshraf2/motif-operator-listings/DelegationManager"
	"github.com/AhmadAshraf2/motif-operator-listings/MotifRegistry"
	"github.com/AhmadAshraf2/motif-operator-listings/eventHandler"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/ethereum/go-ethereum/rpc"
	"github.com/spf13/viper"
)

func SubscribeToBtcKeyRegistered() {
	// Create a new instance of the contract binding

	// oprEthAccount := LoadEthAccount()
	client, err := rpc.Dial(viper.GetString("eth_ws_host"))
	if err != nil {
		log.Fatalf("Failed to connect to the Ethereum client: %v", err)
	}

	ethClient := ethclient.NewClient(client)
	defer ethClient.Close()
	defer client.Close()

	motifRegistryAddr := common.HexToAddress(viper.GetString("motif_stake_contract"))
	motifRegistry, err := MotifRegistry.NewMotifRegistry(motifRegistryAddr, ethClient)
	if err != nil {
		fmt.Println("Failed to instantiate contract:", err)
		panic(err)
	}

	// // Create a channel for the events
	ch := make(chan *MotifRegistry.MotifRegistryOperatorBtcKeyRegistered)

	// // Create a subscription
	sub, err := motifRegistry.WatchOperatorBtcKeyRegistered(
		&bind.WatchOpts{Context: context.Background()},
		ch,
		[]common.Address{},
	)
	if err != nil {
		fmt.Println("Failed to subscribe to operator btc key register events:", err)
		panic(err)
	}

	fmt.Println("Successfully subscribed to operatorbtc key registered events")

	// Handle events in a loop
	for {
		select {
		case err := <-sub.Err():
			fmt.Println("Subscription error deposit:", err)
			if err == nil {
				return
			}
			time.Sleep(1 * time.Minute)

		case event := <-ch:
			eventHandler.HandleMotifOperatorRegistered(event)
		}
	}
}

func SubscribeToBtcKeyDeregistered() {
	// Create a new instance of the contract binding

	// oprEthAccount := LoadEthAccount()
	client, err := rpc.Dial(viper.GetString("eth_ws_host"))
	if err != nil {
		log.Fatalf("Failed to connect to the Ethereum client: %v", err)
	}

	ethClient := ethclient.NewClient(client)
	defer ethClient.Close()
	defer client.Close()

	motifRegistryAddr := common.HexToAddress(viper.GetString("motif_stake_contract"))
	motifRegistry, err := MotifRegistry.NewMotifRegistry(motifRegistryAddr, ethClient)
	if err != nil {
		fmt.Println("Failed to instantiate contract:", err)
		panic(err)
	}

	// // Create a channel for the events
	ch := make(chan *MotifRegistry.MotifRegistryOperatorBtckeyDeregistered)

	// // Create a subscription
	sub, err := motifRegistry.WatchOperatorBtckeyDeregistered(
		&bind.WatchOpts{Context: context.Background()},
		ch,
		[]common.Address{},
	)
	if err != nil {
		fmt.Println("Failed to subscribe to deregister btc events:", err)
		panic(err)
	}

	fmt.Println("Successfully subscribed to operator btc key deregistered events")

	// Handle events in a loop
	for {
		select {
		case err := <-sub.Err():
			fmt.Println("Subscription error deposit:", err)
			if err == nil {
				return
			}
			time.Sleep(1 * time.Minute)

		case event := <-ch:
			eventHandler.HandleMotifOperatorDeregistered(event)
		}
	}
}

func SubscribeToOperatorRegisteredEigenlayer() {
	// Create a new instance of the contract binding

	// oprEthAccount := LoadEthAccount()
	client, err := rpc.Dial(viper.GetString("eth_ws_host"))
	if err != nil {
		log.Fatalf("Failed to connect to the Ethereum client: %v", err)
	}

	ethClient := ethclient.NewClient(client)
	defer ethClient.Close()
	defer client.Close()

	delegationManagerAddr := common.HexToAddress(viper.GetString("eigen_delegation_manager_contract"))
	delegationManager, err := DelegationManager.NewDelegationManager(delegationManagerAddr, ethClient)
	if err != nil {
		fmt.Println("Failed to instantiate contract:", err)
		panic(err)
	}

	// // Create a channel for the events
	ch := make(chan *DelegationManager.DelegationManagerOperatorMetadataURIUpdated)

	// // Create a subscription
	sub, err := delegationManager.WatchOperatorMetadataURIUpdated(
		&bind.WatchOpts{Context: context.Background()},
		ch,
		[]common.Address{},
	)
	if err != nil {
		fmt.Println("Failed to subscribe to metadata update events:", err)
		panic(err)
	}

	fmt.Println("Successfully subscribed to metadata update events")

	// Handle events in a loop
	for {
		select {
		case err := <-sub.Err():
			fmt.Println("Subscription error deposit:", err)
			if err == nil {
				return
			}
			time.Sleep(1 * time.Minute)

		case event := <-ch:
			eventHandler.HandleOperatorMetadataURIUpdatedEigenLayer(event)
		}
	}
}
