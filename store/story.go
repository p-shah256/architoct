// internal/store/post.go
// no abstraction for mvp -- eventually we could have a port called port store and then let 1. mongo implment it and then 2. dynamo and so on...
// abstract away *mongo.Collection

package store

import (
	"context"
	"time"

	"architoct/types"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// TYPE AND CREATE /////////////////////////////////////////////////////////////
type StoryStore struct {
	stories *mongo.Collection
}

func NewPostStore(db *mongo.Database) *StoryStore {
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
	var post types.Story
	err := s.stories.FindOne(ctx, bson.M{"_id": id}).Decode(&post)
	return &post, err
}

// GetRecent supports our homepage feed
// limit: how many posts to fetch
// timeRange: "hour", "day", "week" etc
func (s *StoryStore) GetRecent(ctx context.Context, limit int64, timeRange string) ([]*types.Story, error) {
	// Calculate time threshold based on range
	threshold := time.Now()
	switch timeRange {
	case "hour":
		threshold = threshold.Add(-1 * time.Hour)
	case "day":
		threshold = threshold.Add(-24 * time.Hour)
	case "week":
		threshold = threshold.Add(-7 * 24 * time.Hour)
	}

	// Find posts with sorting
	opts := options.Find().
		SetSort(bson.D{{Key: "created_at", Value: -1}}).
		SetSort(bson.D{{Key: "upvotes_count", Value: -1}}).
		SetLimit(limit)

	cursor, err := s.stories.Find(ctx,
		bson.M{"created_at": bson.M{"$gte": threshold}},
		opts,
	)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var posts []*types.Story
	if err = cursor.All(ctx, &posts); err != nil {
		return nil, err
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

func (s *StoryStore) ToggleUpvote(ctx context.Context, postID string, userID string) error {
    // First try to add upvote
    result := s.stories.FindOneAndUpdate(
        ctx,
        bson.M{"_id": postID, "upvoted_by": bson.M{"$ne": userID}},
        bson.M{
            "$inc": bson.M{"upvote_count": 1},
            "$addToSet": bson.M{"upvoted_by": userID},
        },
    )

    if result.Err() == mongo.ErrNoDocuments {
        // If user already upvoted, remove the upvote
        _, err := s.stories.UpdateOne(
            ctx,
            bson.M{"_id": postID},
            bson.M{
                "$inc": bson.M{"upvote_count": -1},
                "$pull": bson.M{"upvoted_by": userID},
            },
        )
        return err
    }

    return result.Err()
}
