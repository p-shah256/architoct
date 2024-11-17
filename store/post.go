// internal/store/post.go
// no abstraction for mvp -- eventually we could have a port called port store and then let 1. mongo implment it and then 2. dynamo and so on...

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
type PostStore struct {
    posts *mongo.Collection
}

func NewPostStore(db *mongo.Database) *PostStore {
    return &PostStore{
        posts: db.Collection("posts"),
    }
}

// OPERATIONS /////////////////////////////////////////////////////////////////
// Core operations that we'll need for MVP
func (s *PostStore) Create(ctx context.Context, post *types.Post) error {
    // For MVP, we'll do simple inserts without transactions
    _, err := s.posts.InsertOne(ctx, post)
    return err
}

func (s *PostStore) GetByID(ctx context.Context, id string) (*types.Post, error) {
    var post types.Post
    err := s.posts.FindOne(ctx, bson.M{"_id": id}).Decode(&post)
    return &post, err
}

// GetRecent supports our homepage feed
// limit: how many posts to fetch
// timeRange: "hour", "day", "week" etc
func (s *PostStore) GetRecent(ctx context.Context, limit int64, timeRange string) ([]*types.Post, error) {
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
        SetLimit(limit)

    cursor, err := s.posts.Find(ctx,
        bson.M{"created_at": bson.M{"$gte": threshold}},
        opts,
    )
    if err != nil {
        return nil, err
    }
    defer cursor.Close(ctx)

    var posts []*types.Post
    if err = cursor.All(ctx, &posts); err != nil {
        return nil, err
    }
    return posts, nil
}

// IncrementCommentCount updates comment count and adds comment ID
func (s *PostStore) IncrementCommentCount(ctx context.Context, postID string, commentID string) error {
    update := bson.M{
        "$inc": bson.M{"comment_count": 1},
        "$push": bson.M{"comment_ids": commentID},
    }
    _, err := s.posts.UpdateOne(ctx, bson.M{"_id": postID}, update)
    return err
}

// For upvotes we'll keep it simple
func (s *PostStore) IncrementUpvoteCount(ctx context.Context, postID string) error {
    update := bson.M{"$inc": bson.M{"upvote_count": 1}}
    _, err := s.posts.UpdateOne(ctx, bson.M{"_id": postID}, update)
    return err
}
