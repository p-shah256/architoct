package mongos

import (
	"context"
	"log/slog"
	"time"

	"architoct/internal/types"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
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
	var user types.User
	err := s.users.FindOne(ctx, bson.M{"_id": userid}).Decode(&user)
	// if found
	if err == nil {
		return &user, err
	}

	// else Create new user
	user = types.User{
		ID: userid,
		CreatedAt:   time.Now(),
		LastLogin:   time.Now(),
	}
	_, err = s.users.InsertOne(ctx, user)
	slog.Info("user create status", "err", err)
	return &user, err
}

func (s *UserStore) UpdateName(ctx context.Context, userid string, username string) (*types.User, error) {
	slog.Info("updating name......", )
    var updatedUser types.User
    err := s.users.FindOneAndUpdate(ctx,
        bson.M{"_id": userid},
        bson.M{"$set": bson.M{"user_name": username}},
        options.FindOneAndUpdate().SetReturnDocument(options.After),
    ).Decode(&updatedUser)

    if err != nil {
        switch {
        case mongo.IsDuplicateKeyError(err):
			slog.Error("already taken")
            return nil, types.ErrUsernameTaken
        case err == mongo.ErrNoDocuments:
            return nil, types.ErrUsernameTaken
        default:
            return nil, err
        }
    }

	slog.Info("updating name...... no erros", )
    return &updatedUser, nil
}
