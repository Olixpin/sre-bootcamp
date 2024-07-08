package data

import (
	"database/sql"
	"encoding/json"
	"errors"
	"io"
	"log"
	"time"
)

type Student struct {
	ID             uint      `json:"id"`
	FirstName      string    `json:"first_name"`
	LastName       string    `json:"last_name"`
	Age            int       `json:"age"`
	Email          string    `json:"email"`
	EnrollmentDate time.Time `json:"enrollment_date"`
	Class          string    `json:"class"`
	Address        string    `json:"address"`
	PhoneNumber    string    `json:"phone_number"`
	CreatedOn      time.Time `json:"-"`
	UpdatedOn      time.Time `json:"-"`
}

func (s *Student) FromJSON(r io.Reader) error {
	e := json.NewDecoder(r)
	return e.Decode(s)
}

func (s *Student) ToJSON(w io.Writer) error {
	e := json.NewEncoder(w)
	return e.Encode(s)
}

type Students []*Student

func (s *Students) ToJSON(w io.Writer) error {
	e := json.NewEncoder(w)
	return e.Encode(s)
}

var ErrStudentNotFound = errors.New("student not found")

func GetStudents(db *sql.DB) (Students, error) {
	stmt := `SELECT id, first_name, last_name, age, email, enrollment_date, class, address, phone_number FROM students`
	rows, err := db.Query(stmt)
	if err != nil {
		log.Printf("Error executing query: %v", err)
		return nil, err
	}
	defer rows.Close()

	students := Students{}
	for rows.Next() {
		var student Student
		err := rows.Scan(&student.ID, &student.FirstName, &student.LastName, &student.Age, &student.Email, &student.EnrollmentDate, &student.Class, &student.Address, &student.PhoneNumber)
		if err != nil {
			log.Printf("Error scanning row: %v", err)
			return nil, err
		}
		students = append(students, &student)
	}
	return students, nil
}

func GetStudentByID(db *sql.DB, id uint) (*Student, error) {
	stmt := `SELECT id, first_name, last_name, age, email, enrollment_date, class, address, phone_number FROM students WHERE id = $1`
	row := db.QueryRow(stmt, id)

	var student Student
	err := row.Scan(&student.ID, &student.FirstName, &student.LastName, &student.Age, &student.Email, &student.EnrollmentDate, &student.Class, &student.Address, &student.PhoneNumber)
	if err == sql.ErrNoRows {
		return nil, ErrStudentNotFound
	} else if err != nil {
		log.Printf("Error scanning row: %v", err)
		return nil, err
	}
	return &student, nil
}

func AddStudent(db *sql.DB, student *Student) error {
	stmt := `INSERT INTO students (first_name, last_name, age, email, enrollment_date, class, address, phone_number) VALUES ($1, $2, $3, $4, $5, $6, $7, $8) RETURNING id`
	err := db.QueryRow(stmt, student.FirstName, student.LastName, student.Age, student.Email, student.EnrollmentDate, student.Class, student.Address, student.PhoneNumber).Scan(&student.ID)
	if err != nil {
		log.Printf("Error executing insert: %v", err)
	}
	return err
}

func UpdateStudent(db *sql.DB, id uint, student *Student) error {
	stmt := `UPDATE students SET first_name = $1, last_name = $2, age = $3, email = $4, enrollment_date = $5, class = $6, address = $7, phone_number = $8, updated_on = $9 WHERE id = $10`
	res, err := db.Exec(stmt, student.FirstName, student.LastName, student.Age, student.Email, student.EnrollmentDate, student.Class, student.Address, student.PhoneNumber, time.Now(), id)
	if err != nil {
		log.Printf("Error executing update: %v", err)
		return err
	}

	rowsAffected, err := res.RowsAffected()
	if err != nil {
		log.Printf("Error getting rows affected: %v", err)
		return err
	}
	if rowsAffected == 0 {
		return ErrStudentNotFound
	}
	return nil
}

func DeleteStudent(db *sql.DB, id uint) error {
	stmt := `DELETE FROM students WHERE id = $1`
	res, err := db.Exec(stmt, id)
	if err != nil {
		log.Printf("Error executing delete: %v", err)
		return err
	}

	rowsAffected, err := res.RowsAffected()
	if err != nil {
		log.Printf("Error getting rows affected: %v", err)
		return err
	}
	if rowsAffected == 0 {
		return ErrStudentNotFound
	}
	return nil
}
