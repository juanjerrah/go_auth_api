package utils

import (
	"github.com/juanjerrah/go_auth_api/pkg/common"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type mongoUtils struct{}

func NewMongoUtils() common.MongoUtils {
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
