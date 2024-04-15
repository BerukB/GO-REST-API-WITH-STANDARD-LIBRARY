package handlers

import (
	"encoding/json"
	"log"
	"math/rand"
	"net/http"
	"regexp"
	"strconv"

	usermodel "github.com/BerukB/GO-REST-API-WITH-STANDARD-LIBRARY/models"
	"github.com/BerukB/GO-REST-API-WITH-STANDARD-LIBRARY/validation"
	"golang.org/x/crypto/bcrypt"
)

func hashPassword(password string) (string, error) {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}
	return string(hashedPassword), nil
}

type Response struct {
	MetaData struct {
		Page    int `json:"page"`
		PerPage int `json:"per_page"`
	} `json:"meta_data"`
	Data []usermodel.User `json:"data"`
}

type userStore interface {
	Add(user usermodel.User) error
	Get(id string) (usermodel.User, error)
	GetEmail(email string) (usermodel.User, error)
	Update(id string, user usermodel.User) (usermodel.User, error)
	List(page, limit int) ([]usermodel.User, error)
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
	if err := validation.ValidateEmail(user.Email); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	if err := validation.ValidatePhone(user.Phone); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	if err := validation.ValidateAddress(user.Address); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	_, err := h.store.GetEmail(user.Email)
	if err == nil {
		http.Error(w, "User already exists ", http.StatusConflict)
		return
	}

	user.ID = strconv.Itoa(rand.Intn(100000000))
	formattedPhone := user.Phone.Format()
	user.Phone = usermodel.PhoneNumber(formattedPhone)
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
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(user)
}

func (h *UserHandler) ListUsers(w http.ResponseWriter, r *http.Request) {
	// time.Sleep(5 * time.Second)
	page := 1
	limit := 10

	// Read pagination parameters from the query string
	pageQuery := r.URL.Query().Get("page")
	limitQuery := r.URL.Query().Get("limit")

	if pageQuery != "" {
		var err error
		page, err = strconv.Atoi(pageQuery)
		if err != nil {
			http.Error(w, "Invalid page parameter", http.StatusBadRequest)
			return
		}
	}
	if limitQuery != "" {
		var err error
		limit, err = strconv.Atoi(limitQuery)
		if err != nil {
			http.Error(w, "Invalid limit parameter", http.StatusBadRequest)
			return
		}
	}

	// Fetch users with pagination
	resources, err := h.store.List(page, limit)
	if err != nil {
		InternalServerErrorHandler(w, r)
		return
	}
	response := Response{
		MetaData: struct {
			Page    int `json:"page"`
			PerPage int `json:"per_page"`
		}{
			Page:    page,
			PerPage: limit,
		},
		Data: resources,
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(response)
}
func (h *UserHandler) GetUser(w http.ResponseWriter, r *http.Request) {
	matches := UserReWithID.FindStringSubmatch(r.URL.Path)
	if len(matches) < 2 {
		InternalServerErrorHandler(w, r)
		return
	}

	user, err := h.store.Get(matches[1])
	if err != nil {
		if err == usermodel.ErrNotFound {
			NotFoundHandler(w, r)
			return
		}

		InternalServerErrorHandler(w, r)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	err = json.NewEncoder(w).Encode(user)
	if err != nil {
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
