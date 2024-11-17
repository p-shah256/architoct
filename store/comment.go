// internal/store/comment.go

package store

import (
    "context"


	"architoct/types"

    "go.mongodb.org/mongo-driver/bson"
    "go.mongodb.org/mongo-driver/mongo"
    "go.mongodb.org/mongo-driver/mongo/options"
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

func (s *CommentStore) GetCommentsByPostID(ctx context.Context, postID string, limit int64, offset int64) ([]*types.Comment, error) {
    opts := options.Find().
        SetSort(bson.D{{Key: "created_at", Value: -1}}).
        SetLimit(limit).
        SetSkip(offset)

	// TODO: MVP: do not need to do this
	// 1. get the comment ids inside post
	// 2. return all those comment ids
	//
	// REVIEW: cursor is like a stream of response
	// 		   next returns one item and all returns all
    cursor, err := s.comments.Find(ctx,
        bson.M{
            "post_id": postID,
            "is_deleted": false,
        },
        opts,
    )
    if err != nil {
        return nil, err
    }
    defer cursor.Close(ctx)

    var comments []*types.Comment
    if err = cursor.All(ctx, &comments); err != nil {
        return nil, err
    }
    return comments, nil
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
