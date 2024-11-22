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
	users       = make(map[string]*User) // 存储所有用户信息
	sessions    = make(map[string]string) // 存储会话信息，模拟登录状态
	mu          sync.Mutex               // 防止并发读写冲突
	sessionLife = time.Minute * 30        // Session 超时设置
)

// User 结构体定义用户信息
type User struct {
	Username string `json:"username"`
	Nickname string `json:"nickname"`
	Password string `json:"password"`
}

// UserInfoResponse 用于返回用户信息
type UserInfoResponse struct {
	Username string `json:"username"`
	Nickname string `json:"nickname"`
}

// helper function to get session user from cookie or token
func getSessionUser(r *http.Request) *User {
	sessionToken := r.Header.Get("Authorization")
	if sessionToken == "" {
		return nil
	}
	mu.Lock()
	defer mu.Unlock()
	username, exists := sessions[sessionToken]
	if !exists {
		return nil
	}
	user, exists := users[username]
	if !exists {
		return nil
	}
	return user
}

// Register 用户注册
func Register(w http.ResponseWriter, r *http.Request) {
	var user User
	decoder := json.NewDecoder(r.Body)
	err := decoder.Decode(&user)
	if err != nil || user.Username == "" || user.Password == "" {
		http.Error(w, "Invalid request", http.StatusBadRequest)
		return
	}

	mu.Lock()
	defer mu.Unlock()

	if _, exists := users[user.Username]; exists {
		http.Error(w, "User already exists", http.StatusConflict)
		return
	}

	// Store user
	users[user.Username] = &user
	w.WriteHeader(http.StatusCreated)
	fmt.Fprintf(w, "User %s registered successfully", user.Username)
}

// Login 用户登录，生成 session token
func Login(w http.ResponseWriter, r *http.Request) {
	var user User
	decoder := json.NewDecoder(r.Body)
	err := decoder.Decode(&user)
	if err != nil || user.Username == "" || user.Password == "" {
		http.Error(w, "Invalid request", http.StatusBadRequest)
		return
	}

	mu.Lock()
	defer mu.Unlock()

	storedUser, exists := users[user.Username]
	if !exists || storedUser.Password != user.Password {
		http.Error(w, "Invalid username or password", http.StatusUnauthorized)
		return
	}

	// Create session (simple token mechanism)
	sessionToken := fmt.Sprintf("%d", time.Now().UnixNano()) // Generate a simple session token
	sessions[sessionToken] = user.Username
	http.SetCookie(w, &http.Cookie{
		Name:    "session_token",
		Value:   sessionToken,
		Expires: time.Now().Add(sessionLife),
	})

	w.WriteHeader(http.StatusOK)
	fmt.Fprintf(w, "Logged in successfully")
}

// Logout 用户登出，删除 session
func Logout(w http.ResponseWriter, r *http.Request) {
	sessionToken := r.Header.Get("Authorization")
	if sessionToken == "" {
		http.Error(w, "Not authenticated", http.StatusUnauthorized)
		return
	}

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

// ChangePassword 修改密码
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

	user.Password = request.NewPassword
	w.WriteHeader(http.StatusOK)
	fmt.Fprintf(w, "Password changed successfully")
}

// GetUserInfo 获取用户信息
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

// UpdateUserInfo 修改用户信息
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
