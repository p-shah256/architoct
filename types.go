package main

import (
	"html/template"
)

// Story represents a single story item
type Story struct {
	ID           int    `json:"id"`
	Title        string `json:"title"`
	VoteCount    int    `json:"vote_count"`
	User         string `json:"user"`
	TimeAgo      string `json:"time_ago"`
	CommentCount int    `json:"comment_count"`
}

// StoryStore manages the in-memory storage of stories
type StoryStore struct {
	stories []Story
}

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
