package service

// this layer does not care if the handler was HTMX or REACT or just a cli..
// it does not care about rendering or templates or anything
//
// this is the real buisness logic

import (
	"architoct/internal/store/mongos"
	"architoct/internal/types"
	"context"
)

type ArchitoctService struct {
    storyStore *mongos.StoryStore
    commentStore *mongos.CommentStore
    userStore *mongos.UserStore
}

func (architcot *ArchitoctService) HomeFeed(ctx context.Context) ([]*types.Story, error) {
	// 1. get all the top 20 stories for now
	stories, err :=architcot.storyStore.GetRecent(ctx, 20, "week");
	if err != nil {
		return nil, err
	}
	return stories, nil
}
