// api/users.go

package main

import (
    "encoding/json"
    "net/http"
	"time"

    "golang.org/x/crypto/bcrypt"
)

type Credentials struct {
    Email    string `json:"email"`
    Password string `json:"password"`
}

type RegistrationRequest struct {
    Email    string `json:"email"`
    Password string `json:"password"`
    Name     string `json:"name"`
}

type User struct {
    ID    string `json:"id"`
    Email string `json:"email"`
    Name  string `json:"name"`
}

func RegisterHandler(w http.ResponseWriter, r *http.Request) {
    var req RegistrationRequest
    err := json.NewDecoder(r.Body).Decode(&req)
    if err != nil {
        http.Error(w, "Invalid request payload", http.StatusBadRequest)
        return
    }

    hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
    if err != nil {
        http.Error(w, "Error hashing password", http.StatusInternalServerError)
        return
    }

    var user User
    err = db.QueryRow(
        "INSERT INTO users (email, name, password) VALUES ($1, $2, $3) RETURNING id, email, name",
        req.Email, req.Name, string(hashedPassword),
    ).Scan(&user.ID, &user.Email, &user.Name)
    if err != nil {
        http.Error(w, "Error creating user", http.StatusInternalServerError)
        return
    }

    w.Header().Set("Content-Type", "application/json")
    json.NewEncoder(w).Encode(user)
}


func LoginHandler(w http.ResponseWriter, r *http.Request) {
    var creds Credentials
    err := json.NewDecoder(r.Body).Decode(&creds)
    if err != nil {
        http.Error(w, "Invalid request payload", http.StatusBadRequest)
        return
    }

    var storedHashedPassword string
    var userID string
    err = db.QueryRow("SELECT id, password FROM users WHERE email=$1", creds.Email).Scan(&userID, &storedHashedPassword)
    if err != nil {
        http.Error(w, "User not found", http.StatusUnauthorized)
        return
    }

    err = bcrypt.CompareHashAndPassword([]byte(storedHashedPassword), []byte(creds.Password))
    if err != nil {
        http.Error(w, "Incorrect password", http.StatusUnauthorized)
        return
    }

    tokenString, expiresAt, err := GenerateToken(userID)
    if err != nil {
        http.Error(w, "Error generating token", http.StatusInternalServerError)
        return
    }

    response := map[string]interface{}{
        "token":      tokenString,
        "expires_in": expiresAt - time.Now().Unix(),
    }

    w.Header().Set("Content-Type", "application/json")
    json.NewEncoder(w).Encode(response)
}

func UserInfoHandler(w http.ResponseWriter, r *http.Request) {
    userID := r.Context().Value("userID").(string)

    var user User
    err := db.QueryRow("SELECT id, email, name FROM users WHERE id=$1", userID).Scan(&user.ID, &user.Email, &user.Name)
    if err != nil {
        http.Error(w, "User not found", http.StatusNotFound)
        return
    }

    w.Header().Set("Content-Type", "application/json")
    json.NewEncoder(w).Encode(user)
}
