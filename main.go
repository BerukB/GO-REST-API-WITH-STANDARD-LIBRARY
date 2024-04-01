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

	initialUsers := []usermodel.User{
		{ID: "48816866", UserName: "Alice", Email: "alice@gmail.com", PassWord: "$2a$10$EoNuChNRUnvoQVR1p.oucegJgZ.oQ6NMn/uO7SuBvcTLUVyuDb9cq"},
	}

	store := usermodel.NewMemStore(initialUsers)

	userHandler := handler.NewUserHandler(store)

	mux := http.NewServeMux()

	mux.HandleFunc("/login", func(w http.ResponseWriter, r *http.Request) {
		handler.Login(store, w, r)
	})

	// mux.HandleFunc("/login/", handler.Login)

	mux.Handle("/", &homeHandler{})
	mux.Handle("/user", userHandler)
	mux.Handle("/user/", userHandler)

	http.ListenAndServe(":8080", mux)
}
