package main

import (
	"context"
	"database/sql"
	"log"
	"net/http"
	"os"
	"os/signal"
	"time"

	"github.com/gorilla/mux"
	_ "github.com/lib/pq"
	"github.com/nicholasjackson/env"
	"github.com/olixpin/student-api/handlers"
)

var (
	bindAddress = env.String("BIND_ADDRESS", false, ":8080", "Bind address for the server")
	dbURL       = env.String("DB_URL", false, "", "Database connection URL")
)

func main() {
	env.Parse()

	l := log.New(os.Stdout, "students-api ", log.LstdFlags)

	// Test database connection
	db, err := sql.Open("postgres", *dbURL)
	if err != nil {
		l.Fatalf("Error opening database: %v", err)
	}
	defer db.Close()

	err = db.Ping()
	if err != nil {
		l.Fatalf("Error connecting to database: %v", err)
	}
	l.Println("Successfully connected to the database")

	sh := handlers.NewStudents(l, db)
	sm := mux.NewRouter()

	sm.HandleFunc("/api/v1/students", sh.GetStudents).Methods(http.MethodGet)
	sm.HandleFunc("/api/v1/students/{id:[0-9]+}", sh.GetStudentByID).Methods(http.MethodGet)
	sm.HandleFunc("/api/v1/students", sh.AddStudent).Methods(http.MethodPost)
	sm.HandleFunc("/api/v1/students/{id:[0-9]+}", sh.UpdateStudent).Methods(http.MethodPut)
	sm.HandleFunc("/api/v1/students/{id:[0-9]+}", sh.DeleteStudent).Methods(http.MethodDelete)
	sm.HandleFunc("/healthcheck", sh.HealthCheck).Methods(http.MethodGet)

	s := http.Server{
		Addr:         *bindAddress,
		Handler:      sm,
		ErrorLog:     l,
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 10 * time.Second,
		IdleTimeout:  120 * time.Second,
	}

	go func() {
		l.Println("Starting server on port", *bindAddress)
		if err := s.ListenAndServe(); err != nil {
			l.Printf("Error starting server: %s\n", err)
			os.Exit(1)
		}
	}()

	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)
	signal.Notify(c, os.Kill)

	sig := <-c
	l.Println("Got signal:", sig)

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	s.Shutdown(ctx)
}
