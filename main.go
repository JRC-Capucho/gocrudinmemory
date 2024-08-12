package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/rs/xid"
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

var db DB

func main() {
	r := chi.NewMux()

	r.Use(middleware.RequestID)
	r.Use(middleware.RealIP)
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)

	r.Route("/api", func(r chi.Router) {
		r.Route("/users", func(r chi.Router) {
			r.Get("/", HandlerFetchUser)
			r.Post("/", HandlerRegisterUser)
			r.Get("/{id}", HandlerGetUser)
			r.Delete("/{id}", HandlerDeleteUser)
			r.Put("/{id}", HandlerUpdateUser)
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

func HandlerDeleteUser(w http.ResponseWriter, r *http.Request) {
	id := strings.TrimPrefix(r.URL.Path, "/api/users/")

	DeleteUser(id)
}

func HandlerGetUser(w http.ResponseWriter, r *http.Request) {
	id := strings.TrimPrefix(r.URL.Path, "/api/users/")

	user, err := GetUser(id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
	}

	result := fmt.Sprintf("ID: %s,First Name: %s,Last Name: %s,Biography: %s\n", user.ID, user.FirstName, user.LastName, user.Biography)

	fmt.Fprintf(w, result)
}

func HandlerFetchUser(w http.ResponseWriter, r *http.Request) {
	var result string

	for _, user := range FetchUser() {
		result += fmt.Sprintf("ID: %s,First Name: %s,Last Name: %s,Biography: %s\n", user.ID, user.FirstName, user.LastName, user.Biography)
	}

	fmt.Fprintf(w, result)
}

func HandlerUpdateUser(w http.ResponseWriter, r *http.Request) {
	id := strings.TrimPrefix(r.URL.Path, "/api/users/")

	body, err := io.ReadAll(r.Body)
	if err != nil {
		panic(err)
	}

	var user NewUser

	err = json.Unmarshal(body, &user)
	if err != nil {
		panic(err)
	}

	UpdateUser(id, user.FirstName, user.LastName, user.Biography)
}

func HandlerRegisterUser(w http.ResponseWriter, r *http.Request) {
	body, err := io.ReadAll(r.Body)
	if err != nil {
		panic(err)
	}

	var user NewUser

	err = json.Unmarshal(body, &user)
	if err != nil {
		panic(err)
	}

	RegisterUser(user.FirstName, user.LastName, user.Biography)

	w.WriteHeader(http.StatusCreated)
}

func RegisterUser(firstName, lastName, biography string) {
	user := User{
		ID:        xid.New().String(),
		FirstName: firstName,
		LastName:  lastName,
		Biography: biography,
	}
	db.Users = append(db.Users, user)
}

func FetchUser() []User {
	return db.Users
}

func GetUser(id string) (User, error) {
	for _, user := range db.Users {
		if id == user.ID {
			return user, nil
		}
	}
	return User{}, errors.New("user not found")
}

func UpdateUser(id, firstName, lastName, biography string) {
	user, err := GetUser(id)
	if err != nil {
		return
	}

	user.LastName = lastName
	user.Biography = biography
	user.FirstName = firstName
}

func DeleteUser(id string) {
	for i, user := range db.Users {
		if user.ID == id {
			db.Users = append(db.Users[:i], db.Users[i+1:]...)
			return
		}
	}
}
