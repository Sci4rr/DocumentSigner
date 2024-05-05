package main

import (
    "github.com/dgrijalva/jwt-go"
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
        c.JSON(http.StatusBadRequest, gin.H{"error": fmt.Sprintf("Failed to login: %v", err)})
        return
    }

    storedPassword, exists := Users[loginCreds.Username]
    if !exists || bcrypt.CompareHashAndPassword([]byte(storedPassword), []byte(loginCreds.Password)) != nil {
        c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid username or password"})
        return
    }

    expirationTime := time.Now().Add(5 * time.Minute)
    claims := &Claims{
        Username: loginCreds.Username,
        StandardClaims: jwt.StandardClaims{
            ExpiresAt: expirationTime.Unix(),
        },
    }

    token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
    tokenString, err := token.SignedString(jwtKey)
    if err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create token"})
        return
    }

    c.JSON(http.StatusOK, gin.H{"token": tokenString, "message": "Login successful"})
}