package main

import (
	"net/http"

	handler "github.com/BerukB/GO-REST-API-WITH-STANDARD-LIBRARY/handlers"
	usermodel "github.com/BerukB/GO-REST-API-WITH-STANDARD-LIBRARY/models"
)

type homeHandler struct{}

func (h *homeHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("This is my home page"))
}
func main() {

	store := usermodel.NewMemStore()
	userHandler := handler.NewUserHandler(store)

	mux := http.NewServeMux()

	mux.Handle("/", &homeHandler{})
	mux.Handle("/user", userHandler)
	mux.Handle("/user/", userHandler)

	http.ListenAndServe(":8080", mux)
}
