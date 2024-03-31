package handlers

import (
	"encoding/json"
	"net/http"
	"time"

	usermodel "github.com/BerukB/GO-REST-API-WITH-STANDARD-LIBRARY/models"
	"github.com/dgrijalva/jwt-go"
	"golang.org/x/crypto/bcrypt"
)

var jwtSecret = []byte("12345678")

const (
	accessTokenDuration  = time.Minute * 15
	refreshTokenDuration = time.Hour * 24
)

type LoginResponse struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
}

func Login(w http.ResponseWriter, r *http.Request) {
	var user usermodel.User
	if err := json.NewDecoder(r.Body).Decode(&user); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	store := usermodel.NewMemStore()

	foundUser, err := store.Get(user.ID)
	if err != nil {
		http.Error(w, "Invalid credentials", http.StatusUnauthorized)
		return
	}
	hashedPassword := foundUser.PassWord

	err = bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(user.PassWord))
	if err != nil {
		// Passwords don't match
		http.Error(w, "Invalid credentials", http.StatusUnauthorized)
	}

	accessToken, err := generateAccessToken(foundUser.ID)
	if err != nil {
		http.Error(w, "Failed to generate access token", http.StatusInternalServerError)
		return
	}

	refreshToken, err := generateRefreshToken(foundUser.ID)
	if err != nil {
		http.Error(w, "Failed to generate refresh token", http.StatusInternalServerError)
		return
	}

	response := LoginResponse{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
	}
	jsonResponse, err := json.Marshal(response)
	if err != nil {
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(jsonResponse)
}

func generateAccessToken(userID string) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"userID": userID,
		"exp":    time.Now().Add(accessTokenDuration).Unix(),
	})
	return token.SignedString(jwtSecret)
}

func generateRefreshToken(userID string) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"userID": userID,
		"exp":    time.Now().Add(refreshTokenDuration).Unix(),
	})
	return token.SignedString(jwtSecret)
}
