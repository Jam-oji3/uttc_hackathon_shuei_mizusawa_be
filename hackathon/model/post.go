package model

import "time"

type Post struct {
	Id        string
	UserId    string
	Text      string
	ReplyTo   *string
	RepostRef *string
	MediaType *string
	MediaURL  *string
	CreatedAt time.Time
}
type PostWithUserAndCounts struct {
	Post        Post
	User        User
	LikeCount   int
	RepostCount int
}
