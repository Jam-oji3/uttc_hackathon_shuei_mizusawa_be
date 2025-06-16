package model

import "time"

type Post struct {
	Id        string    `json:"id"`
	UserId    string    `json:"userId"`
	Text      string    `json:"text"`
	ReplyTo   *string   `json:"replyTo"`
	RepostRef *string   `json:"repostRef"`
	MediaType *string   `json:"mediaType"`
	MediaURL  *string   `json:"mediaUrl"`
	CreatedAt time.Time `json:"createdAt"`
}
type PostWithUserAndCounts struct {
	Post        Post
	User        User
	LikeCount   int
	RepostCount int
}
