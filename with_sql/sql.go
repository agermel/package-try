package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"sync"
	"time"

	_ "github.com/go-sql-driver/mysql"
)

var (
	db          *sql.DB //数据库的连接池
	mu          sync.Mutex
	sessionLife = time.Minute * 30
)

// 用户信息
type User struct {
	Username string `json:"username"`
	Nickname string `json:"nickname"`
	Password string `json:"password"`
}

// 返回用户信息
type UserInfoResponse struct {
	Username string `json:"username"`
	Nickname string `json:"nickname"`
}

func init() {
	var err error

	// 获取环境变量配置 MySQL 连接信息
	dbHost := os.Getenv("DB_HOST")
	dbPort := os.Getenv("DB_PORT")
	dbUser := os.Getenv("DB_USER")
	dbPassword := os.Getenv("DB_PASSWORD")
	dbName := os.Getenv("DB_NAME")

	// 构建数据库连接字符串
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/", dbUser, dbPassword, dbHost, dbPort)

	// 连接到 MySQL
	db, err = sql.Open("mysql", dsn)
	if err != nil {
		log.Fatal("无法连接到数据库: ", err)
	}

	// 测试数据库连接
	if err := db.Ping(); err != nil {
		log.Fatal("无法 ping 到数据库: ", err)
	}

	// 创建数据库,若无指定数据库则创建新的
	_, err = db.Exec(fmt.Sprintf("CREATE DATABASE IF NOT EXISTS %s", dbName))
	if err != nil {
		log.Fatal("无法创建数据库: ", err)
	}

	// 连接到指定数据库
	dsn = fmt.Sprintf("%s:%s@tcp(%s:%s)/%s", dbUser, dbPassword, dbHost, dbPort, dbName)
	db, err = sql.Open("mysql", dsn)
	if err != nil {
		log.Fatal("没连接到数据库: ", err)
	}

	// 测试连接
	if err := db.Ping(); err != nil {
		log.Fatal("没 ping 到数据库: ", err)
	}

	// 若没有表直接创建
	createTables := []string{
		`CREATE TABLE IF NOT EXISTS users (
			id INT AUTO_INCREMENT PRIMARY KEY,
			username VARCHAR(255) NOT NULL UNIQUE,
			nickname VARCHAR(255) DEFAULT '',
			password VARCHAR(255) NOT NULL
		)`,
		`CREATE TABLE IF NOT EXISTS sessions (
			id INT AUTO_INCREMENT PRIMARY KEY,
			session_token VARCHAR(255) NOT NULL UNIQUE,
			username VARCHAR(255) NOT NULL,
			expire_at DATETIME NOT NULL,
			FOREIGN KEY (username) REFERENCES users(username) ON DELETE CASCADE
		)`,
	}

	for _, query := range createTables {
		_, err = db.Exec(query)
		if err != nil {
			log.Fatal("无法创建表: ", err)
		}
	}

	log.Println("初始化完成")
}

func getSessionUser(r *http.Request) *User {
	cookie, err := r.Cookie("session_token")
	if err != nil {
		return nil
	}
	sessionToken := cookie.Value

	mu.Lock()
	defer mu.Unlock()

	var username string
	var expiresAT time.Time
	err = db.QueryRow("SELECT username, expire_at FROM sessions WHERE session_token = ?", sessionToken).Scan(&username, &expiresAT)
	if err != nil || expiresAT.Before(time.Now()) {
		return nil
	}

	var user User
	err = db.QueryRow("SELECT username, nickname, password FROM users WHERE username = ?", username).Scan(&user.Username, &user.Nickname, &user.Password)
	if err != nil {
		return nil
	}

	return &user
}

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

	var exist int
	err = db.QueryRow("SELECT COUNT(*) FROM users WHERE username = ?", user.Username).Scan(&exist)

	if exist > 0 {
		http.Error(w, "你注册过了喵", http.StatusConflict)
		return
	}

	_, err = db.Exec("INSERT INTO users (username,nickname,password) VALUE (?,?,?)", user.Username, user.Nickname, user.Password)
	if err != nil {
		http.Error(w, "注册不了喵", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
	fmt.Fprintf(w, "%s 你已经注册好了喵", user.Username)
}

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

	var storedUser User

	// 查询登录输入的username是否存在而且查询给username对应的password是否与输入的一致
	err = db.QueryRow("SELECT username , password FROM users WHERE username = ?", user.Username).Scan(&storedUser.Username, &storedUser.Password)
	if err != nil || storedUser.Password != user.Password {
		http.Error(w, "Invalid username or password", http.StatusUnauthorized)
		return
	}

	sessionToken := fmt.Sprintf("%d", time.Now().UnixNano())
	expiresAT := time.Now().Add(sessionLife)

	_, err = db.Exec("INSERT INTO sessions (session_token,username,expire_at) VALUES (?,?,?)", sessionToken, user.Username, expiresAT)
	if err != nil {
		http.Error(w, "Failed to create session", http.StatusInternalServerError)
		return
	}

	http.SetCookie(w, &http.Cookie{
		Name:    "session_token",
		Value:   sessionToken,
		Expires: expiresAT,
	})

	w.WriteHeader(http.StatusOK)
	fmt.Fprintf(w, "登进去了喵")
}

func Logout(w http.ResponseWriter, r *http.Request) {
	cookie, err := r.Cookie("session_token")
	if err != nil {
		http.Error(w, "Not authenticated", http.StatusUnauthorized)
		return
	}

	sessionToken := cookie.Value

	mu.Lock()
	defer mu.Unlock()

	_, err = db.Exec("DELETE FROM sessions WHERE session_token = ?", sessionToken)
	if err != nil {
		http.Error(w, "Failed to logout", http.StatusInternalServerError)
		return
	}

	http.SetCookie(w, &http.Cookie{
		Name:    "session_token",
		Value:   "",
		Expires: time.Now().Add(-time.Minute),
	})

	w.WriteHeader(http.StatusOK)
	fmt.Fprintf(w, "出去了喵")
}

// 修改密码
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

// 获取用户信息
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

// 修改用户信息
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
	defer db.Close()

	http.HandleFunc("/register", Register)
	http.HandleFunc("/login", Login)
	http.HandleFunc("/logout", Logout)
	http.HandleFunc("/change-password", ChangePassword)
	http.HandleFunc("/user-info", GetUserInfo)
	http.HandleFunc("/update-info", UpdateUserInfo)

	log.Println("Starting server on :9000")
	log.Fatal(http.ListenAndServe(":9000", nil))
}
