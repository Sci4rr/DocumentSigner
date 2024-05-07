package main

import (
    "fmt"
    "github.com/dgrijalva/jwt-go"
    "github.com/gin-gonic/gin"
    "golang.org/x/crypto/bcrypt"
    "net/http"
    "time"
)

type Claims struct {
    Username string `json:"username"`
    jwt.StandardClaims
}

var jwtKey = []byte("my_secret_key")

func main() {
}

func Login(c *gin.Context) {
    var loginCreds User
    if err := c.BindJSON(&loginCreds); err != nil {
        respondWithError(c, http.StatusBadRequest, fmt.Sprintf("Failed to login: %v", err))
        return
    }

    if storedPassword, exists := Users[loginCreds.Username]; !exists || bcrypt.CompareHashAndPassword([]byte(storedPassword), []byte(loginCreds.Password)) != nil {
        respondWithError(c, http.StatusUnauthorized, "Invalid username or password")
        return
    }

    tokenString, err := generateToken(loginCreds.Username)
    if err != nil {
        respondWithError(c, http.StatusInternalServerError, "Failed to create token")
        return
    }

    c.JSON(http.StatusOK, gin.H{"token": tokenString, "message": "Login successful"})
}

func respondWithError(c *gin.Context, code int, message string) {
    c.JSON(code, gin.H{"error": message})
}

func generateToken(username string) (string, error) {
    expirationTime := time.Now().Add(5 * time.Minute)
    claims := &Claims{
        Username: username,
        StandardClaims: jwt.StandardClaims{
            ExpiresAt: expirationTime.Unix(),
        },
    }

    token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
    return token.SignedString(jwtKey)
}