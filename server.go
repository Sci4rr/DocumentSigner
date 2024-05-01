package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"golang.org/x/crypto/bcrypt"
)

type User struct {
	Username string
	Password string 
}

var Users = make(map[string]string)

func main() {
	if err := godotenv.Load(); err != nil {
		log.Fatal("Error loading .env file")
	}

	r := gin.Default()

	r.Use(gin.Logger())
	r.Use(gin.Recovery())

	api := r.Group("/api")
	{
		api.POST("/upload", UploadDocument)
		api.GET("/download/:filename", DownloadDocument)
		api.POST("/share", ShareDocument)
		api.POST("/login", Login)
		api.POST("/register", Register)
	}

	err := r.RunTLS(":https", "server.crt", "server.key") 
	if err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}

func UploadDocument(c *gin.Context) {
}

func DownloadDocument(c *gin.Context) {
	filename := c.Param("filename")
}

func ShareDocument(c *gin.Context) {
}

func Login(c *gin.Context) {
	var loginCreds User
	if err := c.BindJSON(&loginCreds); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": fmt.Sprintf("Failed to login: %v", err)})
		return
	}

	storedPassword, exists := Users[loginCreds.Username]
	if !exists || bcrypt.CompareHashAndPassword([]byte(storedPassword), []byte(loginCreds.Password)) != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid username or password"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Login successful"})
}

func Register(c *gin.Context) {
	var newUser User
	if err := c.BindJSON(&newUser); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": fmt.Sprintf("Failed to register: %v", err)})
		return
	}

	if _, exists := Users[newUser.Username]; exists {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Username already taken"})
		return
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(newUser.Password), bcrypt.DefaultCost)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to hash password"})
		return
	}

	Users[newUser.Username] = string(hashedPassword)
	c.JSON(http.StatusOK, gin.H{"message": "User registered successfully"})
}