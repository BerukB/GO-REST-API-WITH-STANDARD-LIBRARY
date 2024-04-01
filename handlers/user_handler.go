package handlers

import (
	"encoding/json"
	"log"
	"math/rand"
	"net/http"
	"regexp"
	"strconv"

	middleware "github.com/BerukB/GO-REST-API-WITH-STANDARD-LIBRARY/middleware"
	usermodel "github.com/BerukB/GO-REST-API-WITH-STANDARD-LIBRARY/models"
	"golang.org/x/crypto/bcrypt"
)

func hashPassword(password string) (string, error) {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}
	return string(hashedPassword), nil
}

type userStore interface {
	Add(user usermodel.User) error
	Get(id string) (usermodel.User, error)
	GetEmail(email string) (usermodel.User, error)
	Update(id string, user usermodel.User) (usermodel.User, error)
	List() ([]usermodel.User, error)
	Remove(id string) error
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
	// Hash the password before storing it
	hashedPassword, err := hashPassword(user.PassWord)
	if err != nil {
		InternalServerErrorHandler(w, r)
		return
	}
	user.PassWord = hashedPassword

	if err := h.store.Add(user); err != nil {
		InternalServerErrorHandler(w, r)
		return
	}
	// Set the status code to 200
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(user)
}

func (h *UserHandler) ListUsers(w http.ResponseWriter, r *http.Request) {
	// Apply middleware before executing ListUsers handler logic
	middleware.JWTAuthMiddleware(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		resources, err := h.store.List()
		if err != nil {
			InternalServerErrorHandler(w, r)
			return
		}
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(resources)
	})).ServeHTTP(w, r)
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
	w.WriteHeader(http.StatusOK)
	// Encode the user object to JSON and write it to the response.
	err = json.NewEncoder(w).Encode(user)
	if err != nil {
		// Log the error and return an internal server error response.
		log.Printf("Error encoding user to JSON: %v", err)
		InternalServerErrorHandler(w, r)
		return
	}

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
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(user)

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
