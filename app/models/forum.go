package models

type Forum struct {
	ID        int    `json:"id" bson:"_id"`
	Name      string `json:"name" bson:"name"`
	CreatorID int    `json:"creatorId" bson:"creatorId"`
}
