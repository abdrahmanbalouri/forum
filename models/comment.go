package models

import "time"

type Comment struct {
	ID           int       `json:"id"`
	PostID       int       `json:"post_id"`
	UserID       int       `json:"user_id"`
	Username     string    `json:"username"`
	Content      string    `json:"content"`
	CreatedAt    time.Time `json:"created_at"`
	LikeCount    int       `json:"like_count"`
	DislikeCount int       `json:"dislike_count"`
	UserReaction string    `json:"user_reaction,omitempty"`
}
