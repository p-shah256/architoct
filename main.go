package main

import (
	"html/template"
	"log"
	"log/slog"
	"net/http"
)

type Story struct {
	ID           int
	Title        string
	VoteCount    int
	User         string
	TimeAgo      string
	CommentCount int
}

type HomepageData struct {
	Stories []Story
}

type StoryPageData struct {
	Story *Story // Note this is a pointer to a single Story
}

func main() {
	// Parse templates
	tmpl, err := template.ParseFiles(
		"views/index.html",
		"views/storyPage.html",
		"views/partials/story.html",
	)
	if err != nil {
		log.Fatal(err)
	}

	// Sample stories data
	stories := []Story{
		{
			ID:           1,
			Title:        "We're (finally) going to the cloud",
			VoteCount:    42,
			User:         "tonyfinn",
			TimeAgo:      "5 hours ago",
			CommentCount: 15,
		},
		{
			ID:           2,
			Title:        "The Future of Microservices",
			VoteCount:    38,
			User:         "sarah_dev",
			TimeAgo:      "3 hours ago",
			CommentCount: 23,
		},
	}

	// Handle root route
	http.HandleFunc("/", func(writer http.ResponseWriter, r *http.Request) {
		data := HomepageData{
			Stories: stories,
		}
		tmpl.ExecuteTemplate(writer, "index.html", data)
	})

	// Handle story route
	http.HandleFunc("/story/", func(writer http.ResponseWriter, r *http.Request) {

		slog.Info("request recieved at /story/")

		story := &stories[0]
		data := StoryPageData{
			Story: story,
		}
		err := tmpl.ExecuteTemplate(writer, "storyPage.html", data)
		if err != nil {
			log.Printf("Template error: %v", err)
			http.Error(writer, "Internal Server Error", http.StatusInternalServerError)
		}

		slog.Info("sent.. EOF /story/")
	})

	log.Println("Server starting on :8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
