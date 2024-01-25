package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/rs/cors"
)

// Item es una estructura para representar los datos
type Item struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
}

func main() {
	router := mux.NewRouter()
	router.HandleFunc("/api/data", GetData).Methods("GET")

	corsHandler := cors.Default().Handler(router)

	fmt.Println("Server is running on :8080")
	log.Fatal(http.ListenAndServe(":8080", corsHandler))
}

func GetData(w http.ResponseWriter, r *http.Request) {
	rows, err := Db.Query("SELECT id, name FROM data")
	if err != nil {
		http.Error(w, "Error querying the database", http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	var items []Item

	for rows.Next() {
		var item Item
		if err := rows.Scan(&item.ID, &item.Name); err != nil {
			http.Error(w, "Error scanning rows", http.StatusInternalServerError)
			return
		}
		items = append(items, item)
	}

	responseJSON, err := json.Marshal(items)
	if err != nil {
		http.Error(w, "Error converting data to JSON", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(responseJSON)
}
