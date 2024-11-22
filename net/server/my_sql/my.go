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
	mu sync.Mutex
	Users = make(map[string]*User)
	sessions = make(map[string]string)
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

//似乎缺失检测用户名被占用的功能
func Register(w http.ResponseWriter,r *http.Request) {
	var user User

	decoder := json.NewDecoder(r.Body)
	err := decoder.Decode(&user)
	if err != nil || user.Username == "" ||user.Password == "" {
		http.Error(w,"Invaild Request",http.StatusBadRequest)
		return
	}

	mu.Lock()
	defer mu.Unlock()

	Users[user.Username] = &user
	w.WriteHeader(http.StatusCreated)
	fmt.Fprintf(w,"User %s has successfully created!",user.Username)
}

func Login(w http.ResponseWriter,r *http.Request) {
	var user User
	
	decoder := json.NewDecoder(r.Body)
	err := decoder.Decode(&user)
	if err!=nil || user.Username == "" || user.Password == "" {
		http.Error(w,"Invaild Request",http.StatusBadRequest)
	}

	mu.Lock()
	defer mu.Unlock()

	storedUser ,exists := Users[user.Username]
	if !exists || user.Password != storedUser.Password{
		http.Error(w,"Invaild Account or Password",http.StatusUnauthorized)
		return
	}

	sessionToken := fmt.Sprintf("%d",time.Now().UnixNano())
	sessions[sessionToken] = user.Username
	http.SetCookie(w,&http.Cookie{
		Name: "session_token",
		Value: sessionToken,
		Expires: time.Now().Add(sessionLife),
	})

	w.Header().Set("Authorization", sessionToken)
	w.WriteHeader(http.StatusOK)
	fmt.Fprintf(w,"Welcome %s",storedUser.Nickname)
}

func Logout(w http.ResponseWriter,r *http.Request) {
	sessionToken := r.Header.Get("Authorization")

	if sessionToken == "" {
		http.Error(w,"Not Authorized",http.StatusUnauthorized)
		return
	}

	mu.Lock()
	defer mu.Unlock()

	delete(sessions,sessionToken)
	http.SetCookie(w,&http.Cookie{
		Name: "session_token",
		Value: "",
		Expires: time.Now().Add(-time.Minute),
	})

	w.WriteHeader(http.StatusOK)
	fmt.Println(w,"Logged out successfully")
}

//寻找在当前request请求的session下的user
func getSessionUser (r *http.Request) *User{
	sessionToken := r.Header.Get("Authorization") //找sessiontoken
	if sessionToken == "" {
		return nil
	}

	mu.Lock()
	defer mu.Unlock()

	username,exists := sessions[sessionToken]  //通过sessiontoken找username
	if !exists {
		return nil
	}

	user,exists := Users[username] //通过username找user
	if !exists {
		return nil
	}

	return user
}

func ChangePassword (w http.ResponseWriter,r *http.Request) {
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
	//密码乱输
	if err != nil || request.OldPassword == "" || request.NewPassword == "" {
		http.Error(w, "Invalid request", http.StatusBadRequest)
		return
	}
	//密码跟老密码不一致
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

func GetUserInfo (w http.ResponseWriter,r *http.Request) {
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
	json.NewEncoder(w).Encode(response)//有点神经
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
	decoder := json.NewDecoder(r.Body) //从请求主体出找出要修改的值并解析到结构体中从而引入要改变的值
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
