package mongo

import "go.mongodb.org/mongo-driver/bson/primitive"

func IsValidObjectID(id string) bool {
	_, err := primitive.ObjectIDFromHex(id)
	return err == nil
}
