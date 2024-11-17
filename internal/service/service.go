package service

// DESIGN:
// this layer does not care if the handler was HTMX or REACT or just a cli..
// it does not care about rendering or templates or anything
// also it does not care if its mongo or something else.. so it has to be passed into it
//
// this is the real buisness logic

import (
	"architoct/internal/store/mongos"
	"architoct/internal/types"
	"context"
	"fmt"
	"log/slog"

	// "log/slog"
	"time"
)

type ArchitoctService struct {
	storyStore   *mongos.StoryStore
	commentStore *mongos.CommentStore
	userStore    *mongos.UserStore
}

func NewArchitoctService(s *mongos.StoryStore, c *mongos.CommentStore, u *mongos.UserStore) *ArchitoctService {
	return &ArchitoctService{
		storyStore:   s,
		commentStore: c,
		userStore:    u,
	}
}

func (architcot *ArchitoctService) HomeFeed(ctx context.Context) ([]types.Story, error) {
	// 1. get all the top 20 stories for now
	stories, err := architcot.storyStore.GetRecent(ctx, 20)
	if err != nil {
		return nil, err
	}

	// 2. Transform the timestamps
	for i := range stories {
		stories[i].TimeAgo = formatTimeAgo(stories[i].CreatedAt)
	}
	return stories, nil
}

// PERF: not sure if this is a bad idea to send entire story
func (architcot *ArchitoctService) Upvote(ctx context.Context, comment bool, id string) (any, error) {
	slog.Info("upvoting...", "comment", comment, "id", id)
	if comment {
		updatedStory, err := architcot.commentStore.ToggleUpvote(ctx, id, "asdb")
		return updatedStory, err
	} else {
		updatedComment, err := architcot.storyStore.ToggleUpvote(ctx, id, "asdb")
		return updatedComment, err
	}

}

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
