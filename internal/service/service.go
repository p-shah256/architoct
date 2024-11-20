package service

// DESIGN:
// this layer does not care if the handler was HTMX or REACT or just a cli..
// it does not care about rendering or templates or anything
// also it does not care if its mongo or something else.. so it has to be passed into it
//
// this is the real buisness logic
//
// also this cannot return any ptrs cause HTMX needs actual data

import (
	"architoct/internal/store/mongos"
	"architoct/internal/types"
	"context"
	"log/slog"
	"time"
	// "log/slog"
)

//ENUM TO KEEP things simple///////////////////////////////////////////////////
type ContentType string

const (
    TypeStory   ContentType = "STORY"
    Comment ContentType = "COMMENT"
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

// GET //////////////////////////////////////////////////////////////////////
// TODO: check types and remove redundant returns and use ptrs
func (architcot *ArchitoctService) HomeFeed(ctx context.Context) ([]types.Story, error) {
	stories, err := architcot.storyStore.GetRecent(ctx, 20)
	if err != nil {
		return nil, err
	}

	for i := range stories {
		formatStory(&stories[i])
	}
	return stories, nil
}

func (architcot *ArchitoctService) StoryPage(ctx context.Context, id string) (types.StoryPage, error) {
	requestedStory, err := architcot.storyStore.GetByID(ctx, id)
	if err != nil {
		return types.StoryPage{}, err
	}

	comments := make([]types.Comment, 0, len(requestedStory.Replies))
	// not using reply count here cause it might be incosistent failing whole fetch
	for i := range len(requestedStory.Replies) {
		comment, err := architcot.commentStore.GetById(ctx, requestedStory.Replies[i])
		if err != nil {
			return types.StoryPage{}, err
		}
		formatComment(comment)
		comments = append(comments, *comment)
	}

	// uses pointer so no return needed
	formatStory(requestedStory)
	storyPage := types.StoryPage {
		Story: *requestedStory,
		Comments: comments,
	}

	return storyPage, nil
}


// POST ////////////////////////////////////////////////////////////////////////
func (architcot *ArchitoctService) Upvote(ctx context.Context, contentType ContentType, id string, userid string) (any, error) {
	slog.Info("upvoting...", "comment", contentType, "id", id)
	if contentType == Comment {
		updatedStory, err := architcot.commentStore.ToggleUpvote(ctx, id, userid)
		return updatedStory, err
	} else {
		updatedComment, err := architcot.storyStore.ToggleUpvote(ctx, id, userid)
		return updatedComment, err
	}
}


func (architcot *ArchitoctService) Comment(ctx context.Context, parentid string, userid string, body string, contentType ContentType) (error) {
	if contentType == Comment {

	} else {
		slog.Info("Commenting on a story", )
		commentID, err := architcot.commentStore.Create(ctx, &types.Comment{
			PostID: parentid,
			Body: body,
			UserID: userid,
			CreatedAt: time.Now(),
		})
		if err != nil {
			return err
		}
		err = architcot.storyStore.AddComment(ctx, parentid, commentID)
		if err != nil {
			return err
		}
	}
	return nil
}
