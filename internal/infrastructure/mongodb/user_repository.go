package mongodb

import (
	"context"

	"github.com/juanjerrah/go_auth_api/internal/domain/user"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type UserRepository struct {
	collection *mongo.Collection
}

func NewUserRepository(db *mongo.Database) user.Repository {
	return &UserRepository{
		collection: db.Collection("users"),
	}
}

// Create implements user.Repository.
func (u *UserRepository) Create(ctx context.Context, user *user.User) error {
	_, err := u.collection.InsertOne(ctx, user)
	return err
}

// Delete implements user.Repository.
func (u *UserRepository) Delete(ctx context.Context, id primitive.ObjectID) error {
	_, err := u.collection.DeleteOne(ctx, bson.M{"_id": id})
	return err
}

// ExistsByEmail implements user.Repository.
func (u *UserRepository) ExistsByEmail(ctx context.Context, email string) (bool, error) {
	count, err := u.collection.CountDocuments(ctx, bson.M{"email": email})
	if err != nil {
		return false, err
	}
	return count > 0, nil
}

// FindByEmail implements user.Repository.
func (u *UserRepository) FindByEmail(ctx context.Context, email string) (*user.User, error) {
	var usr user.User
	err := u.collection.FindOne(ctx, bson.M{"email": email}).Decode(&usr)
	if err != nil {
		return nil, err
	}
	return &usr, nil
}

// FindByID implements user.Repository.
func (u *UserRepository) FindByID(ctx context.Context, id primitive.ObjectID) (*user.User, error) {
	var usr user.User
	err := u.collection.FindOne(ctx, bson.M{"_id": id}).Decode(&usr)
	if err != nil {
		return nil, err
	}
	return &usr, nil
}

// Update implements user.Repository.
func (u *UserRepository) Update(ctx context.Context, user *user.User) error {
	_, err := u.collection.UpdateOne(ctx, bson.M{"_id": user.ID}, bson.M{"$set": user})
	return err
}

