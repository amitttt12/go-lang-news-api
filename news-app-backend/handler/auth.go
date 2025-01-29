package handler

import (
	"encoding/json"
	"net/http"
	"news-app/config"
	model "news-app/models"

	"github.com/go-sql-driver/mysql"
	"golang.org/x/crypto/bcrypt"
)

//function to signup the users

func Signup(w http.ResponseWriter, r *http.Request) {
	var user model.User

	err := json.NewDecoder(r.Body).Decode(&user)
	if err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
	if err != nil {
		http.Error(w, "Error hashing password", http.StatusInternalServerError)
		return
	}
	user.Password = string(hashedPassword)

	//create a connection to database

	db, err := config.InitDatabase()

	if err != nil {
		http.Error(w, "Could not connect to database", http.StatusInternalServerError)
		return
	}
	defer db.Close()

	query := `INSERT INTO users (username, email, password, phone, profile_picture,isAdmin) 
              VALUES (?, ?, ?, ?, ?, ?)`

	_, err = db.Exec(query, user.Username, user.Email, user.Password, user.Phone, user.ProfilePicture, user.IsAdmin)

	if err != nil {
		if mysqlErr, ok := err.(*mysql.MySQLError); ok && mysqlErr.Number == 1062 {
			http.Error(w, "Email already exists", http.StatusConflict)
		} else {
			http.Error(w, "Error creating user", http.StatusInternalServerError)
		}
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{"message": "User created successfully"})
}

//Sign in function

func Signin(w http.ResponseWriter, r *http.Request) {
	var loginData struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}

	// Decode the JSON request body
	err := json.NewDecoder(r.Body).Decode(&loginData)
	if err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	// Initialize database connection
	db, err := config.InitDatabase()
	if err != nil {
		http.Error(w, "Could not connect to database", http.StatusInternalServerError)
		return
	}
	defer db.Close()

	// Fetch the user data for the given email
	var user model.User
	var hashedPassword string
	query := "SELECT id, username, email, phone, profile_picture, isAdmin, password FROM users WHERE email = ?"
	err = db.QueryRow(query, loginData.Email).Scan(
		&user.ID, &user.Username, &user.Email, &user.Phone, &user.ProfilePicture, &user.IsAdmin, &hashedPassword,
	)
	if err != nil {
		http.Error(w, "Invalid email or password", http.StatusUnauthorized)
		return
	}

	// Compare the hashed password with the password provided by the client
	err = bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(loginData.Password))
	if err != nil {
		http.Error(w, "Invalid email or password", http.StatusUnauthorized)
		return
	}

	// Respond with the user's data
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(user)
}
