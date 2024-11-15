package main

import (
	"html/template"
)

// Templates handles rendering of HTML templates
type Templates struct {
	templates *template.Template
}

// App holds the application state
type App struct {
	store     *StoryStore
	templates *Templates
}

// PageData wraps our data for the template
type PageData struct {
	Stories []Story
}

// Story represents a single story with comments.
type Story struct {
	ID           int
	Title        string
	VoteCount    int
	User         string
	TimeAgo      string
	Body		 string
	CommentCount int
	Comments     []Comment
}

// Comment represents a single comment on a story.
type Comment struct {
	ID      int
	Content string
	User    string
	TimeAgo string
	VoteCount int
	Replies []Comment
}

// StoryStore holds a collection of stories.
type StoryStore struct {
	stories []Story
}
