package controllers

import (
	"context"

	"github.com/zachmshort/emoney-backend/config"
	"github.com/zachmshort/emoney-backend/models"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

type Transfer = models.Transfer

func PlayerTransfer(transfer Transfer) error {
	session, err := config.DB.Client().StartSession()
	if err != nil {
		return err
	}
	defer session.EndSession(context.Background())

	err = session.StartTransaction()
	if err != nil {
		return err
	}

	err = mongo.WithSession(context.Background(), session, func(sc mongo.SessionContext) error {
		playerColl := config.DB.Collection("Player")

		fromUpdate := bson.M{"$inc": bson.M{"balance": -transfer.Amount}}
		toUpdate := bson.M{"$inc": bson.M{"balance": transfer.Amount}}

		if err := playerColl.FindOneAndUpdate(sc, bson.M{"_id": transfer.FromPlayerID}, fromUpdate).Err(); err != nil {
			return err
		}

		if err := playerColl.FindOneAndUpdate(sc, bson.M{"_id": transfer.ToPlayerID}, toUpdate).Err(); err != nil {
			return err
		}

		transferColl := config.DB.Collection("Transfer")
		_, err := transferColl.InsertOne(sc, transfer)
		return err
	})

	if err != nil {
		session.AbortTransaction(context.Background())
		return err
	}

	return session.CommitTransaction(context.Background())
}

func BankTransfer(transfer Transfer) error {
	return nil
}

func RequestTransfer(transfer Transfer) error { return nil }
