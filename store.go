// store.go
package main

// NewStoryStore creates a new story store with sample data.
func NewStoryStore() *StoryStore {
	return &StoryStore{
		stories: []Story{
			{
				ID:           1,
				Title:        "The Future of Smart Cities",
				VoteCount:    38,
				User:         "architect_joe",
				TimeAgo:      "3 hours ago",
				CommentCount: 3,
				Body: "Smart cities are the future of urban development, combining technology and sustainability to create more livable environments. With advancements in IoT (Internet of Things), cities can become more efficient, with better traffic management, waste disposal, and energy usage. As the adoption of AI-driven urban planning continues, we can expect to see more data-driven decisions that can improve public services. However, the challenge lies in ensuring privacy and security while implementing these technologies. Moreover, integrating green spaces and ensuring that the infrastructure is adaptable to climate change will be critical for long-term success.",
				Comments: []Comment{
					{
						ID:        1,
						Content:   "Smart cities are great, but the integration of technology is challenging, especially in older buildings.",
						User:      "urban_planner",
						TimeAgo:   "2 hours ago",
						VoteCount: 5,
						Replies: []Comment{
							{
								ID:        2,
								Content:   "True, but retrofitting with newer tech solutions is becoming more cost-effective.",
								User:      "tech_builder",
								TimeAgo:   "1 hour ago",
								VoteCount: 3,
							},
						},
					},
					{
						ID:        3,
						Content:   "Sustainability in architecture is key for future cities. Energy-efficient buildings should be prioritized.",
						User:      "eco_architect",
						TimeAgo:   "1 hour ago",
						VoteCount: 8,
					},
				},
			},
			{
				ID:           2,
				Title:        "The Rise of Modular Construction and its rigidities are on the rise too",
				VoteCount:    42,
				User:         "modular_mike",
				TimeAgo:      "5 hours ago",
				CommentCount: 2,
				Comments: []Comment{
					{
						ID:        4,
						Content:   "Modular construction can significantly speed up the building process and reduce costs.",
						User:      "cost_saver",
						TimeAgo:   "4 hours ago",
						VoteCount: 10,
					},
					{
						ID:        5,
						Content:   "Agreed! The ability to assemble and disassemble parts easily is a game-changer.",
						User:      "builder_bob",
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
