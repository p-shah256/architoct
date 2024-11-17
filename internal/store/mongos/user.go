package mongos

import (
	"context"
	"time"

	"architoct/internal/types"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

// TYPE AND CREATE /////////////////////////////////////////////////////////////
type UserStore struct {
	users *mongo.Collection
}

func NewUserStore(db *mongo.Database) *UserStore {
	return &UserStore{
		users: db.Collection("users"),
	}
}

// OPERATIONS /////////////////////////////////////////////////////////////////
func (s *UserStore) GetOrCreate(ctx context.Context, fingerprint string) (*types.User, error) {
	// Try to find existing user
	var user types.User
	err := s.users.FindOne(ctx, bson.M{"fingerprint": fingerprint}).Decode(&user)
	// if found
	if err == nil {
		// Update last login
		_, err = s.users.UpdateOne(ctx,
			bson.M{"fingerprint": fingerprint},
			bson.M{"$set": bson.M{"last_login": time.Now()}},
		)
		return &user, err
	}

	// else Create new user
	user = types.User{
		ID: fingerprint,
		CreatedAt:   time.Now(),
		LastLogin:   time.Now(),
	}
	_, err = s.users.InsertOne(ctx, user)
	return &user, err
}
