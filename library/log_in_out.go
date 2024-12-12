package main

import (
	"log"
	"net/http"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v4"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

var DB *gorm.DB
var jwtKey = []byte("this_is_a_secret_key") //这是对token加密的密钥

type User struct {
	ID uint `gorm:"primaryKey"`
	//在 GORM 中，当你将 ID 字段定义为 uint 或 int 类型并且标记为 primaryKey，它会自动成为自增主键。
	Username string `json:"username" binding:"required" gorm:"unique;not null"`
	Password string `json:"password" binding:"required" `
	Role     string `json:"role" binding:"required,oneof=admin user"` // "admin" or "user"
}

type Book struct {
	ID     uint   `json:"id" binding:"required"`
	ISBN   string `json:"ISBN" binding:"required"` 
	Title  string `json:"title" binding:"required"`
	Author string `json:"author" binding:"required"`
}

type Claims struct {
	Username             string `json:"username"`
	Role                 string `json:"role"`
	jwt.RegisteredClaims        //Registered Claims（注册声明）是 JWT 标准中预定义的字段，用于提供常见信息。 //这是匿名嵌套结构体使 jwt.RegisteredClaims 的所有字段变成 Claims 的直接字段
}

//是jwt中payload的一部分，带上一些有用的信息

func initDB() {
	var err error
	dsn := "root:WEEdru1351.@tcp(123.56.118.201:3306)/library?charset=utf8mb4&parseTime=True&loc=Local"
	DB, err = gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatal("连接数据库失败:", err)
	}

	err = DB.AutoMigrate(&User{}, &Book{})
	if err != nil {
		log.Fatal("数据库迁移失败:", err)
	}
}

func login(c *gin.Context) {
	var user User
	if err := c.ShouldBindJSON(&user); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "请求参数错误"})
		return
	}

	//检测密码正确与否
	var dbUser User
	if err := DB.Where("username = ? AND password = ?", user.Username, user.Password).Take(&dbUser).Error; err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "用户名或密码错误"})
		return
	}

	expirationTime := time.Now().Add(24 * time.Hour)
	claims := &Claims{
		Username: dbUser.Username,
		Role:     dbUser.Role,
		//匿名嵌套的字段虽然可以省略嵌套结构体的名称来访问，但本质上它们仍然是属于嵌套结构体的。
		/*
			在 Go 中，结构体中的匿名嵌套字段（如 jwt.RegisteredClaims）必须显式指定字段名来初始化，即使它是匿名的。

			匿名嵌套的本质是：

			字段访问可以直接通过外层结构体访问（如 claims.ExpiresAt），
			但初始化时，仍然需要指明它是哪个嵌套字段的值。
		*/
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(expirationTime),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims) //创建新Token
	tokenString, err := token.SignedString(jwtKey)             //用密钥对Token加密，生成tokenString
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "生成Token失败"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"token": tokenString}) //构造JSON响应，返回生成的JWT Token字符串
	//例如{"token": "eyJhbGciOiJIUzI1NiIsIn..."}  客户端可以使用该 Token 在后续请求中进行身份验证

	//需要前端把这个token拿到请求头中

}

