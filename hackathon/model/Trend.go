package model

import "time"

type Trend struct {
	Id    string    `json:"id"`
	Word  string    `json:"word"`
	Count int       `json:"count"`
	Hour  time.Time `json:"hour"`
}

type TrendSummary struct {
	Word  string `json:"word"`
	Count int    `json:"count"`
}
