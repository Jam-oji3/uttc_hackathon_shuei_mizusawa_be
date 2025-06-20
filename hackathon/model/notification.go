package model

type Notification struct {
	Type  string `json:"type"` //like or repost or follow
	Actor struct {
		Username string `json:"username"`
		IconUrl  string `json:"iconUrl"`
	} `json:"actor"`
	TargetId  *string `json:"targetId"`
	CreatedAt string  `json:"createdAt"`
}
