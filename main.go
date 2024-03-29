package main

import (
	"net/http"
	"regexp"
)

type homeHandler struct{}

func (h *homeHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("This is my home page"))
}

type UserHandler struct{}

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
func (h *UserHandler) ListUsers(w http.ResponseWriter, r *http.Request)  {}
func (h *UserHandler) GetUser(w http.ResponseWriter, r *http.Request)    {}
func (h *UserHandler) UpdateUser(w http.ResponseWriter, r *http.Request) {}
func (h *UserHandler) DeleteUser(w http.ResponseWriter, r *http.Request) {}
func (h *UserHandler) CreateUser(w http.ResponseWriter, r *http.Request) {}

var (
	UserRe       = regexp.MustCompile(`^/user/*$`)
	UserReWithID = regexp.MustCompile(`^/user/([a-z0-9]+(?:-[a-z0-9]+)+)$`)
)

func main() {

	// Create a new request multiplexer
	// Take incoming requests and dispatch them to the matching handlers
	mux := http.NewServeMux()

	// Register the routes and handlers
	mux.Handle("/", &homeHandler{})
	mux.Handle("/user", &UserHandler{})
	mux.Handle("/user/", &UserHandler{})

	// Run the server
	http.ListenAndServe(":8080", mux)
}
