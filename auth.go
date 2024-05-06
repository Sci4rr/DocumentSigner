package main

import (
    "context"
    "encoding/json"
    "fmt"
    "log"
    "net/http"
    "os"
    "time"

    "github.com/gorilla/mux"
    "golang.org/x/crypto/bcrypt"
    "github.com/golang-jwt/jwt/v4"
)

var jwtKey = []byte(os.Getenv("JWT_SECRET_KEY"))

type Credentials struct {
    Password string `json:"password"`
    Username string `json:"username"`
}

type Claims struct {
    Username string `json:"username"`
    jwt.StandardClaims
}

type UpdatePasswordRequest struct {
    NewPassword string `json:"newPassword"`
}

var users = map[string]string{}

func main() {
    router := mux.NewRouter()

    router.HandleFunc("/register", registerHandler).Methods("POST")
    router.HandleFunc("/login", loginHandler).Methods("POST")
    router.HandleFunc("/welcome", authMiddleware(welcomeHandler)).Methods("GET")
    router.HandleFunc("/updatePassword", authMiddleware(updatePasswordHandler)).Methods("POST")

    log.Fatal(http.ListenAndServe(":8080", router))
}

func registerHandler(w http.ResponseWriter, r *http.Request) {
    var creds Credentials
    err := json.NewDecoder(r.Body).Decode(&creds)
    if err != nil {
        http.Error(w, "Bad request", http.StatusBadRequest)
        return
    }
    defer r.Body.Close()

    hashedPassword, err := bcrypt.GenerateFromPassword([]byte(creds.Password), bcrypt.DefaultCost)
    if err != nil {
        http.Error(w, "Internal server error", http.StatusInternalServerError)
        return
    }

    users[creds.Username] = string(hashedPassword)

    w.WriteHeader(http.StatusCreated)
}

func loginHandler(w http.ResponseWriter, r *http.Request) {
    var creds Credentials
    err := json.NewDecoder(r.Body).Decode(&creds)
    if err != nil {
        http.Error(w, "Bad request", http.StatusBadRequest)
        return
    }
    defer r.Body.Close()

    expectedPassword, ok := users[creds.Username]

    if !ok || bcrypt.CompareHashAndPassword([]byte(expectedPassword), []byte(creds.Password)) != nil {
        http.Error(w, "Unauthorized", http.StatusUnauthorized)
        return
    }

    expirationTime := time.Now().Add(5 * time.Minute)
    claims := &Claims{
        Username: creds.Username,
        StandardClaims: jwt.StandardClaims{
            ExpiresAt: expirationTime.Unix(),
        },
    }

    token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
    tokenString, err := token.SignedString(jwtKey)

    if err != nil {
        http.Error(w, "Internal server error", http.StatusInternalServerError)
        return
    }

    http.SetCookie(w, &http.Cookie{
        Name:    "token",
        Value:   tokenString,
        Expires: expirationTime,
    })
}

func welcomeHandler(w http.ResponseWriter, r *http.Request) {
    _, err := w.Write([]byte("Welcome!"))
    if err != nil {
        log.Printf("Failed to write response: %v", err)
    }
}

func updatePasswordHandler(w http.ResponseWriter, r *http.Request) {
    var updatePasswordReq UpdatePasswordRequest
    var claims = r.Context().Value("claims").(*Claims)

    err := json.NewDecoder(r.Body).Decode(&updatePasswordReq)
    if err != nil {
        http.Error(w, "Bad request", http.StatusBadRequest)
        return
    }
    defer r.Body.Close()

    hashedPassword, err := bcrypt.GenerateFromPassword([]byte(updatePasswordReq.NewPassword), bcrypt.DefaultCost)
    if err != nil {
        http.Error(w, "Internal server error", http.StatusInternalServerError)
        return
    }

    users[claims.Username] = string(hashedPassword)

    w.WriteHeader(http.StatusOK)
    message := fmt.Sprintf("Password updated successfully for user %s", claims.Username)
    _, err = w.Write([]byte(message))
    if err != nil {
        log.Printf("Failed to write response: %v", err)
    }
}

func authMiddleware(next http.HandlerFunc) http.HandlerFunc {
    return func(w http.ResponseWriter, r *http.Request) {
        cookie, err := r.Cookie("token")
        if err != nil {
            if err == http.ErrNoCookie {
                http.Error(w, "Unauthorized", http.StatusUnauthorized)
                return
            }
            http.Error(w, "Bad request", http.StatusBadRequest)
            return
        }

        tokenStr := cookie.Value
        claims := &Claims{}

        tkn, err := jwt.ParseWithClaims(tokenStr, claims, func(token *jwt.Token) (interface{}, error) {
            if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
                return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
            }
            return jwtKey, nil
        })

        if err != nil {
            if err == jwt.ErrSignatureInvalid {
                http.Error(w, "Unauthorized", http.StatusUnauthorized)
                return
            }
            http.Error(w, "Bad request", http.StatusBadRequest)
            return
        }

        if !tkn.Valid {
            http.Error(w, "Unauthorized", http.StatusUnauthorized)
            return
        }

        ctx := context.WithValue(r.Context(), "claims", claims)
        next.ServeHTTP(w, r.WithContext(ctx))
    }
}