package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"sync"
	"time"
)

type User struct {
	Username string `json:"username"`
	Nickname string `json:"nickname"`
	Password string `json:"password"`
}

var (
	users    = make(map[string]*User)
	sessions = make(map[string]string)
	mu       sync.Mutex
)

func Register(w http.ResponseWriter, r *http.Request) {
	var user *User
	decoder := json.NewDecoder(r.Body)
	err := decoder.Decode(&user)
	if err != nil || user.Username == "" || user.Password == "" {
		fmt.Fprintf(w, "Invaild requset")
		return
	}

	if _, exists := users[user.Username]; exists {
		http.Error(w, "Username is already used!", http.StatusBadRequest)
		return
	}

	mu.Lock()
	defer mu.Unlock()

	users[user.Username] = user
	w.WriteHeader(http.StatusCreated)
	fmt.Fprintln(w, "Register succeeded!")
}

func Login(w http.ResponseWriter, r *http.Request) {
	var user *User
	decoder := json.NewDecoder(r.Body)
	err := decoder.Decode(&user)
	if err != nil || users[user.Username] == nil || user.Password == "" {
		fmt.Fprintln(w, "invaild request!")
	}

	mu.Lock()
	defer mu.Unlock()

	user = users[user.Username]

	session_token := fmt.Sprintf("%d", time.Now().UnixNano())
	sessions[session_token] = user.Username
	http.SetCookie(w, &http.Cookie{
		Name:    "session_token",
		Value:   session_token,
		Expires: time.Now().Add(time.Minute * 30),
	})
	w.WriteHeader(http.StatusOK)
	fmt.Fprintf(w, "%s is successfully logged in", user.Username)
}

func logout(w http.ResponseWriter, r *http.Request) {
	cookie, err := r.Cookie("session_token")
	if err != nil {
		http.Error(w, "Not authenticated", http.StatusUnauthorized)
		return
	}

	sessionToken := cookie.Value
	mu.Lock()
	defer mu.Unlock()

	delete(sessions, sessionToken)
	http.SetCookie(w, &http.Cookie{
		Name:    "session_token",
		Value:   "",
		Expires: time.Now().Add(-time.Minute),
	})

	w.WriteHeader(http.StatusOK)
	fmt.Fprintf(w, "Logged out successfully")
}

func main() {
	http.HandleFunc("/register", Register)
	http.HandleFunc("/login", Login)
	http.HandleFunc("/logout", logout)

	log.Println("Starting server on :9090")
	log.Fatal(http.ListenAndServe(":9090", nil))
}
