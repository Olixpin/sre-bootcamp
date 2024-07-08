package handlers

import (
	"context"
	"log"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
	"github.com/olixpin/student-api/student-api/data"
)

type Students struct {
	l *log.Logger
}

func NewStudents(l *log.Logger) *Students {
	return &Students{l}
}

func (s *Students) GetStudents(rw http.ResponseWriter, r *http.Request) {
	s.l.Println("Handle GET Students")

	ls := data.GetStudents()

	err := ls.ToJSON(rw)
	if err != nil {
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
	student, err := data.GetStudentByID(uint(id))

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

	student := r.Context().Value(KeyStudent{}).(*data.Student)
	data.AddStudent(student)
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
	student := r.Context().Value(KeyStudent{}).(*data.Student)

	err = data.UpdateStudent(uint(id), student)
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

	err = data.DeleteStudent(uint(id))
	if err == data.ErrStudentNotFound {
		http.Error(rw, "Student not found", http.StatusNotFound)
		return
	} else if err != nil {
		http.Error(rw, "Unable to delete student", http.StatusInternalServerError)
		return
	}

	rw.WriteHeader(http.StatusNoContent)
}

type KeyStudent struct{}

func (s *Students) MiddlewareStudentValidation(next http.Handler) http.Handler {
	return http.HandlerFunc(func(rw http.ResponseWriter, r *http.Request) {
		student := &data.Student{}

		err := student.FromJSON(r.Body)
		if err != nil {
			http.Error(rw, "Unable to unmarshal json", http.StatusBadRequest)
			return
		}

		err = student.Validate()
		if err != nil {
			http.Error(rw, err.Error(), http.StatusBadRequest)
			return
		}

		ctx := context.WithValue(r.Context(), KeyStudent{}, student)
		r = r.WithContext(ctx)

		next.ServeHTTP(rw, r)
	})
}
