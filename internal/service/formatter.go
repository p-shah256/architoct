package service

import (
	"architoct/internal/types"
	"fmt"
	"time"
)

// Helper function in the service package
func formatTimeAgo(t time.Time) string {
	now := time.Now()
	diff := now.Sub(t)

	switch {
	case diff < time.Minute:
		return "just now"
	case diff < time.Hour:
		mins := int(diff.Minutes())
		return fmt.Sprintf("%dm ago", mins)
	case diff < 24*time.Hour:
		hours := int(diff.Hours())
		return fmt.Sprintf("%dh ago", hours)
	case diff < 7*24*time.Hour:
		days := int(diff.Hours() / 24)
		return fmt.Sprintf("%dd ago", days)
	default:
		weeks := int(diff.Hours() / 24 / 7)
		return fmt.Sprintf("%dw ago", weeks)
	}
}

// handles all the required formatting from db to service layer
// eg:
// created at -> time ago
func formatStory(story *types.Story) {
	story.TimeAgo = formatTimeAgo(story.CreatedAt);
}

func formatComment(comment *types.Comment) {
	comment.TimeAgo = formatTimeAgo(comment.CreatedAt);
}
