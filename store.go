// NOTE: can be implemented by many diff types of stores... mongo atm... sql later
package main

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
