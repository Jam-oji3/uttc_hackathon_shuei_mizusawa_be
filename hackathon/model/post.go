package model

import (
	"time"
)

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
	Id        string    `json:"id"`
	Text      string    `json:"text"`
	ReplyTo   *string   `json:"replyTo"`
	RepostRef *string   `json:"repostRef"`
	MediaType *string   `json:"mediaType"`
	MediaURL  *string   `json:"mediaUrl"`
	CreatedAt time.Time `json:"createdAt"`
	Author    struct {
		Id          string `json:"id"`
		Username    string `json:"username"`
		DisplayName string `json:"displayName"`
		IconURL     string `json:"iconUrl"`
	} `json:"author"`
	Stats struct {
		Likes    int `json:"likes"`
		Reposts  int `json:"reposts"`
		Comments int `json:"comments"`
	} `json:"stats"`
	UserActions struct {
		Liked    bool `json:"liked"`
		Reposted bool `json:"reposted"`
	} `json:"userActions"`
}
