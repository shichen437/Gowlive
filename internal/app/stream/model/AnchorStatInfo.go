package model

type AnchorStatInfo struct {
	WeekFollowersIncr  int               `json:"weekFollowersIncr"`
	WeekLikeNumIncr    int               `json:"weekLikeNumIncr"`
	MonthFollowersIncr int               `json:"monthFollowersIncr"`
	HistoryData        []*AnchorStatData `json:"historyData"`
}

type AnchorStatData struct {
	RecordDate string `json:"recordDate"`
	Followers  int    `json:"followers"`
	LikeCount  int    `json:"likeCount"`
}
