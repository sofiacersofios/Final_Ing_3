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

	InitDB("root:root@tcp(localhost:3306)/mydatabase")

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

	InitDB("root:root@tcp(localhost:3306)/mydatabase")

	Db = db

	newItem := Item{ID: 1, Name: "test"}
	jsonData, err := json.Marshal(newItem)
	if err != nil {
		t.Fatal(err)
	}

	// Expect INSERT query and return a result with ID 1
	mock.ExpectExec("INSERT INTO data").WithArgs(newItem.Name).WillReturnResult(sqlmock.NewResult(1, 1))

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

	// Configuración de la base de datos de prueba
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatal(err)
	}
	defer db.Close()

	// Configura la variable global Db con la base de datos de prueba
	Db = db

	// Configura la expectativa para la consulta DELETE cuando el elemento no existe
	mock.ExpectExec("DELETE FROM data WHERE id = ?").
		WithArgs(1).
		WillReturnResult(sqlmock.NewResult(0, 0))

	// Realiza la solicitud para borrar el elemento con ID conocido (1)
	reqDelete, err := http.NewRequest("DELETE", "/api/data/1", nil)
	if err != nil {
		t.Fatal(err)
	}

	rrDelete := httptest.NewRecorder()
	handlerDelete := http.HandlerFunc(DeleteData)
	handlerDelete.ServeHTTP(rrDelete, reqDelete)

	// Verifica el código de estado de la respuesta HTTP para el caso de elemento no encontrado
	if status := rrDelete.Code; status != http.StatusNotFound {
		t.Errorf("Handler for deleting data should return http.StatusNotFound for missing item: got %v want %v",
			status, http.StatusNotFound)
	}

	// Verifica el cuerpo de la respuesta HTTP para el caso de elemento no encontrado
	expectedBodyNotFound := "Item not found"
	if bodyNotFound := strings.TrimSpace(rrDelete.Body.String()); bodyNotFound != expectedBodyNotFound {
		t.Errorf("Handler for deleting data should not return a body, but got %v", bodyNotFound)
	}

	// Verifica que se cumplan todas las expectativas de sqlmock
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("SQL expectations were not met: %v", err)
	}

	// Reinicia las expectativas de sqlmock para el caso en que el elemento existe
	mock.ExpectExec("DELETE FROM data WHERE id = ?").
		WithArgs(1).
		WillReturnResult(sqlmock.NewResult(0, 1))

	// Realiza la solicitud para borrar el elemento con ID conocido (1)
	reqDeleteExisting, err := http.NewRequest("DELETE", "/api/data/1", nil)
	if err != nil {
		t.Fatal(err)
	}

	rrDeleteExisting := httptest.NewRecorder()
	handlerDeleteExisting := http.HandlerFunc(DeleteData)
	handlerDeleteExisting.ServeHTTP(rrDeleteExisting, reqDeleteExisting)

	// Verifica el código de estado de la respuesta HTTP para el caso de elemento existente
	if statusExisting := rrDeleteExisting.Code; statusExisting != http.StatusOK {
		t.Errorf("Handler for deleting data returned wrong status code: got %v want %v",
			statusExisting, http.StatusOK)
	}

	// Verifica que no haya cuerpo de respuesta para el caso de elemento existente
	if bodyExisting := strings.TrimSpace(rrDeleteExisting.Body.String()); bodyExisting != "" {
		t.Errorf("Handler for deleting data should not return a body, but got %v", bodyExisting)
	}

	// Verifica que se cumplan todas las expectativas de sqlmock para el caso de elemento existente
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("SQL expectations were not met: %v", err)
	}
}
*/
