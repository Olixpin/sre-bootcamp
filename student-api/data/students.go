package data

import (
	"encoding/json"
	"errors"
	"io"
	"time"
)

type Student struct {
	ID             uint      `json:"id" gorm:"primary_key"`
	FirstName      string    `json:"first_name"`
	LastName       string    `json:"last_name"`
	Age            int       `json:"age"`
	Email          string    `json:"email"`
	EnrollmentDate time.Time `json:"enrollment_date"`
	Class          string    `json:"class"`
	Address        string    `json:"address"`
	PhoneNumber    string    `json:"phone_number"`
	CreatedOn      string    `json:"-"`
	UpdatedOn      string    `json:"-"`
	DeletedOn      string    `json:"-"`
}

func (p *Student) FromJSON(r io.Reader) error {
	e := json.NewDecoder(r)
	return e.Decode(p)
}

func (p *Student) ToJSON(w io.Writer) error {
	e := json.NewEncoder(w)
	return e.Encode(p)
}

func (p *Student) Validate() error {
	// Implement validation logic
	return nil
}

type Students []*Student

func (s *Students) ToJSON(w io.Writer) error {
	enc := json.NewEncoder(w)
	return enc.Encode(s)
}

func (s *Students) FromJSON(r io.Reader) error {
	e := json.NewDecoder(r)
	return e.Decode(s)
}

func GetStudents() Students {
	return studentList
}

func GetStudentByID(id uint) (*Student, error) {
	student, _, err := findStudent(id)
	return student, err
}

func AddStudent(s *Student) {
	s.ID = getNextID()
	studentList = append(studentList, s)
}

func UpdateStudent(id uint, s *Student) error {
	_, i, err := findStudent(id)
	if err != nil {
		return err
	}
	studentList[i] = s
	return nil
}

func DeleteStudent(id uint) error {
	_, i, err := findStudent(id)
	if err != nil {
		return err
	}
	studentList = append(studentList[:i], studentList[i+1:]...)
	return nil
}

var ErrStudentNotFound = errors.New("student not found")

func findStudent(id uint) (*Student, int, error) {
	for i, s := range studentList {
		if s.ID == id {
			return s, i, nil
		}
	}
	return nil, -1, ErrStudentNotFound
}

func getNextID() uint {
	if len(studentList) == 0 {
		return 1
	}
	return studentList[len(studentList)-1].ID + 1
}

var studentList = Students{
	{
		ID:             1,
		FirstName:      "John",
		LastName:       "Doe",
		Age:            20,
		Email:          "john.doe@example.com",
		EnrollmentDate: time.Now(),
		Class:          "A",
		Address:        "123 Main St, Anytown USA",
		PhoneNumber:    "555-1234",
	},
	{
		ID:             2,
		FirstName:      "Jane",
		LastName:       "Smith",
		Age:            22,
		Email:          "jane.smith@example.com",
		EnrollmentDate: time.Now(),
		Class:          "B",
		Address:        "456 Oak Rd, Anytown USA",
		PhoneNumber:    "555-5678",
	},
}
