package eventListener

import (
	"context"
	"log"
	"strings"

	"github.com/AhmadAshraf2/motif-operator-listings/eventHandler"
	"github.com/AhmadAshraf2/motif-operator-listings/utils"
	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/common"
	"github.com/spf13/viper"
)

func ListenUpdateMetaData() {
	const abiFilePath = "abis/eigen_layer_delegation_manager.abi"
	delegationManagerAddress := viper.GetString("eigen_delegation_manager_contract")
	delegationManagerABI, err := utils.ReadABIFromFile(abiFilePath)
	if err != nil {
		log.Fatalf("Failed to read ABI from file: %v", err)
	}
	client, err := utils.GetEthClient()
	if err != nil {
		log.Fatalf("Failed to connect to Ethereum node: %v", err)
	}

	// Parse the contract ABI
	parsedABI, err := abi.JSON(strings.NewReader(delegationManagerABI))
	if err != nil {
		log.Fatalf("Failed to parse contract ABI: %v", err)
	}

	// Define the contract address
	contractAddress := common.HexToAddress(delegationManagerAddress)

	// Define the event topics
	operatorRegisteredSig := parsedABI.Events["OperatorRegistered"].ID.Hex()
	operatorMetadataURIUpdatedSig := parsedABI.Events["OperatorMetadataURIUpdated"].ID.Hex()

	// Create a query to filter logs for both events
	query := ethereum.FilterQuery{
		Addresses: []common.Address{contractAddress},
		Topics:    [][]common.Hash{{common.HexToHash(operatorRegisteredSig), common.HexToHash(operatorMetadataURIUpdatedSig)}},
	}

	// Fetch logs
	logs, err := client.FilterLogs(context.Background(), query)
	if err != nil {
		log.Fatalf("Failed to fetch logs: %v", err)
	}

	// Process logs
	for _, vLog := range logs {
		switch vLog.Topics[0].Hex() {
		case operatorMetadataURIUpdatedSig:
			eventHandler.HandleOperatorMetadataURIUpdated(parsedABI, vLog)
		}
	}

}

func ListenMotifOperatorRegistered() {
	const abiFilePath = "abis/motif_stake_registry.abi" // Path to the ABI file
	motifStakeRegistryAddress := viper.GetString("motif_stake_contract")

	// Read the ABI from the file
	delegationManagerABI, err := utils.ReadABIFromFile(abiFilePath)
	if err != nil {
		log.Fatalf("Failed to read ABI from file: %v", err)
	}

	// Connect to the Ethereum client
	client, err := utils.GetEthClient()
	if err != nil {
		log.Fatalf("Failed to connect to Ethereum node: %v", err)
	}

	// Parse the contract ABI
	parsedABI, err := abi.JSON(strings.NewReader(delegationManagerABI))
	if err != nil {
		log.Fatalf("Failed to parse contract ABI: %v", err)
	}

	// Define the contract address
	contractAddress := common.HexToAddress(motifStakeRegistryAddress)

	// Define the event topic
	operatorBtcKeyRegisteredSig := parsedABI.Events["OperatorBtcKeyRegistered"].ID.Hex()

	// Create a query to filter logs for the event
	query := ethereum.FilterQuery{
		Addresses: []common.Address{contractAddress},
		Topics:    [][]common.Hash{{common.HexToHash(operatorBtcKeyRegisteredSig)}},
	}

	// Fetch logs
	logs, err := client.FilterLogs(context.Background(), query)
	if err != nil {
		log.Fatalf("Failed to fetch logs: %v", err)
	}

	// Process logs
	for _, vLog := range logs {
		eventHandler.HandleMotifOperatorRegistered(parsedABI, vLog)
	}
}

func ListenOperatorBtckeyDeregistered() {
	const abiFilePath = "abis/motif_stake_registry.abi" // Path to the ABI file
	motifStakeRegistryAddress := viper.GetString("motif_stake_contract")

	// Read the ABI from the file
	delegationManagerABI, err := utils.ReadABIFromFile(abiFilePath)
	if err != nil {
		log.Fatalf("Failed to read ABI from file: %v", err)
	}

	// Connect to the Ethereum client
	client, err := utils.GetEthClient()
	if err != nil {
		log.Fatalf("Failed to connect to Ethereum node: %v", err)
	}

	// Parse the contract ABI
	parsedABI, err := abi.JSON(strings.NewReader(delegationManagerABI))
	if err != nil {
		log.Fatalf("Failed to parse contract ABI: %v", err)
	}

	// Define the contract address
	contractAddress := common.HexToAddress(motifStakeRegistryAddress)

	// Define the event topic
	operatorBtckeyDeregisteredSig := parsedABI.Events["OperatorBtckeyDeregistered"].ID.Hex()

	// Create a query to filter logs for the event
	query := ethereum.FilterQuery{
		Addresses: []common.Address{contractAddress},
		Topics:    [][]common.Hash{{common.HexToHash(operatorBtckeyDeregisteredSig)}},
	}

	// Fetch logs
	logs, err := client.FilterLogs(context.Background(), query)
	if err != nil {
		log.Fatalf("Failed to fetch logs: %v", err)
	}

	// Process logs
	for _, vLog := range logs {
		eventHandler.HandleOperatorBtckeyDeregistered(parsedABI, vLog)
	}
}
