package model

type UserInfo struct {
	Id       int    `json:"id"`
	Username string `json:"username"`
	Nickname string `json:"nickname"`
	Sex      int    `json:"sex"`
	Status   int    `json:"status"`
}
