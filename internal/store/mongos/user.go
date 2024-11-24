package mongos

import (
	"context"
	"time"

	"architoct/internal/logger"
	"architoct/internal/types"

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
func (s *UserStore) Create(ctx context.Context, userid string) (*types.User, error) {
	logger.Debug().Str("userid", userid)
	var user types.User
	user = types.User{
		ID: userid,
		CreatedAt:   time.Now(),
		LastLogin:   time.Now(),
	}
	_, err := s.users.InsertOne(ctx, user)
	logger.Debug().Err(err).Msg("user create status")
	return &user, err
}
