package eventHandler

import (
	"database/sql"
	"encoding/hex"
	"fmt"
	"log"

	"github.com/AhmadAshraf2/motif-operator-listings/DelegationManager"
	"github.com/AhmadAshraf2/motif-operator-listings/MotifRegistry"
	"github.com/AhmadAshraf2/motif-operator-listings/db"
	"github.com/AhmadAshraf2/motif-operator-listings/types"
	"github.com/AhmadAshraf2/motif-operator-listings/utils"
)

var database *sql.DB

func HandleOperatorMetadataURIUpdatedEigenLayer(event *DelegationManager.DelegationManagerOperatorMetadataURIUpdated) {
	// Decode the event
	dbconn := db.InitDB()
	defer dbconn.Close()

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

func HandleMotifOperatorRegistered(event *MotifRegistry.MotifRegistryOperatorBtcKeyRegistered) {
	dbconn := db.InitDB()
	defer dbconn.Close()

	btcPubKey := hex.EncodeToString(event.BtcPublicKey[:])
	fmt.Printf("OperatorBtcKeyRegistered:\n")
	fmt.Printf("  Operator: %s\n", event.Operator.Hex())
	fmt.Printf("  BTC Public Key: %x\n", event.BtcPublicKey)
	err := db.AddOrUpdateBtcPublicKey(dbconn, event.Operator.Hex(), btcPubKey)
	if err != nil {
		fmt.Println(err)
	}
}

func HandleMotifOperatorDeregistered(event *MotifRegistry.MotifRegistryOperatorBtckeyDeregistered) {
	dbconn := db.InitDB()
	defer dbconn.Close()
	// Log the event details
	fmt.Printf("OperatorBtckeyDeregistered:\n")
	fmt.Printf("  Operator: %s\n", event.Operator.Hex())

	err := db.DeleteFromMotif(database, event.Operator.Hex())
	if err != nil {
		log.Printf("Failed to remove BTC public key for operator: %v", err)
		return
	}
	fmt.Printf("  BTC public key removed for operator %s\n", event.Operator.Hex())
}
