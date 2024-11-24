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
	CreatedAt time.Time `bson:"created_at,omitempty"`
	LastLogin time.Time `bson:"last_login,omitempty"`
}

// // comment /////////////////////////////////////////////////////////////////
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

	// extra meta data user specific
	HasUpvoted bool `bson:"-"`
}

func (s *Story) SetUserSpecificData(userID string) {
    for _, upvoterID := range s.UpvotedBy {
        if upvoterID == userID {
            s.HasUpvoted = true
            return
        }
    }
    s.HasUpvoted = false
}

// // comment /////////////////////////////////////////////////////////////////
// Comment represents a comment on a story
type Comment struct {
	// Comment Data
	PostID    string    `bson:"post_id"`
	UserID    string    `bson:"user_id"`
	Body      string    `bson:"body"`
	CreatedAt time.Time `bson:"created_at"`
	TimeAgo   string

    ID        string `bson:"_id,omitempty"`  // Changed type and fixed tag
	IsDeleted bool      `bson:"is_deleted,omitempty"`

	// Metadata
	UpvoteCount int      `bson:"upvote_count,omitempty"`
	UpvotedBy   []string `bson:"upvoted_by,omitempty"` // user_ids
	ReplyCount  int      `bson:"reply_count,omitempty"`
	Replies     []string `bson:"replies,omitempty"` // comment_ids

	// extra meta data user specific
	HasUpvoted bool `bson:"-"`
}

func (c *Comment) SetUserSpecificData(userID string) {
    for _, upvoterID := range c.UpvotedBy {
        if upvoterID == userID {
            c.HasUpvoted = true
            return
        }
    }
    c.HasUpvoted = false
}

//SERVICE TYPES/////////////////////////////////////////////////////////////
type StoryPage struct {
	Story Story
	Comments []Comment
}
