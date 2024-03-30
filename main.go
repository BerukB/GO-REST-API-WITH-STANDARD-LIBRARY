package main

import (
	"encoding/json"
	"log"
	"math/rand"
	"net/http"
	"regexp"
	"strconv"

	usermodel "github.com/BerukB/GO-REST-API-WITH-STANDARD-LIBRARY/models"
)

type userStore interface {
	Add(user usermodel.User) error
	Get(id string) (usermodel.User, error)
	Update(id string, user usermodel.User) (usermodel.User, error)
	List() ([]usermodel.User, error)
	Remove(id string) error
}

type homeHandler struct{}

func (h *homeHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("This is my home page"))
}

type UserHandler struct {
	store userStore
}

func NewUserHandler(s userStore) *UserHandler {
	return &UserHandler{
		store: s,
	}
}

func (h *UserHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	switch {
	case r.Method == http.MethodPost && UserRe.MatchString(r.URL.Path):
		h.CreateUser(w, r)
		return
	case r.Method == http.MethodGet && UserRe.MatchString(r.URL.Path):
		h.ListUsers(w, r)
		return
	case r.Method == http.MethodGet && UserReWithID.MatchString(r.URL.Path):
		h.GetUser(w, r)
		return
	case r.Method == http.MethodPut && UserReWithID.MatchString(r.URL.Path):
		h.UpdateUser(w, r)
		return
	case r.Method == http.MethodDelete && UserReWithID.MatchString(r.URL.Path):
		h.DeleteUser(w, r)
		return
	default:
		return
	}
}
func InternalServerErrorHandler(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusInternalServerError)
	w.Write([]byte("500 Internal Server Error"))
}

func NotFoundHandler(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusNotFound)
	w.Write([]byte("404 Not Found"))
}

func (h *UserHandler) CreateUser(w http.ResponseWriter, r *http.Request) {
	var user usermodel.User
	if err := json.NewDecoder(r.Body).Decode(&user); err != nil {
		InternalServerErrorHandler(w, r)
		return
	}
	user.ID = strconv.Itoa(rand.Intn(100000000))
	if err := h.store.Add(user); err != nil {
		InternalServerErrorHandler(w, r)
		return
	}

	json.NewEncoder(w).Encode(user)

	// Set the status code to 200
	w.WriteHeader(http.StatusOK)
}

func (h *UserHandler) ListUsers(w http.ResponseWriter, r *http.Request) {
	resources, err := h.store.List()

	// jsonBytes, err := json.Marshal(resources)
	if err != nil {
		InternalServerErrorHandler(w, r)
		return
	}
	json.NewEncoder(w).Encode(resources)

	w.WriteHeader(http.StatusOK)
	// w.Write(jsonBytes)
}
func (h *UserHandler) GetUser(w http.ResponseWriter, r *http.Request) {
	// Extract the resource ID/slug using a regex
	matches := UserReWithID.FindStringSubmatch(r.URL.Path)
	// Expect matches to be length >= 2 (full string + 1 matching group)
	if len(matches) < 2 {
		InternalServerErrorHandler(w, r)
		return
	}

	// Retrieve user from the store
	user, err := h.store.Get(matches[1])
	if err != nil {
		// Special case of NotFound Error
		if err == usermodel.ErrNotFound {
			NotFoundHandler(w, r)
			return
		}

		// Every other error
		InternalServerErrorHandler(w, r)
		return
	}

	// Set the "Content-Type: application/json" header on the response.
	w.Header().Set("Content-Type", "application/json")

	// Encode the user object to JSON and write it to the response.
	err = json.NewEncoder(w).Encode(user)
	if err != nil {
		// Log the error and return an internal server error response.
		log.Printf("Error encoding user to JSON: %v", err)
		InternalServerErrorHandler(w, r)
		return
	}

	w.WriteHeader(http.StatusOK)
}

func (h *UserHandler) UpdateUser(w http.ResponseWriter, r *http.Request) {
	matches := UserReWithID.FindStringSubmatch(r.URL.Path)
	if len(matches) < 2 {
		InternalServerErrorHandler(w, r)
		return
	}

	// User object that will be populated from JSON payload
	var user usermodel.User
	if err := json.NewDecoder(r.Body).Decode(&user); err != nil {
		InternalServerErrorHandler(w, r)
		return

	}
	user, err := h.store.Update(matches[1], user)
	if err != nil {
		if err == usermodel.ErrNotFound {
			NotFoundHandler(w, r)
			return
		}
		InternalServerErrorHandler(w, r)
		return
	}
	json.NewEncoder(w).Encode(user)
	w.WriteHeader(http.StatusOK)
}
func (h *UserHandler) DeleteUser(w http.ResponseWriter, r *http.Request) {
	matches := UserReWithID.FindStringSubmatch(r.URL.Path)
	if len(matches) < 2 {
		InternalServerErrorHandler(w, r)
		return
	}

	if err := h.store.Remove(matches[1]); err != nil {
		InternalServerErrorHandler(w, r)
		return
	}

	w.WriteHeader(http.StatusOK)
}

var (
	UserRe       = regexp.MustCompile(`^/user/*$`)
	UserReWithID = regexp.MustCompile(`^/user/([0-9]+)$`)
)

func main() {

	store := usermodel.NewMemStore()
	userHandler := NewUserHandler(store)

	// Create a new request multiplexer
	// Take incoming requests and dispatch them to the matching handlers
	mux := http.NewServeMux()

	// Register the routes and handlers
	mux.Handle("/", &homeHandler{})
	mux.Handle("/user", userHandler)
	mux.Handle("/user/", userHandler)

	// Run the server
	http.ListenAndServe(":8080", mux)
}
