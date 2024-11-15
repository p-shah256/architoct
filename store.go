// store.go
package main

// NewStoryStore creates a new story store with sample data
func NewStoryStore() *StoryStore {
	return &StoryStore{
		stories: []Story{
			{
				ID:           1,
				Title:        "The Future of Microservices",
				VoteCount:    38,
				User:         "sarah_dev",
				TimeAgo:      "3 hours ago",
				CommentCount: 23,
			},
			{
				ID:           2,
				Title:        "Why Go is Great for Web Development",
				VoteCount:    42,
				User:         "gopher_pro",
				TimeAgo:      "5 hours ago",
				CommentCount: 15,
			},
{
				ID:           1,
				Title:        "The Future of Microservices",
				VoteCount:    38,
				User:         "sarah_dev",
				TimeAgo:      "3 hours ago",
				CommentCount: 23,
			},
			{
				ID:           2,
				Title:        "Why Go is Great for Web Development",
				VoteCount:    42,
				User:         "gopher_pro",
				TimeAgo:      "5 hours ago",
				CommentCount: 15,
			},
		},
	}
}

// GetStory returns a story by ID
func (s *StoryStore) GetStory(id int) (Story, bool) {
	for _, story := range s.stories {
		if story.ID == id {
			return story, true
		}
	}
	return Story{}, false
}

// GetAllStories returns all stories
func (s *StoryStore) GetAllStories() []Story {
	return s.stories
}

// IncrementVote increments the vote count for a story
func (s *StoryStore) IncrementVote(id int) bool {
	for i := range s.stories {
		if s.stories[i].ID == id {
			s.stories[i].VoteCount++
			return true
		}
	}
	return false
}
