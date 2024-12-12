package main

import (
	"log"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

var DB *gorm.DB

type Book struct {
	ID     string `json:"id" binding:"required"`
	ISBN   string `json:"ISBN" binding:"required"`
	Title  string `json:"title" binding:"required"`
	Author string `json:"author" binding:"required"`
}

func initDB() {
	var err error
	dsn := "root:WEEdru1351.@tcp(123.56.118.201:3306)/library?charset=utf8mb4&parseTime=True&loc=Local"
	DB, err = gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatal("连接数据库失败:", err)
	}

	err = DB.AutoMigrate(&Book{})
	if err != nil {
		log.Fatal("数据库迁移失败:", err)
	}
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

	// 添加跨域中间件
	r.Use(corsMiddleware())

	r.POST("/addBook", addBook)
	r.GET("/getBooks", getBooks)
	r.DELETE("/deleteBook", deleteBook)
	r.PUT("/changeBook", changeBook)

	if err := r.Run(":9090"); err != nil {
		log.Fatal("启动服务器失败:", err)
	}
}
