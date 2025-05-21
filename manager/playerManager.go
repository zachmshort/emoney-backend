package manager

import (
	"context"
	"fmt"

	"github.com/zachmshort/emoney-backend/config"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func UpdatePlayerBalance(playerID primitive.ObjectID, amount int) error {
	playersColl := config.DB.Collection("Player")
	filter := bson.M{"_id": playerID}
	update := bson.M{"$inc": bson.M{"balance": -amount}}

	_, err := playersColl.UpdateOne(context.TODO(), filter, update)
	if err != nil {
		return fmt.Errorf("failed to update player balance: %w", err)
	}

	return nil
}
