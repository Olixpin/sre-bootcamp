package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"time"

	"github.com/gorilla/mux"
	"github.com/olixpin/student-api/handlers"
)

var bindAddress = ":8080"

func main() {
	l := log.New(os.Stdout, "students-api ", log.LstdFlags)
	sh := handlers.NewStudents(l)
	sm := mux.NewRouter()

	getRouter := sm.Methods(http.MethodGet).Subrouter()
	getRouter.HandleFunc("/students", sh.GetStudents)
	getRouter.HandleFunc("/students/{id:[0-9]+}", sh.GetStudentByID)

	postRouter := sm.Methods(http.MethodPost).Subrouter()
	postRouter.HandleFunc("/students", sh.AddStudent)
	postRouter.Use(sh.MiddlewareStudentValidation)

	putRouter := sm.Methods(http.MethodPut).Subrouter()
	putRouter.HandleFunc("/students/{id:[0-9]+}", sh.UpdateStudent)
	putRouter.Use(sh.MiddlewareStudentValidation)

	deleteRouter := sm.Methods(http.MethodDelete).Subrouter()
	deleteRouter.HandleFunc("/students/{id:[0-9]+}", sh.DeleteStudent)

	s := http.Server{
		Addr:         bindAddress,
		Handler:      sm,
		ErrorLog:     l,
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 10 * time.Second,
		IdleTimeout:  120 * time.Second,
	}

	go func() {
		l.Println("Starting server on port", bindAddress)
		err := s.ListenAndServe()
		if err != nil {
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
