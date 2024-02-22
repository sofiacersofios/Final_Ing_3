// main_test.go
package main

import (
	"bytes"
	"encoding/json"
	"log"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
)

func TestGetData(t *testing.T) {
	log.Println("Running TestGetData")
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatal(err)
	}
	defer db.Close()

	Db = db

	rows := sqlmock.NewRows([]string{"id", "name"}).AddRow(1, "example") // Ajusta según los datos de prueba
	mock.ExpectQuery("SELECT id, name FROM data").WillReturnRows(rows)

	req, err := http.NewRequest("GET", "/api/data", nil)
	if err != nil {
		t.Fatal(err)
	}

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(GetData)

	handler.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusOK {
		t.Errorf("Handler returned wrong status code: got %v want %v",
			status, http.StatusOK)
	}

	log.Printf("Status Code: %d", rr.Code)
	log.Printf("Response Body: %s", rr.Body.String())

	// Verifica la respuesta
	expected := `[{"id":1,"name":"example"}]` // Ajusta según los datos de prueba
	if rr.Body.String() != expected {
		t.Errorf("Handler returned unexpected body: got %v want %v",
			rr.Body.String(), expected)
	}
}

func TestCreateData(t *testing.T) {
	log.Println("Running TestCreateData")

	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatal(err)
	}
	defer db.Close()

	Db = db

	newItem := Item{ID: 1, Name: "test"}
	jsonData, err := json.Marshal(newItem)
	if err != nil {
		t.Fatal(err)
	}

	// Expect INSERT query and return a result with ID 1
	mock.ExpectExec("INSERT INTO data").WithArgs(newItem.Name).WillReturnResult(sqlmock.NewResult(1, 1)) //aca se pone el id

	req, err := http.NewRequest("POST", "/api/data", bytes.NewBuffer(jsonData))
	if err != nil {
		t.Fatal(err)
	}

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(CreateData)

	handler.ServeHTTP(rr, req)

	log.Printf("Status Code: %d", rr.Code)
	log.Printf("Response Body: %s", rr.Body.String())

	if status := rr.Code; status != http.StatusCreated {
		t.Errorf("Handler returned wrong status code: got %v want %v",
			status, http.StatusCreated)
	}

	// Verificar que el cuerpo de la respuesta contiene el ID
	var response map[string]int
	err = json.Unmarshal(rr.Body.Bytes(), &response)
	if err != nil {
		t.Fatal(err)
	}

	createdID := response["id"]

	if createdID != 1 {
		t.Errorf("Expected created ID to be 1, but got %v", createdID)
	}
}

/*
func TestDeleteData(t *testing.T) {
	log.Println("Running TestDeleteData")
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatal(err)
	}
	defer db.Close()

	Db = db

	// Ajusta según los datos de prueba
	itemID := "1"

	// Expectativa de la consulta DELETE
	mock.ExpectExec("DELETE FROM data WHERE id = ?").WithArgs(itemID).WillReturnResult(sqlmock.NewResult(0, 1))

	req, err := http.NewRequest("DELETE", fmt.Sprintf("/api/data/%s", itemID), nil)
	if err != nil {
		t.Fatal(err)
	}

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(DeleteData)

	handler.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusOK {
		t.Errorf("Handler returned wrong status code: got %v want %v",
			status, http.StatusOK)
	}

	log.Printf("Status Code: %d", rr.Code)
	log.Printf("Response Body: %s", rr.Body.String())

	// Verificar que se llamó a la consulta DELETE con el ID esperado
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("Expectations not met: %s", err)
	}
}
*/
