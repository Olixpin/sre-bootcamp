package handlers

import (
	"bytes"
	"database/sql"
	"encoding/json"
	"log"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"
	"time"

	"github.com/gorilla/mux"
	_ "github.com/mattn/go-sqlite3"
	"github.com/olixpin/student-api/student-api/data"
)

func setupTestDB(t *testing.T) *sql.DB {
	db, err := sql.Open("sqlite3", ":memory:")
	if err != nil {
		t.Fatal(err)
	}

	createTable := `
	CREATE TABLE students (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		first_name TEXT NOT NULL,
		last_name TEXT NOT NULL,
		age INTEGER NOT NULL,
		email TEXT NOT NULL,
		enrollment_date TIMESTAMP NOT NULL,
		class TEXT NOT NULL,
		address TEXT NOT NULL,
		phone_number TEXT NOT NULL,
		created_on TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
		updated_on TIMESTAMP DEFAULT CURRENT_TIMESTAMP
	);`

	_, err = db.Exec(createTable)
	if err != nil {
		t.Fatal(err)
	}

	return db
}

func TestGetStudents(t *testing.T) {
	db := setupTestDB(t)
	defer db.Close()

	l := log.New(os.Stdout, "students-api-test ", log.LstdFlags)
	studentsHandler := NewStudents(l, db)

	req, err := http.NewRequest("GET", "/api/v1/students", nil)
	if err != nil {
		t.Fatal(err)
	}

	rr := httptest.NewRecorder()
	router := mux.NewRouter()
	router.HandleFunc("/api/v1/students", studentsHandler.GetStudents)
	router.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusOK)
	}

	var students data.Students
	err = json.NewDecoder(rr.Body).Decode(&students)
	if err != nil {
		t.Errorf("handler returned invalid body: %v", err)
	}
}

func TestAddStudent(t *testing.T) {
	db := setupTestDB(t)
	defer db.Close()

	l := log.New(os.Stdout, "students-api-test ", log.LstdFlags)
	studentsHandler := NewStudents(l, db)

	newStudent := &data.Student{
		FirstName:      "Test",
		LastName:       "Student",
		Age:            21,
		Email:          "test.student@example.com",
		EnrollmentDate: time.Now(),
		Class:          "Test Class",
		Address:        "Test Address",
		PhoneNumber:    "1234567890",
	}
	body, _ := json.Marshal(newStudent)

	req, err := http.NewRequest("POST", "/api/v1/students", bytes.NewBuffer(body))
	if err != nil {
		t.Fatal(err)
	}
	req.Header.Set("Content-Type", "application/json")

	rr := httptest.NewRecorder()
	router := mux.NewRouter()
	router.HandleFunc("/api/v1/students", studentsHandler.AddStudent)
	router.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusCreated {
		t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusCreated)
	}
}

func TestUpdateStudent(t *testing.T) {
	db := setupTestDB(t)
	defer db.Close()

	l := log.New(os.Stdout, "students-api-test ", log.LstdFlags)
	studentsHandler := NewStudents(l, db)

	// Insert a student to update
	newStudent := &data.Student{
		FirstName:      "Test",
		LastName:       "Student",
		Age:            21,
		Email:          "test.student@example.com",
		EnrollmentDate: time.Now(),
		Class:          "Test Class",
		Address:        "Test Address",
		PhoneNumber:    "1234567890",
	}
	_, err := db.Exec("INSERT INTO students (first_name, last_name, age, email, enrollment_date, class, address, phone_number) VALUES (?, ?, ?, ?, ?, ?, ?, ?)",
		newStudent.FirstName, newStudent.LastName, newStudent.Age, newStudent.Email, newStudent.EnrollmentDate, newStudent.Class, newStudent.Address, newStudent.PhoneNumber)
	if err != nil {
		t.Fatal(err)
	}

	updatedStudent := &data.Student{
		ID:             1,
		FirstName:      "Updated",
		LastName:       "Student",
		Age:            22,
		Email:          "updated.student@example.com",
		EnrollmentDate: time.Now(),
		Class:          "Updated Class",
		Address:        "Updated Address",
		PhoneNumber:    "0987654321",
	}
	body, _ := json.Marshal(updatedStudent)

	req, err := http.NewRequest("PUT", "/api/v1/students/1", bytes.NewBuffer(body))
	if err != nil {
		t.Fatal(err)
	}
	req.Header.Set("Content-Type", "application/json")

	rr := httptest.NewRecorder()
	router := mux.NewRouter()
	router.HandleFunc("/api/v1/students/{id:[0-9]+}", studentsHandler.UpdateStudent)
	router.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusNoContent {
		t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusNoContent)
	}
}

func TestDeleteStudent(t *testing.T) {
	db := setupTestDB(t)
	defer db.Close()

	l := log.New(os.Stdout, "students-api-test ", log.LstdFlags)
	studentsHandler := NewStudents(l, db)

	// Insert a student to delete
	_, err := db.Exec("INSERT INTO students (first_name, last_name, age, email, enrollment_date, class, address, phone_number) VALUES (?, ?, ?, ?, ?, ?, ?, ?)",
		"Test", "Student", 21, "test.student@example.com", time.Now(), "Test Class", "Test Address", "1234567890")
	if err != nil {
		t.Fatal(err)
	}

	req, err := http.NewRequest("DELETE", "/api/v1/students/1", nil)
	if err != nil {
		t.Fatal(err)
	}

	rr := httptest.NewRecorder()
	router := mux.NewRouter()
	router.HandleFunc("/api/v1/students/{id:[0-9]+}", studentsHandler.DeleteStudent)
	router.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusNoContent {
		t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusNoContent)
	}
}

func TestHealthCheck(t *testing.T) {
	db := setupTestDB(t)
	defer db.Close()

	l := log.New(os.Stdout, "students-api-test ", log.LstdFlags)
	studentsHandler := NewStudents(l, db)

	req, err := http.NewRequest("GET", "/healthcheck", nil)
	if err != nil {
		t.Fatal(err)
	}

	rr := httptest.NewRecorder()
	router := mux.NewRouter()
	router.HandleFunc("/healthcheck", studentsHandler.HealthCheck)
	router.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusOK)
	}

	expected := "OK"
	if rr.Body.String() != expected {
		t.Errorf("handler returned unexpected body: got %v want %v", rr.Body.String(), expected)
	}
}
