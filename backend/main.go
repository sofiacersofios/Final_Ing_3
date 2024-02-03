package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strings"

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
	router.HandleFunc("/api/data", CreateData).Methods("POST")
	router.HandleFunc("/api/data/{id}", DeleteData).Methods("DELETE")

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

func CreateData(w http.ResponseWriter, r *http.Request) {
	var newItem Item
	err := json.NewDecoder(r.Body).Decode(&newItem)
	if err != nil {
		http.Error(w, "Error decoding request body", http.StatusBadRequest)
		return
	}

	result, err := Db.Exec("INSERT INTO data (name) VALUES (?)", newItem.Name)
	if err != nil {
		http.Error(w, "Error inserting data into the database", http.StatusInternalServerError)
		return
	}

	// Obtener el ID del elemento recién insertado
	createdID, err := result.LastInsertId()
	if err != nil {
		http.Error(w, "Error getting last insert ID", http.StatusInternalServerError)
		return
	}

	// Crear la respuesta con el ID
	response := map[string]int{"id": int(createdID)}

	// Convertir a JSON y enviar la respuesta
	jsonResponse, err := json.Marshal(response)
	if err != nil {
		http.Error(w, "Error converting response to JSON", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	w.Write(jsonResponse)
}

func DeleteData(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	itemID, ok := vars["id"]
	if !ok {
		http.Error(w, "Missing item ID", http.StatusBadRequest)
		return
	}

	_, err := Db.Exec("DELETE FROM data WHERE id = ?", itemID)
	if err != nil {
		// Verifica si el error es debido a que el elemento no existe
		if strings.Contains(err.Error(), "no rows in result set") {
			http.Error(w, "Item not found", http.StatusNotFound)
			return
		}

		http.Error(w, "Error deleting data from the database", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}
