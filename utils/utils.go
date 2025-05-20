package utils

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"

	"github.com/AhmadAshraf2/motif-operator-listings/types"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/spf13/viper"
)

func FetchMetadataContent(uri string) (*types.Metadata, error) {
	resp, err := http.Get(uri)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch content from URL: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("failed to fetch content, HTTP status: %s", resp.Status)
	}

	content, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read content from URL: %w", err)
	}

	var metadata types.Metadata
	err = json.Unmarshal(content, &metadata)
	if err != nil {
		return nil, fmt.Errorf("failed to parse JSON content: %w", err)
	}

	return &metadata, nil
}

func ReadABIFromFile(filePath string) (string, error) {
	// Open the file
	file, err := os.Open(filePath)
	if err != nil {
		return "", fmt.Errorf("failed to open file: %w", err)
	}
	defer file.Close()

	// Read the file content
	content, err := ioutil.ReadAll(file)
	if err != nil {
		return "", fmt.Errorf("failed to read file: %w", err)
	}

	return string(content), nil
}

func GetEthClient() (*ethclient.Client, error) {
	client, err := ethclient.Dial(viper.GetString("eth_rpc_host"))
	if err != nil {
		fmt.Println("Failed to connect to the Ethereum client: %v", err)
		return nil, err
	}
	return client, nil
}
