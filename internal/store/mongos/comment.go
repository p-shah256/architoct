// internal/store/comment.go

package mongos

import (
	"architoct/internal/logger"
	"architoct/internal/types"
	"context"

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
func (s *CommentStore) InsertToDB(ctx context.Context, comment *types.Comment) (string, error) {
	result, err := s.comments.InsertOne(ctx, comment)
	if err != nil {
		return "", err
	}
	logger.Debug().
		Str("id", result.InsertedID.(primitive.ObjectID).Hex()).
		Msg("returning ID")
	id := result.InsertedID.(primitive.ObjectID).Hex()
	return id, nil
}

func (s *CommentStore) GetById(ctx context.Context, commentId string) (*types.Comment, error) {
	var comment types.Comment
	id, err := primitive.ObjectIDFromHex(commentId)
	if err != nil {
		return nil, err
	}
	err = s.comments.FindOne(ctx, bson.M{"_id": id}).Decode(&comment)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, err
		}
		return nil, err
	}
	return &comment, nil
}

func (s *CommentStore) AddToRepliesArray(ctx context.Context, parentCommentId string, commentId string) error {
    commentid, _ := primitive.ObjectIDFromHex(commentId)
    parentid, _ := primitive.ObjectIDFromHex(parentCommentId)
    _ , err := s.comments.UpdateOne(
        ctx,
        bson.M{"_id": parentid},
        bson.M{
            "$push": bson.M{"replies": commentid},
            "$inc": bson.M{"reply_count": 1},
        },
    )
    if err != nil {
        if err == mongo.ErrNoDocuments {
            return err
        }
		return err
    }
    return nil
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
	id, _ := primitive.ObjectIDFromHex(commentID)
    // First try to add upvote
    var updatedComment types.Comment
	err := s.comments.FindOneAndUpdate(
        ctx,
        bson.M{"_id": id, "upvoted_by": bson.M{"$ne": userID}},
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
            bson.M{"_id": id},
            bson.M{
                "$inc": bson.M{"upvote_count": -1},
                "$pull": bson.M{"upvoted_by": userID},
            },
            options.FindOneAndUpdate().SetReturnDocument(options.After),
        ).Decode(&updatedComment)
        return updatedComment, err
    }

    return updatedComment, nil
}
