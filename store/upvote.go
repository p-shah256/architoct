// internal/store/upvote.go

package store

import (
    "context"

	"architoct/types"
    "go.mongodb.org/mongo-driver/bson"
    "go.mongodb.org/mongo-driver/mongo"
)

// TYPE AND CREATE /////////////////////////////////////////////////////////////
type UpvoteStore struct {
    upvotes *mongo.Collection
}

func NewUpvoteStore(db *mongo.Database) *UpvoteStore {
    return &UpvoteStore{
        upvotes: db.Collection("upvotes"),
    }
}

// OPERATIONS /////////////////////////////////////////////////////////////////
func (s *UpvoteStore) Create(ctx context.Context, upvote *types.Upvote) error {
    // Check if upvote already exists
    exists, err := s.HasUserUpvoted(ctx, upvote.ParentID, upvote.UserID)
    if err != nil {
        return err
    }
    if exists {
        // return ErrAlreadyUpvoted // You'll need to define this error
    }

    _, err = s.upvotes.InsertOne(ctx, upvote)
    return err
}

// TODO: a lot of optimizations required here
func (s *UpvoteStore) HasUserUpvoted(ctx context.Context, parentID, userID string) (bool, error) {
    count, err := s.upvotes.CountDocuments(ctx, bson.M{
        "parent_id": parentID,
        "user_id":   userID,
    })
    return count > 0, err
}

func (s *UpvoteStore) GetUpvoteCount(ctx context.Context, parentID string) (int64, error) {
    count, err := s.upvotes.CountDocuments(ctx, bson.M{"parent_id": parentID})
    return count, err
}
