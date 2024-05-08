package main

import (
	"encoding/json"
	"net/http"
	"os"
)

type User struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

func main() {
	http.HandleFunc("/sign-up", handleSignUp)
	http.ListenAndServe(":8080", nil)
}

func handleSignUp(w http.ResponseWriter, r *http.Request) {

	email := r.URL.Query().Get("email")
	password := r.URL.Query().Get("password")

	user := User{
		Email:    email,
		Password: password,
	}

	file, err := os.OpenFile("users.json", os.O_RDWR|os.O_CREATE, 0644)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer file.Close()

	var users []User
	if err := json.NewDecoder(file).Decode(&users); err != nil && err.Error() != "EOF" {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	users = append(users, user)

	if _, err := file.Seek(0, 0); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if err := json.NewEncoder(file).Encode(users); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")

	if err := json.NewEncoder(w).Encode(user); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}
