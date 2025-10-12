package anchor

type AnchorInfo struct {
	Platform       string `json:"platform"`
	AnchorName     string `json:"anchor_name"`
	UniqueId       string `json:"unique_id"`
	Signature      string `json:"signature"`
	FollowerCount  int    `json:"follower_count"`
	FollowingCount int    `json:"following_count"`
	LikeCount      int    `json:"like_count"`
	VideoCount     int    `json:"video_count"`
}
