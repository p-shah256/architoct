// internal/store/comment.go

package mongos

import (
	"architoct/internal/types"
	"context"
	"log/slog"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
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
func (s *CommentStore) Create(ctx context.Context, comment *types.Comment) (string, error) {
	result, err := s.comments.InsertOne(ctx, comment)
	if err != nil {
		return "", err
	}
	slog.Info("returning ID", "id", result.InsertedID.(primitive.ObjectID).Hex())
	id := result.InsertedID.(primitive.ObjectID).Hex()
	slog.Info("returning ID", "id", id)
	return id, nil
}

func (s *CommentStore) GetById(ctx context.Context, commentId primitive.ObjectID) (*types.Comment, error) {
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

func (s *CommentStore) ToggleUpvote(ctx context.Context, commentID string, userID string) (types.Comment, error) {
    // First try to add upvote
    var updatedComment types.Comment
	err := s.comments.FindOneAndUpdate(
        ctx,
        bson.M{"_id": commentID, "upvoted_by": bson.M{"$ne": userID}},
        bson.M{
            "$inc": bson.M{"upvote_count": 1},
            "$addToSet": bson.M{"upvoted_by": userID},
        },
        options.FindOneAndUpdate().SetReturnDocument(options.After),
    ).Decode(&updatedComment)

    if err == mongo.ErrNoDocuments {
        // If user already upvoted, remove the upvote, and return false
        err = s.comments.FindOneAndUpdate(
            ctx,
            bson.M{"_id": commentID},
            bson.M{
                "$inc": bson.M{"upvote_count": -1},
                "$pull": bson.M{"upvoted_by": userID},
            },
            options.FindOneAndUpdate().SetReturnDocument(options.After),
        ).Decode(&updatedComment)
        return updatedComment, err
    }

	// else return true
    return updatedComment, err
}
