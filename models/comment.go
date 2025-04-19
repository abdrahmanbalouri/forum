package models

type Comment struct {
	ID                   int
	PostID               int
	UserID               int
	Username             string
	Content              string
	CreatedAt            string
	LikeCount            int
	DislikeCount         int
	UserReaction         string
	MinutesSinceCreation int
}
