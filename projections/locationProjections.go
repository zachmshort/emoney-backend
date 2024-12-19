package projections

import "go.mongodb.org/mongo-driver/bson"

func CreateLocationProjection(fields string) bson.M {
	defaultFields := []string{"_id"}

	return CreateProjection(fields, defaultFields)
}
