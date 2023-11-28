package models

import "time"

type Message struct {
	ID        int       `json:"id" bson:"_id"`
	Username  string    `json:"username" bson:"username"`
	Message   string    `json:"message" bson:"message"`
	ForumID   int       `json:"forumId" bson:"forumId"`
	Timestamp time.Time `json:"timestamp" bson:"timestamp"`
}
