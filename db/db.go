package db

import (
	"database/sql"
	"fmt"
	"log"

	"github.com/AhmadAshraf2/motif-operator-listings/types"
	"github.com/spf13/viper"
)

func InitDB() *sql.DB {
	psqlconn := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable", viper.Get("DB_host"), viper.Get("DB_port"), viper.Get("DB_user"), viper.Get("DB_password"), viper.Get("DB_name"))
	db, err := sql.Open("postgres", psqlconn)
	if err != nil {
		log.Println("DB error : ", err)
		panic(err)
	}
	fmt.Println("DB initialized")
	return db
}

func InsertOperator(db *sql.DB, operator types.Operator) error {
	query := `INSERT INTO operators (
        eth_address, name, description, logo_uri, ip_address, btc_public_key, website, twitter
    ) VALUES ($1, $2, $3, $4, $5, $6, $7, $8)`

	_, err := db.Exec(query, operator.EthAddress, operator.Name, operator.Description, operator.LogoURI, operator.IPAddress, operator.BtcPublicKey, operator.Website, operator.Twitter)
	if err != nil {
		return fmt.Errorf("failed to insert operator: %w", err)
	}
	return nil
}

func GetOperators(db *sql.DB) ([]types.Operator, error) {
	query := `SELECT eth_address, name, description, logo_uri, ip_address, btc_public_key, website, twitter FROM operators where is_motif = true`

	rows, err := db.Query(query)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch operators: %w", err)
	}
	defer rows.Close()

	var operators []types.Operator
	for rows.Next() {
		var operator types.Operator
		err := rows.Scan(&operator.EthAddress, &operator.Name, &operator.Description, &operator.LogoURI, &operator.IPAddress, &operator.BtcPublicKey, &operator.Website, &operator.Twitter)
		if err != nil {
			return nil, fmt.Errorf("failed to scan operator: %w", err)
		}
		operators = append(operators, operator)
	}

	return operators, nil
}

func UpsertOperator(db *sql.DB, operator types.Operator) error {
	queryCheck := `SELECT COUNT(*) FROM operators WHERE eth_address = $1`
	var count int
	err := db.QueryRow(queryCheck, operator.EthAddress).Scan(&count)
	if err != nil {
		return fmt.Errorf("failed to check if eth_address exists: %w", err)
	}

	if count > 0 {
		// Update the existing entry
		queryUpdate := `UPDATE operators 
            SET name = $1, description = $2, logo_uri = $3, ip_address = $4, website = $5, twitter = $6
            WHERE eth_address = $7`
		_, err := db.Exec(queryUpdate, operator.Name, operator.Description, operator.LogoURI, operator.IPAddress, operator.Website, operator.Twitter, operator.EthAddress)
		if err != nil {
			return fmt.Errorf("failed to update operator: %w", err)
		}
	} else {
		// Insert a new entry
		queryInsert := `INSERT INTO operators (
            eth_address, name, description, logo_uri, ip_address, btc_public_key, website, twitter,
        ) VALUES ($1, $2, $3, $4, $5, $6, $7)`
		_, err := db.Exec(queryInsert, operator.EthAddress, operator.Name, operator.Description, operator.LogoURI, operator.IPAddress, operator.BtcPublicKey, operator.Website, operator.Twitter)
		if err != nil {
			return fmt.Errorf("failed to insert operator: %w", err)
		}
	}

	return nil
}

func AddOrUpdateBtcPublicKey(db *sql.DB, ethAddress string, btcPublicKey string) error {
	queryCheck := `SELECT COUNT(*) FROM operators WHERE eth_address = $1`
	var count int
	err := db.QueryRow(queryCheck, ethAddress).Scan(&count)
	if err != nil {
		return fmt.Errorf("failed to check if eth_address exists: %w", err)
	}

	if count > 0 {
		// Update the BTC public key for the existing entry
		queryUpdate := `UPDATE operators SET btc_public_key = $1, is_motif = true WHERE eth_address = $2`
		_, err := db.Exec(queryUpdate, btcPublicKey, ethAddress)
		if err != nil {
			return fmt.Errorf("failed to update BTC public key: %w", err)
		}
	} else {
		// Insert a new entry with the BTC public key
		queryInsert := `INSERT INTO operators (eth_address, btc_public_key, is_motif) VALUES ($1, $2, true)`
		_, err := db.Exec(queryInsert, ethAddress, btcPublicKey)
		if err != nil {
			return fmt.Errorf("failed to insert BTC public key: %w", err)
		}
	}

	return nil
}

func DeleteFromMotif(db *sql.DB, ethAddress string) error {
	query := `UPDATE operators SET is_motif = false WHERE eth_address = $1`

	_, err := db.Exec(query, ethAddress)
	if err != nil {
		return fmt.Errorf("failed to update is_motif to false for eth_address %s: %w", ethAddress, err)
	}

	return nil
}
