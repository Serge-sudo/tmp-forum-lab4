package widgets

import (
	"context"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func GetNextForumID(ctx context.Context, client *mongo.Client) (int, error) {
	counterCollection := client.Database("go_chat_app").Collection("counters")
	filter := bson.M{"_id": "forumID"}
	update := bson.M{"$inc": bson.M{"seq": 1}}
	options_ := options.FindOneAndUpdate().SetUpsert(true).SetReturnDocument(options.After)

	var counterDoc bson.M
	err := counterCollection.FindOneAndUpdate(ctx, filter, update, options_).Decode(&counterDoc)
	if err != nil {
		return 0, err
	}

	return int(counterDoc["seq"].(int32)), nil
}

func GetNextMessageID(ctx context.Context, client *mongo.Client) (int, error) {
	counterCollection := client.Database("go_chat_app").Collection("counters")
	filter := bson.M{"_id": "msgID"}
	update := bson.M{"$inc": bson.M{"seq": 1}}
	options_ := options.FindOneAndUpdate().SetUpsert(true).SetReturnDocument(options.After)

	var counterDoc bson.M
	err := counterCollection.FindOneAndUpdate(ctx, filter, update, options_).Decode(&counterDoc)
	if err != nil {
		return 0, err
	}

	return int(counterDoc["seq"].(int32)), nil
}

func GetNextUserID(ctx context.Context, client *mongo.Client) (int, error) {
	counterCollection := client.Database("go_chat_app").Collection("counters")
	filter := bson.M{"_id": "userID"}
	update := bson.M{"$inc": bson.M{"seq": 1}}
	options_ := options.FindOneAndUpdate().SetUpsert(true).SetReturnDocument(options.After)

	var counterDoc bson.M
	err := counterCollection.FindOneAndUpdate(ctx, filter, update, options_).Decode(&counterDoc)
	if err != nil {
		return 0, err
	}

	return int(counterDoc["seq"].(int32)), nil
}
