package types

import (
	"time"
)

// NOTE: maybe external?
// DB TYPES/////////////////////////////////////////////////////////////////////

// User represents a user identified by browser fingerprint
// TODO: look into how this can be extended with auth
// User represents a user in the system
type User struct {
	ID        string    `bson:"_id"`
	CreatedAt time.Time `bson:"created_at"`
	LastLogin time.Time `bson:"last_login,omitempty"`
}

// Story represents a story/post in the system
type Story struct {
	// Story Components
	ID        string    `bson:"_id"`
	UserID    string    `bson:"user_id"`
	Title     string    `bson:"title"`
	Body      string    `bson:"body"`
	CreatedAt time.Time `bson:"created_at"`
	TimeAgo   string

	// Metadata
	UpvoteCount int      `bson:"upvote_count,omitempty"`
	UpvotedBy   []string `bson:"upvoted_by,omitempty"` // user_ids
	ReplyCount  int      `bson:"reply_count,omitempty"`
	Replies     []string `bson:"replies,omitempty"` // comment_ids
}

// Comment represents a comment on a story
type Comment struct {
	// Comment Data
	ID        string    `bson:"_id"`
	PostID    string    `bson:"post_id"`
	UserID    string    `bson:"user_id"`
	Body      string    `bson:"body"`
	CreatedAt time.Time `bson:"created_at"`
	IsDeleted bool      `bson:"is_deleted,omitempty"`

	// Metadata
	UpvoteCount int      `bson:"upvote_count,omitempty"`
	UpvotedBy   []string `bson:"upvoted_by,omitempty"` // user_ids
	ReplyCount  int      `bson:"reply_count,omitempty"`
	Replies     []string `bson:"replies,omitempty"` // comment_ids
}
