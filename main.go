package main

import (
	"database/sql"
	"encoding/json"
	"log"
	"net/http"

	"github.com/gorilla/mux"
	_ "github.com/lib/pq"
)

type Operator struct {
	LogoURI    string `json:"logo_uri"`
	IPAddress  string `json:"ip_address"`
	Name       string `json:"name"`
	EthAddress string `json:"eth_address"`
}

var db *sql.DB

func main() {
	var err error
	db, err = sql.Open("postgres", "user=forkscanner dbname=forkscanner sslmode=disable password=forkscanner")
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	router := mux.NewRouter()

	router.HandleFunc("/operator", createOperator).Methods("POST")
	router.HandleFunc("/operators", getOperators).Methods("GET")

	log.Fatal(http.ListenAndServe(":8000", router))
}

func createOperator(w http.ResponseWriter, r *http.Request) {
	var operator Operator
	if err := json.NewDecoder(r.Body).Decode(&operator); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	_, err := db.Exec(`INSERT INTO operators (logo_uri, ip_address, name, eth_address) VALUES ($1, $2, $3, $4)`,
		operator.LogoURI, operator.IPAddress, operator.Name, operator.EthAddress)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{"status": "operator created"})
}

func getOperators(w http.ResponseWriter, r *http.Request) {
	rows, err := db.Query("SELECT logo_uri, ip_address, name, eth_address FROM operators")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	var operators []Operator
	for rows.Next() {
		var operator Operator
		if err := rows.Scan(&operator.LogoURI, &operator.IPAddress, &operator.Name, &operator.EthAddress); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		operators = append(operators, operator)
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(operators)
}
