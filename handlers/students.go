package handlers

import (
	"database/sql"
	"log"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
	_ "github.com/lib/pq"
	"github.com/olixpin/student-api/student-api/data"
)

type Students struct {
	l  *log.Logger
	db *sql.DB
}

func NewStudents(l *log.Logger, db *sql.DB) *Students {
	return &Students{l, db}
}

func (s *Students) GetStudents(rw http.ResponseWriter, r *http.Request) {
	s.l.Println("Handle GET Students")

	students, err := data.GetStudents(s.db)
	if err != nil {
		s.l.Printf("Error fetching students: %v", err)
		http.Error(rw, "Unable to fetch students", http.StatusInternalServerError)
		return
	}

	err = students.ToJSON(rw)
	if err != nil {
		s.l.Printf("Error marshalling students: %v", err)
		http.Error(rw, "Unable to marshal json", http.StatusInternalServerError)
	}
}

func (s *Students) GetStudentByID(rw http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, err := strconv.Atoi(vars["id"])
	if err != nil {
		http.Error(rw, "Unable to convert id", http.StatusBadRequest)
		return
	}

	s.l.Println("Handle GET Student", id)
	student, err := data.GetStudentByID(s.db, uint(id))

	if err == data.ErrStudentNotFound {
		http.Error(rw, "Student not found", http.StatusNotFound)
		return
	} else if err != nil {
		http.Error(rw, "Unable to get student", http.StatusInternalServerError)
		return
	}

	err = student.ToJSON(rw)
	if err != nil {
		http.Error(rw, "Unable to marshal json", http.StatusInternalServerError)
	}
}

func (s *Students) AddStudent(rw http.ResponseWriter, r *http.Request) {
	s.l.Println("Handle POST Students")

	student := &data.Student{}
	err := student.FromJSON(r.Body)
	if err != nil {
		http.Error(rw, "Unable to unmarshal json", http.StatusBadRequest)
		return
	}

	err = data.AddStudent(s.db, student)
	if err != nil {
		http.Error(rw, "Unable to add student", http.StatusInternalServerError)
		return
	}

	rw.WriteHeader(http.StatusCreated)
}

func (s *Students) UpdateStudent(rw http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, err := strconv.Atoi(vars["id"])
	if err != nil {
		http.Error(rw, "Unable to convert id", http.StatusBadRequest)
		return
	}

	s.l.Println("Handle PUT Student", id)
	student := &data.Student{}
	err = student.FromJSON(r.Body)
	if err != nil {
		http.Error(rw, "Unable to unmarshal json", http.StatusBadRequest)
		return
	}

	err = data.UpdateStudent(s.db, uint(id), student)
	if err == data.ErrStudentNotFound {
		http.Error(rw, "Student not found", http.StatusNotFound)
		return
	} else if err != nil {
		http.Error(rw, "Unable to update student", http.StatusInternalServerError)
		return
	}

	rw.WriteHeader(http.StatusNoContent)
}

func (s *Students) DeleteStudent(rw http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, err := strconv.Atoi(vars["id"])
	if err != nil {
		http.Error(rw, "Unable to convert id", http.StatusBadRequest)
		return
	}

	s.l.Println("Handle DELETE Student", id)

	err = data.DeleteStudent(s.db, uint(id))
	if err == data.ErrStudentNotFound {
		http.Error(rw, "Student not found", http.StatusNotFound)
		return
	} else if err != nil {
		http.Error(rw, "Unable to delete student", http.StatusInternalServerError)
		return
	}

	rw.WriteHeader(http.StatusNoContent)
}

func (s *Students) HealthCheck(rw http.ResponseWriter, r *http.Request) {
	s.l.Println("Handle Health Check")

	err := s.db.Ping()
	if err != nil {
		s.l.Printf("Database connection error: %v", err)
		http.Error(rw, "Database connection error", http.StatusInternalServerError)
		return
	}

	rw.WriteHeader(http.StatusOK)
	rw.Write([]byte("OK"))
}
