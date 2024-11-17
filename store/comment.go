// internal/store/comment.go

package store

import (
    "context"


	"architoct/types"

    "go.mongodb.org/mongo-driver/bson"
    "go.mongodb.org/mongo-driver/mongo"
)

// TYPE AND CREATE /////////////////////////////////////////////////////////////
type CommentStore struct {
    comments *mongo.Collection
}

func NewCommentStore(db *mongo.Database) *CommentStore {
    return &CommentStore{
        comments: db.Collection("comments"),
    }
}

// OPERATIONS /////////////////////////////////////////////////////////////////
func (s *CommentStore) Create(ctx context.Context, comment *types.Comment) error {
    _, err := s.comments.InsertOne(ctx, comment)
    return err
}

func (s *CommentStore) GetById(ctx context.Context, commentId string) (*types.Comment, error) {
    var comment types.Comment
    err := s.comments.FindOne(ctx, bson.M{"_id": commentId}).Decode(&comment)
    if err != nil {
        if err == mongo.ErrNoDocuments {
            return nil, err
        }
        return nil, err
    }
    return &comment, nil
}

func (s *CommentStore) SoftDelete(ctx context.Context, commentID string) error {
	// REVIEW: what is bson.M?
	// bson.M is just a data structure = map[string]interface{}
	// basically allows us to send queries without creating structs
    updateAction := bson.M{"$set": bson.M{"is_deleted": true}}
    _, err := s.comments.UpdateOne(ctx, bson.M{"_id": commentID}, updateAction)
    return err
}

func (s *CommentStore) IncrementUpvoteCount(ctx context.Context, commentID string) error {
    updateAction := bson.M{"$inc": bson.M{"upvote_count": 1}}
    _, err := s.comments.UpdateOne(ctx, bson.M{"_id": commentID}, updateAction)
    return err
}
