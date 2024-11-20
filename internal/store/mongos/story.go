// internal/store/post.go
// no abstraction for mvp -- eventually we could have a port called port store and then let 1. mongo implment it and then 2. dynamo and so on...
// abstract away *mongo.Collection

package mongos

import (
	"context"
	"fmt"

	"architoct/internal/types"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// TYPE AND CREATE /////////////////////////////////////////////////////////////
type StoryStore struct {
	stories *mongo.Collection
}

func NewStoryStore(db *mongo.Database) *StoryStore {
	return &StoryStore{
		stories: db.Collection("stories"),
	}
}

// OPERATIONS /////////////////////////////////////////////////////////////////
// Core operations that we'll need for MVP
func (s *StoryStore) Create(ctx context.Context, story *types.Story) error {
	// For MVP, we'll do simple inserts without transactions
	_, err := s.stories.InsertOne(ctx, story)
	return err
}

func (s *StoryStore) GetByID(ctx context.Context, id string) (*types.Story, error) {
	var story types.Story
	err := s.stories.FindOne(ctx, bson.M{"_id": id}).Decode(&story)
	return &story, err
}

// GetRecent supports our homepage feed
// limit: how many posts to fetch
// timeRange: "hour", "day", "week" etc
func (s *StoryStore) GetRecent(ctx context.Context, limit int64) ([]types.Story, error) {
    // First, let's count total documents
    opts := options.Find().
        SetSort(bson.D{{Key: "upvote_count", Value: -1}}).
        SetLimit(limit)

    // Let's also print what a sample document looks like
    cursor, err := s.stories.Find(ctx, bson.D{}, opts)
    if err != nil {
        return nil, fmt.Errorf("find error: %w", err)
    }
    defer cursor.Close(ctx)

    var posts []types.Story
    if err = cursor.All(ctx, &posts); err != nil {
        return nil, fmt.Errorf("cursor.All error: %w", err)
    }
    return posts, nil
}

// IncrementCommentCount updates comment count and adds comment ID
func (s *StoryStore) AddComment(ctx context.Context, postID string, commentID string) error {
    // add comment id to replies and increment reply count by 1
	update := bson.M{
		"$push": bson.M{"replies": commentID},
		"$inc": bson.M{"reply_count": 1},
	}
	_, err := s.stories.UpdateOne(ctx, bson.M{"_id": postID}, update)

	return err
}

// TODO: create a helper to use it between story and comment
func (s *StoryStore) ToggleUpvote(ctx context.Context, postID string, userID string) (types.Story, error) {
    // For first try, get the updated doc
    var updatedStory types.Story

    // Try to add upvote
    err := s.stories.FindOneAndUpdate(
        ctx,
        bson.M{"_id": postID, "upvoted_by": bson.M{"$ne": userID}},
        bson.M{
            "$inc": bson.M{"upvote_count": 1},
            "$addToSet": bson.M{"upvoted_by": userID},
        },
        options.FindOneAndUpdate().SetReturnDocument(options.After),
    ).Decode(&updatedStory)

    if err == mongo.ErrNoDocuments {
        // If user already upvoted, remove the upvote
        err = s.stories.FindOneAndUpdate(
            ctx,
            bson.M{"_id": postID},
            bson.M{
                "$inc": bson.M{"upvote_count": -1},
                "$pull": bson.M{"upvoted_by": userID},
            },
            options.FindOneAndUpdate().SetReturnDocument(options.After),
        ).Decode(&updatedStory)
        return updatedStory, err
    }

    return updatedStory, err
}
