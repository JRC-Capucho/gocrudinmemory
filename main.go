package main

import (
	"net/http"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

type User struct {
	ID        string
	FirstName string
	LastName  string
	Biography string
}

type NewUser struct {
	FirstName string
	LastName  string
	Biography string
}

type Response struct {
	Error string `json:"error,omitempty"`
	Data  any    `json:"data,omitempty"`
}

type DB struct {
	Users []User
}

func main() {
	r := chi.NewMux()

	r.Use(middleware.RequestID)
	r.Use(middleware.RealIP)
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)

	r.Route("/api", func(r chi.Router) {
		r.Route("/users", func(r chi.Router) {
			r.Get("/", func(w http.ResponseWriter, r *http.Request) {})
			r.Post("/", HandlerRegisterUser)
		})
	})

	s := http.Server{
		Addr:         ":3333",
		IdleTimeout:  10 * time.Second,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 10 * time.Second,
		Handler:      r,
	}

	s.ListenAndServe()
}

func ResponseApi(w http.ResponseWriter, res Response, status int) {
	w.WriteHeader(status)
}

func HandlerRegisterUser(w http.ResponseWriter, r *http.Request) {
}

func (db *DB) RegisterUser(firstName, lastName, biography string) {
	user := User{}
	db.Users = append(db.Users, user)
}
