package utils

import "go.mongodb.org/mongo-driver/bson/primitive"

type MongoUtils interface {
	ToObjectID(id string) (primitive.ObjectID)
	GenerateObjectID() (primitive.ObjectID)
}

type mongoUtils struct{}

func NewMongoUtils() MongoUtils {
	return &mongoUtils{}
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

