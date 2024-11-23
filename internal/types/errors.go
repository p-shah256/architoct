package types

import "errors"

var (
    ErrUsernameTaken = errors.New("username already taken")
    ErrUserNotFound  = errors.New("user not found")
    ErrStoryNotFound  = errors.New("story not found")
    ErrCommentNotFound  = errors.New("Comment not found")
    ErrCommentNotPosted  = errors.New("Comment not posted")
)
