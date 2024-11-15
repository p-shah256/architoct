// store.go
package main

// NewStoryStore creates a new story store with sample data.
func NewStoryStore() *StoryStore {
	return &StoryStore{
		stories: []Story{
			{
				ID:           1,
				Title:        "The Future of Microservices",
				VoteCount:    38,
				User:         "sarah_dev",
				TimeAgo:      "3 hours ago",
				CommentCount: 3,
				Body: "Microservices have been a game-changer in the world of software development, offering flexibility, scalability, and better fault tolerance. As more organizations adopt microservice architectures, it's important to consider the future trends in this space. The growing adoption of containerization and Kubernetes will continue to shape the landscape, allowing for easier management of distributed systems. Additionally, the integration of AI and machine learning with microservices can enable smarter, more adaptive applications. However, as complexity increases, monitoring and managing these systems will become even more critical. It's essential to build the right tools and practices to address these challenges.",
				Comments: []Comment{
					{
						ID:        1,
						Content:   "Microservices are great, but they come with a lot of complexity. Microservices are great, but they come with a lot of complexity.",
						User:      "tech_guru",
						TimeAgo:   "2 hours ago",
						VoteCount: 5,
						Replies: []Comment{
							{
								ID:        2,
								Content:   "True, but tools like Kubernetes are making it easier.",
								User:      "cloud_master",
								TimeAgo:   "1 hour ago",
								VoteCount: 3,
							},
						},
					},
					{
						ID:        3,
						Content:   "Monoliths are still valid for smaller teams!",
						User:      "old_school_dev",
						TimeAgo:   "1 hour ago",
						VoteCount: 8,
					},
				},
			},
			{
				ID:           2,
				Title:        "Why Go is Great for Web Development",
				VoteCount:    42,
				User:         "gopher_pro",
				TimeAgo:      "5 hours ago",
				CommentCount: 2,
				Comments: []Comment{
					{
						ID:        4,
						Content:   "Go's simplicity is what makes it shine for web apps.",
						User:      "simplicity_fan",
						TimeAgo:   "4 hours ago",
						VoteCount: 10,
					},
					{
						ID:        5,
						Content:   "Agreed! The standard library is a huge plus.",
						User:      "lib_lover",
						TimeAgo:   "3 hours ago",
						VoteCount: 7,
					},
				},
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
