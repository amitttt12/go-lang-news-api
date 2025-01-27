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

	query := `INSERT INTO users (username, email, password, phone, profile_picture) 
              VALUES (?, ?, ?, ?, ?)`

	_, err = db.Exec(query, user.Username, user.Email, user.Password, user.Phone, user.ProfilePicture)

	if err != nil {
		if mysqlErr, ok := err.(*mysql.MySQLError); ok && mysqlErr.Number == 1062 {
			http.Error(w, "Email already exists", http.StatusConflict)
		} else {
			http.Error(w, "Error creating user", http.StatusInternalServerError)
		}
		return
	}

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(map[string]string{"message": "User created successfully"})
}