// 定义了一个 Gin 框架的中间件函数,gin.HandlerFunc 类型：Gin 的中间件类型，用于封装处理 HTTP 请求的逻辑。
func authenticateMiddleware() gin.HandlerFunc {

	//返回gin.HandlerFunc函数类型
	//签名为func(c *gin.Context)  c *gin.Context 是 Gin 框架提供的上下文对象，它保存了与 HTTP 请求相关的所有信息
	//如在中间件中，你可以使用 Context 来控制请求的流转，并将数据传递给后续的处理函数。
	//可以通过 Context 访问请求中的参数、header、URL 等信息。
	//可以通过 Context 设置 HTTP 响应的状态码、header 和 body。
	//因此，authenticateMiddleware 的返回类型实际上是一个 gin.HandlerFunc（即 func(c *gin.Context) 类型）。这是 Gin 中间件的标准定义方式

	return func(c *gin.Context) {
		tokenString := c.GetHeader("Authorization") //从请求头中拿到的tokenString
		if tokenString == "" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "未提供Token"})
			c.Abort()
			//c.Abort() 是 Gin 框架中的一个方法，表示中止当前请求的处理流程，并且后续的处理函数不会再执行。
			//调用 c.Abort() 后，Gin 会停止进一步的处理，直接返回响应,实际上只是阻止后继中间件继续运行
			//用return来彻底终结程序
			return
		}

		//解析并验证 JWT
		//jwt.ParseWithClaims：用于解析带有自定义 Claims 的 JWT。
		//tokenString 是你从请求中提取到的 JWT 字符串
		//&Claims{} 表示将解析后的 JWT 存储到一个 Claims 类型的变量中。Claims 是一个结构体类型，包含了 JWT 中的自定义字段（如用户名和角色）

		token, err := jwt.ParseWithClaims(tokenString, &Claims{}, func(token *jwt.Token) (interface{}, error) {
			return jwtKey, nil
		}) //解析

		//jwt.ParseWithClaims 会根据 JWT 字符串中的数据解析出 JWT 的 payload 和 header，并且进行签名验证。然后返回一个 token 和 err：
		//token 是解析后的 JWT 对象（包含 Claims 数据）。
		//err 是解析过程中的错误信息。
		//在 jwt.ParseWithClaims 中，JWT 的签名会根据传入的 jwtKey 进行验证
		//如果 JWT 的签名与提供的密钥匹配，且 JWT 没有过期、没有被篡改，那么 token.Valid 会被设置为 true。否则，它将被设置为 false。
		//回调函数（Callback Function）是传递给其他函数的一个函数，它将在某个特定事件发生时执行。在这里，回调函数的作用是提供签名验证的密钥.
		//这个回调函数的作用是告诉 ParseWithClaims 如何从 token 中提取签名的密钥，以便对 JWT 的签名进行验证。
		//jwt.ParseWithClaims定义中就要一个函数来返回密匙,看定义中的keyfunc

		if err != nil || !token.Valid {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "无效的Token"})
			c.Abort()
			return
		}

		//Go 中的接口类型（interface{}）可以存储任何类型的值，而不需要事先知道具体的类型。这种特性被称为动态类型，因为接口类型的具体类型在运行时才确定。
		//类型断言（Type Assertion）是 Go 中的一种机制，用于从接口类型中提取具体类型的值。
		//x, ok := i.(T)
		//i 是接口类型，T 是你希望转换成的具体类型
		//如果 i 的动态类型是 T，x 会被赋值为 i 中的值，ok 为 true
		//如果 i 的动态类型不是 T，x 会是类型 T 的零值，ok 为 false

		claims, ok := token.Claims.(*Claims)
		if !ok {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "解析Token失败"})
			c.Abort()
			return
		}

		c.Set("username", claims.Username)
		c.Set("role", claims.Role)
		//将解析出的 Username 和 Role 存入 Gin 的 Context
		//Context 提供了 Set 和 Get 方法，可以在请求的生命周期内存储和传递数据。例如，在 authenticateMiddleware 中，你可以通过 c.Set("username", claims.Username) 将 username 存入上下文中，以便后续的路由处理函数使用。
		//供后续的处理逻辑（如路由处理函数或其他中间件）使用
		c.Next()
		//继续调用下一个中间件或实际的路由处理函数
		//只有当 Token 验证通过，才会执行后续的业务逻辑
	}
}

func adminMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		role, exists := c.Get("role")
		if !exists || role != "admin" {
			c.JSON(http.StatusForbidden, gin.H{"error": "需要管理员权限"})
			c.Abort()
			return
		}
		c.Next()
	}
}

func register(c *gin.Context) {
	var user User
	if err := c.ShouldBindJSON(&user); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "请求参数错误"})
		return
	}

	if err := DB.Create(&user).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "注册失败"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "注册成功"})
}

