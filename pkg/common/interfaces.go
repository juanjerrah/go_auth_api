package common

import "go.mongodb.org/mongo-driver/bson/primitive"

type PasswordHasher interface {
	Hash(password string) (string, error)
	Verify(password, hash string) error
}

type MongoUtils interface {
	ToObjectID(id string) primitive.ObjectID
	GenerateObjectID() primitive.ObjectID
}
