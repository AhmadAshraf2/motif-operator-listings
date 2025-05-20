package eventHandler

import (
	"database/sql"
	"fmt"
	"log"

	"github.com/AhmadAshraf2/motif-operator-listings/db"
	"github.com/AhmadAshraf2/motif-operator-listings/types"
	"github.com/AhmadAshraf2/motif-operator-listings/utils"
	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/common"
	etypes "github.com/ethereum/go-ethereum/core/types"
)

var database *sql.DB

func HandleOperatorMetadataURIUpdated(parsedABI abi.ABI, vLog etypes.Log) {
	// Decode the event
	dbconn := db.InitDB()
	defer dbconn.Close()
	event := struct {
		Operator    common.Address
		MetadataURI string
	}{}

	err := parsedABI.UnpackIntoInterface(&event, "OperatorMetadataURIUpdated", vLog.Data)
	if err != nil {
		log.Printf("Failed to unpack OperatorMetadataURIUpdated event: %v", err)
		return
	}

	// Log the event details
	fmt.Printf("OperatorMetadataURIUpdated:\n")
	fmt.Printf("  Operator: %s\n", event.Operator.Hex())
	fmt.Printf("  Metadata URI: %s\n", event.MetadataURI)
	opr := types.Operator{}
	opr.EthAddress = event.Operator.Hex()
	metadata, err := utils.FetchMetadataContent(event.MetadataURI)
	if err != nil {
		log.Printf("Failed to fetch metadata content: %v", err)
		return
	}
	opr.Name = metadata.Name
	opr.Description = metadata.Description
	opr.LogoURI = metadata.Logo
	opr.IPAddress = metadata.IPAddress
	opr.Website = metadata.Website
	opr.Twitter = metadata.Twitter
	opr.BtcPublicKey = ""
	// Insert or update the operator in the database
	err = db.UpsertOperator(dbconn, opr)
	if err != nil {
		log.Printf("Failed to upsert operator: %v", err)
		return
	}
	fmt.Printf("  Operator %s upserted successfully\n", opr.EthAddress)
}

func HandleMotifOperatorRegistered(parsedABI abi.ABI, vLog etypes.Log) {
	dbconn := db.InitDB()
	defer dbconn.Close()
	// Decode the event
	event := struct {
		Operator     common.Address
		BtcPublicKey []byte
	}{}

	err := parsedABI.UnpackIntoInterface(&event, "OperatorBtcKeyRegistered", vLog.Data)
	if err != nil {
		log.Printf("Failed to unpack OperatorBtcKeyRegistered event: %v", err)
		return
	}

	// Log the event details
	fmt.Printf("OperatorBtcKeyRegistered:\n")
	fmt.Printf("  Operator: %s\n", event.Operator.Hex())
	fmt.Printf("  BTC Public Key: %x\n", event.BtcPublicKey)

	db.AddOrUpdateBtcPublicKey(dbconn, event.Operator.Hex(), string(event.BtcPublicKey))
}

func HandleOperatorBtckeyDeregistered(parsedABI abi.ABI, vLog etypes.Log) {
	dbconn := db.InitDB()
	defer dbconn.Close()
	// Decode the event
	event := struct {
		Operator common.Address
	}{}

	err := parsedABI.UnpackIntoInterface(&event, "OperatorBtckeyDeregistered", vLog.Data)
	if err != nil {
		log.Printf("Failed to unpack OperatorBtckeyDeregistered event: %v", err)
		return
	}

	// Log the event details
	fmt.Printf("OperatorBtckeyDeregistered:\n")
	fmt.Printf("  Operator: %s\n", event.Operator.Hex())

	// Example: Remove the BTC public key from the database
	err = db.DeleteFromMotif(database, event.Operator.Hex())
	if err != nil {
		log.Printf("Failed to remove BTC public key for operator: %v", err)
		return
	}
	fmt.Printf("  BTC public key removed for operator %s\n", event.Operator.Hex())
}
