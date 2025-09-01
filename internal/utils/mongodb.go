package utils

import "go.mongodb.org/mongo-driver/bson/primitive"

type MongoUtils interface {
	ToObjectID(id string) (primitive.ObjectID)
	FromObjectID(id primitive.ObjectID) (string)
	GenerateObjectID() (primitive.ObjectID)
}

type mongoUtils struct{}

func NewMongoUtils() MongoUtils {
	return &mongoUtils{}
}

// FromObjectID implements MongoUtils.
func (m *mongoUtils) FromObjectID(id primitive.ObjectID) string {
	return id.Hex()
}

// GenerateObjectID implements MongoUtils.
func (m *mongoUtils) GenerateObjectID() primitive.ObjectID {
	return primitive.NewObjectID()
}

// ToObjectID implements MongoUtils.
func (m *mongoUtils) ToObjectID(id string) primitive.ObjectID {
	objectID, _ := primitive.ObjectIDFromHex(id)
	return objectID
}

