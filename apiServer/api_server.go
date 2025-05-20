package apiServer

import (
	"database/sql"
	"encoding/json"
	"log"
	"net/http"

	"github.com/AhmadAshraf2/motif-operator-listings/db"
	"github.com/gorilla/mux"
)

var database *sql.DB

func StartApiServer() {
	router := mux.NewRouter()
	router.HandleFunc("/operators", getOperators).Methods("GET")
	log.Fatal(http.ListenAndServe(":8000", router))
}

func getOperators(w http.ResponseWriter, r *http.Request) {
	dbconn := db.InitDB()
	defer dbconn.Close()
	operators, err := db.GetOperators(dbconn)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(operators)
}
