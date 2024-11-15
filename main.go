package main

import (
    "html/template"
    "log"
    "net/http"
)

type Story struct {
    Title       string
    VoteCount   int
    User        string
    TimeAgo     string
    CommentCount int
}

func main() {
    // Parse template
    tmpl, err := template.ParseFiles("templates/index.html")
    if err != nil {
        log.Fatal(err)
    }

    // Handle root route
    http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
        story := Story{
            Title:       "We're (finally) going to the cloud",
            VoteCount:   42,
            User:        "tonyfinn",
            TimeAgo:     "5 hours ago",
            CommentCount: 15,
        }

        tmpl.Execute(w, story)
    })

    log.Println("Server starting on :8080")
    log.Fatal(http.ListenAndServe(":8080", nil))
}
