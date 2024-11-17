package types

import (
	"html/template"
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

// Templates handles rendering of HTML templates
type Templates struct {
	templates *template.Template
}

// NOTE: maybe external?
// DB TYPES/////////////////////////////////////////////////////////////////////
type Story struct {
	ID           string    `bson:"_id"`
	UserID       string    `bson:"user_id"`
	Title        string    `bson:"title"`
	Body         string    `bson:"body"`
	CreatedAt    time.Time `bson:"created_at"`
	UpvoteCount  int32     `bson:"upvote_count"`
	CommentCount int32     `bson:"comment_count"`
	CommentIDs   []string  `bson:"comment_ids,omitempty"`
}

// User represents a user identified by browser fingerprint
// TODO: look into how this can be extended with auth
type User struct {
	// ID          primitive.ObjectID `bson:"_id,omitempty"`
	Fingerprint string             `bson:"fingerprint"`
	CreatedAt   time.Time          `bson:"created_at"`
	LastLogin   time.Time          `bson:"last_login"`
}

// Post represents a forum post
type Post struct {
	ID           primitive.ObjectID `bson:"_id,omitempty"`
	UserID       string             `bson:"user_id"`
	Title        string             `bson:"title"`
	Body         string             `bson:"body"`
	CreatedAt    time.Time          `bson:"created_at"`
	UpvoteCount  int32              `bson:"upvote_count"`
	CommentCount int32              `bson:"comment_count"`
	CommentIDs   []string           `bson:"comment_ids,omitempty"`
}

type ParentType string

const (
	ParentTypePost    ParentType = "post"
	ParentTypeComment ParentType = "comment"
)

// Upvote represents an upvote on a post or comment
type Upvote struct {
	ID         primitive.ObjectID `bson:"_id,omitempty"`
	ParentID   string             `bson:"parent_id"`
	ParentType ParentType         `bson:"parent_type"`
	UserID     string             `bson:"user_id"`
	CreatedAt  time.Time          `bson:"created_at"`
}

type Comment struct {
	ID          primitive.ObjectID `bson:"_id,omitempty" json:"id"`
	PostID      string             `bson:"post_id" json:"post_id"`
	UserID      string             `bson:"user_id" json:"user_id"`
	ParentID    *string            `bson:"parent_id,omitempty" json:"parent_id,omitempty"`
	Body        string             `bson:"body" json:"body"`
	CreatedAt   time.Time          `bson:"created_at" json:"created_at"`
	UpvoteCount int32              `bson:"upvote_count" json:"upvote_count"`
	IsDeleted   bool               `bson:"is_deleted" json:"is_deleted"`
}
