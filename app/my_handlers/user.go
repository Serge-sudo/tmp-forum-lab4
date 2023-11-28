package my_handlers

import (
	"context"
	"encoding/json"
	"github.com/dgrijalva/jwt-go"
	"go-chat-app/config"
	"go-chat-app/database"
	"go-chat-app/models"
	"go-chat-app/my_handlers/widgets"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"golang.org/x/crypto/bcrypt"
	"net/http"
	"strconv"
	"time"
)

func Register(w http.ResponseWriter, r *http.Request) {
	client := database.Connect()
	collection := client.Database("go_chat_app").Collection("users")
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	indexModel := mongo.IndexModel{
		Keys:    bson.M{"username": 1},
		Options: options.Index().SetUnique(true),
	}
	_, err := collection.Indexes().CreateOne(ctx, indexModel)
	if err != nil {
		http.Error(w, "Error creating unique username index: "+err.Error(), http.StatusInternalServerError)
		return
	}

	var user models.User
	err = json.NewDecoder(r.Body).Decode(&user)
	if err != nil {
		http.Error(w, "Invalid request body: "+err.Error(), http.StatusBadRequest)
		return
	}

	if len(user.Username) < 4 {
		http.Error(w, "Username must be at least 4 symbols long.", http.StatusBadRequest)
		return
	}
	if len(user.Password) < 4 {
		http.Error(w, "Password must be at least 4 symbols long.", http.StatusBadRequest)
		return
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
	if err != nil {
		http.Error(w, "Error hashing password", http.StatusInternalServerError)
		return
	}
	user.Password = string(hashedPassword)

	nextID, err := widgets.GetNextUserID(ctx, client)
	if err != nil {
		http.Error(w, "Error generating user ID", http.StatusInternalServerError)
		return
	}

	user.ID = nextID

	_, err = collection.InsertOne(ctx, user)
	if err != nil {
		// Check for duplicate username error
		if mongo.IsDuplicateKeyError(err) {
			http.Error(w, "Username already exists", http.StatusConflict)
			return
		}

		http.Error(w, "Error creating user", http.StatusInternalServerError)
		return
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"username": user.Username,
		"id":       user.ID,
		"exp":      time.Now().Add(time.Hour * 24).Unix(),
	})
	tokenString, err := token.SignedString([]byte(config.SECRET_KEY))
	if err != nil {
		http.Error(w, "Error generating token", http.StatusInternalServerError)
		return
	}
	err = json.NewEncoder(w).Encode(map[string]string{"token": tokenString})
	if err != nil {
		http.Error(w, "Error encoding token", http.StatusInternalServerError)
		return
	}
}

func Login(w http.ResponseWriter, r *http.Request) {
	var user models.User
	err := json.NewDecoder(r.Body).Decode(&user)
	if err != nil {
		http.Error(w, "Invalid request body: "+err.Error(), http.StatusBadRequest)
		return
	}
	client := database.Connect()
	collection := client.Database("go_chat_app").Collection("users")
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	var foundUser models.User
	err = collection.FindOne(ctx, bson.M{"username": user.Username}).Decode(&foundUser)
	if err != nil {
		http.Error(w, "Invalid username or password", http.StatusUnauthorized)
		return
	}
	err = bcrypt.CompareHashAndPassword([]byte(foundUser.Password), []byte(user.Password))
	if err != nil {
		http.Error(w, "Invalid username or password", http.StatusUnauthorized)
		return
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"username": foundUser.Username,
		"id":       foundUser.ID,
		"exp":      time.Now().Add(time.Hour * 24).Unix(),
	})
	tokenString, err := token.SignedString([]byte(config.SECRET_KEY))
	if err != nil {
		http.Error(w, "Error generating token", http.StatusInternalServerError)
		return
	}
	err = json.NewEncoder(w).Encode(map[string]string{"token": tokenString})
	if err != nil {
		http.Error(w, "Error encoding token", http.StatusInternalServerError)
		return
	}
}

func GetUserInfo(w http.ResponseWriter, r *http.Request) {
	token, claims, err := widgets.CheckAuth(r.Header.Get("Authorization"))

	if err == nil && token.Valid {
		userID := int(claims["id"].(float64))
		username := claims["username"].(string)
		err := json.NewEncoder(w).Encode(map[string]string{"username": username, "id": strconv.Itoa(userID)})
		if err != nil {
			http.Error(w, "Error encoding user_info", http.StatusInternalServerError)
			return
		}
	} else {
		http.Error(w, err.Error(), http.StatusUnauthorized)
	}
}
