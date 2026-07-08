package model

import "time"

type Follow struct {
	Id         string    `json:"id"`
	FollowerId string    `json:"follower_id"`
	FollowedId string    `json:"followed_id"`
	CreatedAt  time.Time `json:"created_at"`
}