func addBook(c *gin.Context) {
	var book Book
	if err := c.ShouldBindJSON(&book); err != nil {
		errFields := []string{}
		errMsg := err.Error()
		if strings.Contains(errMsg, "Key: 'Book.ID'") {
			errFields = append(errFields, "ID")
		}
		if strings.Contains(errMsg, "Key: 'Book.ISBN'") {
			errFields = append(errFields, "ISBN")
		}
		if strings.Contains(errMsg, "Key: 'Book.Title'") {
			errFields = append(errFields, "Title")
		}
		if strings.Contains(errMsg, "Key: 'Book.Author'") {
			errFields = append(errFields, "Author")
		}
		c.JSON(400, gin.H{
			"error":   "字段校验失败",
			"missing": errFields,
		})
		return
	}

	if book.ID == " " || book.ISBN == " " || book.Title == " " || book.Author == " " {
		c.JSON(400, gin.H{"error": "字段不能为空: ID, ISBN, Title, Author 均为必填项"})
		return
	}

	result := DB.Create(&book)
	if result.Error != nil {
		if strings.Contains(result.Error.Error(), "Duplicate entry") {
			c.JSON(400, gin.H{"error": "ID已存在，无法插入"})
		} else {
			c.JSON(500, gin.H{"error": "数据库操作失败: " + result.Error.Error()})
		}
		return
	}
	c.JSON(201, gin.H{"message": "书籍添加成功", "book": book})
}

func getBooks(c *gin.Context) {
	var books []Book
	if result := DB.Find(&books); result.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch books"})
		return
	}
	c.JSON(http.StatusOK, books)
}

func deleteBook(c *gin.Context) {
	var request struct {
		ID string `json:"id" binding:"required"`
	}

	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(400, gin.H{"error": "请求参数错误: " + err.Error()})
		return
	}

	// 删除对应的记录
	result := DB.Delete(&Book{}, request.ID)
	if result.Error != nil {
		c.JSON(500, gin.H{"error": "数据库操作失败: " + result.Error.Error()})
		//将你传递给它的数据（通常是一个 Go 的结构体、切片或其他数据类型）自动转换为 JSON 格式，然后作为 HTTP 响应发送给客户端。
		return
	}
	if result.RowsAffected == 0 {
		c.JSON(404, gin.H{"error": "未找到要删除的书籍"})
		return
	}

	c.JSON(200, gin.H{"message": "书籍删除成功", "id": request.ID})
}

func changeBook(c *gin.Context) {
	var newBook struct {
		ID     string `json:"id" binding:"required"`
		ISBN   string `json:"ISBN" binding:"required"`
		Title  string `json:"title" binding:"required"`
		Author string `json:"author" binding:"required"`
	}

	if err := c.ShouldBindJSON(&newBook); err != nil {
		c.JSON(400, gin.H{"error": "请求参数出错：" + err.Error()})
		return
	}

	result := DB.Model(&Book{}).Where(("id = ?"), newBook.ID).Updates(Book{
		ISBN:   newBook.ISBN,
		Title:  newBook.Title,
		Author: newBook.Author,
	})
	// 判断更新是否成功
	if result.Error != nil {
		c.JSON(500, gin.H{"error": "数据库更新失败：" + result.Error.Error()})
		return
	}

	// 返回更新成功的响应
	if result.RowsAffected == 0 {
		c.JSON(404, gin.H{"error": "未找到要更新的书籍"})
		return
	}

	c.JSON(200, gin.H{"message": "书籍更新成功", "updatedBook": newBook})
}

func corsMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		// 设置跨域头
		c.Writer.Header().Set("Access-Control-Allow-Origin", "*")                                              // 允许所有域名
		c.Writer.Header().Set("Access-Control-Allow-Credentials", "true")                                      // 是否允许发送Cookie
		c.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization, X-Requested-With") // 允许的请求头
		c.Writer.Header().Set("Access-Control-Allow-Methods", "GET, POST, OPTIONS, PUT, DELETE")               // 允许的请求方法

		// 对于 OPTIONS 请求，直接返回
		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(http.StatusNoContent)
			return
		}

		c.Next()
	}
}

func main() {
	initDB()
	r := gin.Default()
	r.Use(corsMiddleware())

	r.POST("/login", login)
	r.POST("/register", register)

	auth := r.Group("/").Use(authenticateMiddleware())
	{
		auth.POST("/addBook", adminMiddleware(), addBook)
		auth.GET("/getBooks", getBooks)
		auth.DELETE("/deleteBook", adminMiddleware(), deleteBook)
		auth.PUT("/changeBook", adminMiddleware(), changeBook)
	}

	if err := r.Run(":9090"); err != nil {
		log.Fatal("启动服务器失败:", err)
	}
}

//记得用这个docker run --restart always -d your_image_name
//docker容器ip似乎与服务器ip不太一样
