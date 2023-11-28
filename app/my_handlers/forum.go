package my_handlers

import (
	"context"
	"encoding/json"
	"go-chat-app/database"
	"go-chat-app/models"
	"go-chat-app/my_handlers/widgets"
	"go.mongodb.org/mongo-driver/bson"
	"net/http"
	"strconv"
	"time"
)

func CreateForum(w http.ResponseWriter, r *http.Request) {

	token, claims, err := widgets.CheckAuth(r.Header.Get("Authorization"))

	if err == nil && token.Valid {
		userID := int(claims["id"].(float64))
		var forum models.Forum
		client := database.Connect()
		collection := client.Database("go_chat_app").Collection("forums")

		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()

		forum.Name = ""
		err2 := json.NewDecoder(r.Body).Decode(&forum)
		if err2 != nil {
			http.Error(w, "Invalid request body: "+err2.Error(), http.StatusBadRequest)
			return
		}
		if len(forum.Name) == 0 {
			http.Error(w, "Forum name can't be empty", http.StatusBadRequest)
			return
		}

		nextID, err := widgets.GetNextForumID(ctx, client)
		if err != nil {
			http.Error(w, "Error generating user ID", http.StatusInternalServerError)
			return
		}

		forum.ID = nextID
		forum.CreatorID = userID

		_, err = collection.InsertOne(ctx, forum)
		if err != nil {
			http.Error(w, "Error creating forum", http.StatusInternalServerError)
			return
		}

		cursor, err := collection.Find(ctx, bson.M{})
		if err != nil {
			http.Error(w, "Error fetching forums", http.StatusInternalServerError)
			return
		}

		var forums []models.Forum
		err = cursor.All(ctx, &forums)
		if err != nil {
			http.Error(w, "Error decoding forums", http.StatusInternalServerError)
			return
		}

		err = json.NewEncoder(w).Encode(forums)
		if err != nil {
			http.Error(w, "Error encoding forums", http.StatusInternalServerError)
			return
		}
	} else {
		http.Error(w, err.Error(), http.StatusUnauthorized)
	}
}

func GetForums(w http.ResponseWriter, r *http.Request) {
	token, _, err := widgets.CheckAuth(r.Header.Get("Authorization"))

	if err == nil && token.Valid {
		client := database.Connect()
		collection := client.Database("go_chat_app").Collection("forums")

		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()

		cursor, err := collection.Find(ctx, bson.M{})
		if err != nil {
			http.Error(w, "Error fetching forums", http.StatusInternalServerError)
			return
		}

		var forums []models.Forum
		err = cursor.All(ctx, &forums)
		if err != nil {
			http.Error(w, "Error decoding forums", http.StatusInternalServerError)
			return
		}

		err = json.NewEncoder(w).Encode(forums)
		if err != nil {
			http.Error(w, "Error encoding forums", http.StatusInternalServerError)
			return
		}
	} else {
		http.Error(w, err.Error(), http.StatusUnauthorized)
	}

}

func RemoveForum(w http.ResponseWriter, r *http.Request) {
	token, claims, err := widgets.CheckAuth(r.Header.Get("Authorization"))

	if err == nil && token.Valid {
		userID := int(claims["id"].(float64))
		forumIDStr := r.URL.Query().Get("forumId")
		forumID, err := strconv.Atoi(forumIDStr)

		if err != nil {
			http.Error(w, "Invalid forum ID", http.StatusBadRequest)
			return
		}

		client := database.Connect()
		forumsCollection := client.Database("go_chat_app").Collection("forums")
		messagesCollection := client.Database("go_chat_app").Collection("messages")

		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()

		var forum models.Forum
		err = forumsCollection.FindOne(ctx, bson.M{"_id": forumID}).Decode(&forum)
		if err != nil {
			http.Error(w, "Forum not found", http.StatusNotFound)
			return
		}

		if forum.CreatorID != userID {
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}

		_, err = forumsCollection.DeleteOne(ctx, bson.M{"_id": forumID})
		if err != nil {
			http.Error(w, "Error deleting forum", http.StatusInternalServerError)
			return
		}

		_, err = messagesCollection.DeleteMany(ctx, bson.M{"forumId": forumID})
		if err != nil {
			http.Error(w, "Error deleting messages", http.StatusInternalServerError)
			return
		}

		cursor, err := forumsCollection.Find(ctx, bson.M{})
		if err != nil {
			http.Error(w, "Error fetching forums", http.StatusInternalServerError)
			return
		}

		var forums []models.Forum
		err = cursor.All(ctx, &forums)
		if err != nil {
			http.Error(w, "Error decoding forums", http.StatusInternalServerError)
			return
		}

		err = json.NewEncoder(w).Encode(forums)
		if err != nil {
			http.Error(w, "Error encoding forums", http.StatusInternalServerError)
			return
		}

	} else {
		http.Error(w, err.Error(), http.StatusUnauthorized)
	}
}
