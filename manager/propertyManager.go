package manager

import (
	"context"
	"errors"
	"fmt"
	"log"
	"strings"

	"github.com/zachmshort/monopoly-backend/config"
	"github.com/zachmshort/monopoly-backend/models"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func HandleHotelManagement(roomObjID primitive.ObjectID, manageType string, propertyDetails models.PropertyDetails) error {
	propertyColl := config.DB.Collection("Property")

	for _, detail := range propertyDetails {
		propertyObjID, err := primitive.ObjectIDFromHex(detail.PropertyID)
		if err != nil {
			return fmt.Errorf("invalid property ID: %s", detail.PropertyID)
		}

		filter := bson.M{"_id": propertyObjID, "roomId": roomObjID}
		update := bson.M{}

		switch manageType {
		case "ADD_HOTELS":
			update = bson.M{
				"$set": bson.M{"developmentLevel": 5},
			}
		case "REMOVE_HOTELS":
			update = bson.M{"$set": bson.M{"developmentLevel": 5}}
		default:
			return fmt.Errorf("invalid management type: %s", manageType)
		}

		_, err = propertyColl.UpdateOne(context.TODO(), filter, update)
		if err != nil {
			return fmt.Errorf("failed to update property: %w", err)
		}
	}

	return nil
}

func HandleHouseManagement(roomObjID primitive.ObjectID, manageType string, propertyDetails models.PropertyDetails) error {
	propertyColl := config.DB.Collection("Property")

	for _, detail := range propertyDetails {
		propertyObjID, err := primitive.ObjectIDFromHex(detail.PropertyID)
		if err != nil {
			return fmt.Errorf("invalid property ID: %s", detail.PropertyID)
		}
		filter := bson.M{"_id": propertyObjID, "roomId": roomObjID}
		update := bson.M{}
		switch manageType {
		case "ADD_HOUSES":
			update = bson.M{"$set": bson.M{"developmentLevel": detail.Count}}
		case "REMOVE_HOUSES":
			update = bson.M{"$set": bson.M{"developmentLevel": detail.Count}}
		default:
			return fmt.Errorf("invalid management type: %s", manageType)
		}

		_, err = propertyColl.UpdateOne(context.TODO(), filter, update)
		if err != nil {
			return fmt.Errorf("failed to update property: %w", err)
		}
	}

	return nil
}

func HandlePropertySaleMortgage(roomObjID primitive.ObjectID, manageType string, propertyDetails models.PropertyDetails) error {
	manageType = strings.TrimSpace(manageType)
	propertyColl := config.DB.Collection("Property")
	for _, detail := range propertyDetails {
		propertyObjID, err := primitive.ObjectIDFromHex(detail.PropertyID)
		if err != nil {
			return fmt.Errorf("invalid property id")
		}

		filter := bson.M{"_id": propertyObjID, "roomId": roomObjID}
		update := bson.M{}
		switch manageType {
		case "MORTGAGE":
			update = bson.M{"$set": bson.M{"isMortgaged": true}}
		case "UNMORTGAGE":
			update = bson.M{"$set": bson.M{"isMortgaged": false}}
		case "SELL":
			update = bson.M{"$set": bson.M{"playerId": nil, "isMortgaged": false}}
		default:
			return fmt.Errorf("invalid management type, %s", manageType)
		}

		result, err := propertyColl.UpdateOne(context.TODO(), filter, update)
		if err != nil {
			log.Printf("Failed to update property for %s: %v", manageType, err)
			return fmt.Errorf("failed to update property for %s: %w", manageType, err)
		}
		log.Printf("ManageType: %s | Matched: %d | Modified: %d", manageType, result.MatchedCount, result.ModifiedCount)
	}

	return nil
}

func ExtractPropertyDetails(raw interface{}) (models.PropertyDetails, error) {
	properties, ok := raw.([]interface{})
	if !ok {
		return nil, errors.New("properties field is not an array")
	}

	var propertyDetails models.PropertyDetails
	for _, prop := range properties {
		propMap, ok := prop.(map[string]interface{})
		if !ok {
			return nil, errors.New("property item is not a valid object")
		}

		propertyID, ok := propMap["propertyId"].(string)
		if !ok {
			return nil, errors.New("missing or invalid propertyId field")
		}

		count, ok := propMap["count"].(float64)
		if !ok {
			count = 0
		}

		propertyDetails = append(propertyDetails, struct {
			PropertyID string
			Count      int
		}{
			PropertyID: propertyID,
			Count:      int(count),
		})
	}

	return propertyDetails, nil
}
