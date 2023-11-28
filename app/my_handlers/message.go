package my_handlers

import (
	"context"
	"encoding/json"
	"go-chat-app/database"
	"go-chat-app/models"
	"go-chat-app/my_handlers/widgets"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo/options"
	"net/http"
	"strconv"
	"time"
)

func CreateMessage(w http.ResponseWriter, r *http.Request) {
	token, claims, err := widgets.CheckAuth(r.Header.Get("Authorization"))

	if err == nil && token.Valid {
		var message models.Message
		client := database.Connect()
		collection := client.Database("go_chat_app").Collection("messages")
		forumsCollection := client.Database("go_chat_app").Collection("forums")

		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()

		message.Message = ""
		message.ForumID = -1
		err = json.NewDecoder(r.Body).Decode(&message)
		if err != nil {
			http.Error(w, "Invalid request body: "+err.Error(), http.StatusBadRequest)
			return
		}
		if len(message.Message) == 0 {
			http.Error(w, "Message can't be empty", http.StatusBadRequest)
			return
		}
		if message.ForumID == -1 {
			http.Error(w, "ForumId should be provided", http.StatusBadRequest)
			return
		}

		opts := options.Find()
		opts.SetSort(bson.D{{"_id", -1}})

		cursorForum, err := forumsCollection.Find(ctx, bson.M{"_id": message.ForumID}, opts)
		if err != nil {
			http.Error(w, "Error fetching forum", http.StatusInternalServerError)
			return
		}
		if cursorForum.RemainingBatchLength() == 0 {
			http.Error(w, "Invalid forumId", http.StatusNotFound)
			return
		}

		nextID, err := widgets.GetNextMessageID(ctx, client)
		if err != nil {
			http.Error(w, "Error generating user ID", http.StatusInternalServerError)
			return
		}
		message.ID = nextID

		message.Timestamp = time.Now()

		message.Username = claims["username"].(string)

		_, err = collection.InsertOne(ctx, message)
		if err != nil {
			http.Error(w, "Error creating message", http.StatusInternalServerError)
			return
		}

		cursor, err := collection.Find(ctx, bson.M{"forumId": message.ForumID}, opts)
		if err != nil {
			http.Error(w, "Error fetching messages", http.StatusInternalServerError)
			return
		}

		var messages []models.Message
		err = cursor.All(ctx, &messages)
		if err != nil {
			http.Error(w, "Error decoding messages", http.StatusInternalServerError)
			return
		}
		err = json.NewEncoder(w).Encode(messages)
		if err != nil {
			http.Error(w, "Error encoding messages", http.StatusInternalServerError)
			return
		}
	} else {
		http.Error(w, err.Error(), http.StatusUnauthorized)
	}
}

func GetMessages(w http.ResponseWriter, r *http.Request) {
	token, _, err := widgets.CheckAuth(r.Header.Get("Authorization"))

	if err == nil && token.Valid {
		client := database.Connect()
		collection := client.Database("go_chat_app").Collection("messages")
		collectionForum := client.Database("go_chat_app").Collection("forums")
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()
		opts := options.Find()
		opts.SetSort(bson.D{{"_id", -1}})
		forumIDStr := r.URL.Query().Get("forumId")
		forumID, err := strconv.Atoi(forumIDStr)
		if err != nil {
			http.Error(w, "Invalid forumId format", http.StatusNotFound)
			return
		}

		cursorForum, err := collectionForum.Find(ctx, bson.M{"_id": forumID}, opts)
		if err != nil {
			http.Error(w, "Error fetching messages", http.StatusInternalServerError)
			return
		}
		if cursorForum.RemainingBatchLength() == 0 {
			http.Error(w, "Invalid forumId", http.StatusNotFound)
			return
		}

		cursor, err := collection.Find(ctx, bson.M{"forumId": forumID}, opts)
		if err != nil {
			http.Error(w, "Error fetching messages", http.StatusInternalServerError)
			return
		}

		var messages []models.Message
		err = cursor.All(ctx, &messages)
		if err != nil {
			http.Error(w, "Error decoding messages", http.StatusInternalServerError)
			return
		}
		err = json.NewEncoder(w).Encode(messages)
		if err != nil {
			http.Error(w, "Error encoding messages", http.StatusInternalServerError)
			return
		}
	} else {
		http.Error(w, err.Error(), http.StatusUnauthorized)
	}
}
