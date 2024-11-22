package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"sync"
	"time"
)

var (
	mu         sync.Mutex
	Users      = make(map[string]*User)
	sessions   = make(map[string]string)
	sessionLife = time.Minute * 30
)

type User struct {
	Username string `json:"username"`
	Nickname string `json:"nickname"`
	Password string `json:"password"`
}

type UserInfoResponse struct {
	Username string `json:"username"`
	Nickname string `json:"nickname"`
}

// 检测用户名是否被占用
func Register(w http.ResponseWriter, r *http.Request) {
	var user User

	decoder := json.NewDecoder(r.Body)
	err := decoder.Decode(&user)
	if err != nil || user.Username == "" || user.Password == "" {
		http.Error(w, "Invalid Request", http.StatusBadRequest)
		return
	}

	mu.Lock()
	defer mu.Unlock()

	// 检查用户名是否已存在
	if _, exists := Users[user.Username]; exists {
		http.Error(w, "Username already exists", http.StatusBadRequest)
		return
	}

	Users[user.Username] = &user
	w.WriteHeader(http.StatusCreated)
	fmt.Fprintf(w, "User %s has successfully created!", user.Username)
}

func Login(w http.ResponseWriter, r *http.Request) {
	var user User
	
	decoder := json.NewDecoder(r.Body)
	err := decoder.Decode(&user)
	if err != nil || user.Username == "" || user.Password == "" {
		http.Error(w, "Invalid Request", http.StatusBadRequest)
		return
	}

	mu.Lock()
	defer mu.Unlock()

	storedUser, exists := Users[user.Username]
	if !exists || user.Password != storedUser.Password {
		http.Error(w, "Invalid Account or Password", http.StatusUnauthorized)
		return
	}

	sessionToken := fmt.Sprintf("%d", time.Now().UnixNano())
	sessions[sessionToken] = user.Username
	http.SetCookie(w, &http.Cookie{
		Name:    "session_token",
		Value:   sessionToken,
		Expires: time.Now().Add(sessionLife),
	})
	
	w.WriteHeader(http.StatusOK)
	fmt.Fprintf(w, "Welcome %s", storedUser.Nickname)
}

func Logout(w http.ResponseWriter, r *http.Request) {
	// 从请求的 Cookie 中获取 session_token
	sessionCookie, err := r.Cookie("session_token")
	if err != nil || sessionCookie == nil {
		http.Error(w, "Not Authorized", http.StatusUnauthorized)
		return
	}

	sessionToken := sessionCookie.Value

	mu.Lock()
	defer mu.Unlock()

	// 删除 sessionToken
	delete(sessions, sessionToken)

	// 清空 Cookie
	http.SetCookie(w, &http.Cookie{
		Name:    "session_token",
		Value:   "",
		Expires: time.Now().Add(-time.Minute),
	})

	w.WriteHeader(http.StatusOK)
	fmt.Fprintf(w, "Logged out successfully")
}

// 根据当前请求的 sessionToken 获取用户
func getSessionUser(r *http.Request) *User {
	sessionCookie, err := r.Cookie("session_token") // 从 Cookie 中获取 session_token
	if err != nil || sessionCookie == nil {
		return nil
	}

	sessionToken := sessionCookie.Value

	mu.Lock()
	defer mu.Unlock()

	username, exists := sessions[sessionToken] // 根据 sessionToken 获取用户名
	if !exists {
		return nil
	}

	user, exists := Users[username] // 根据用户名获取用户
	if !exists {
		return nil
	}

	return user
}

func ChangePassword(w http.ResponseWriter, r *http.Request) {
	user := getSessionUser(r)
	if user == nil {
		http.Error(w, "Not authenticated", http.StatusUnauthorized)
		return
	}

	var request struct {
		OldPassword string `json:"old_password"`
		NewPassword string `json:"new_password"`
	}

	decoder := json.NewDecoder(r.Body)
	err := decoder.Decode(&request)
	if err != nil || request.OldPassword == "" || request.NewPassword == "" {
		http.Error(w, "Invalid request", http.StatusBadRequest)
		return
	}

	if user.Password != request.OldPassword {
		http.Error(w, "Old password incorrect", http.StatusUnauthorized)
		return
	}

	mu.Lock()
	defer mu.Unlock()

	user.Password = request.NewPassword
	w.WriteHeader(http.StatusOK)
	fmt.Fprintf(w, "Password changed successfully")
}

func GetUserInfo(w http.ResponseWriter, r *http.Request) {
	user := getSessionUser(r)
	if user == nil {
		http.Error(w, "Not authenticated", http.StatusUnauthorized)
		return
	}

	response := UserInfoResponse{
		Username: user.Username,
		Nickname: user.Nickname,
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

func UpdateUserInfo(w http.ResponseWriter, r *http.Request) {
	user := getSessionUser(r)
	if user == nil {
		http.Error(w, "Not authenticated", http.StatusUnauthorized)
		return
	}

	var request struct {
		Nickname string `json:"nickname"`
	}
	decoder := json.NewDecoder(r.Body)
	err := decoder.Decode(&request)
	if err != nil || request.Nickname == "" {
		http.Error(w, "Invalid request", http.StatusBadRequest)
		return
	}

	user.Nickname = request.Nickname
	w.WriteHeader(http.StatusOK)
	fmt.Fprintf(w, "Nickname updated successfully")
}

func main() {
	http.HandleFunc("/register", Register)
	http.HandleFunc("/login", Login)
	http.HandleFunc("/logout", Logout)
	http.HandleFunc("/change-password", ChangePassword)
	http.HandleFunc("/user-info", GetUserInfo)
	http.HandleFunc("/update-info", UpdateUserInfo)

	log.Println("Starting server on :8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
