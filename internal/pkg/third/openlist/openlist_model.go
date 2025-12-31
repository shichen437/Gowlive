package openlist

var (
	loginUrl            = "/api/auth/login"
	userInfoUrl         = "/api/me"
	dirInfoUrl          = "/api/fs/get"
	uploadFileStreamUrl = "/api/fs/put"
	mkdirUrl            = "/api/fs/mkdir"
	logoutUrl           = "/api/auth/logout"
)

type CommonResp struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
}

type LoginRespModel struct {
	CommonResp
	Data struct {
		Token string `json:"token"`
	} `json:"data"`
}

type UserInfoRespModel struct {
	CommonResp
	Data struct {
		ID         string `json:"id"`
		Username   string `json:"username"`
		BasePath   string `json:"base_path"`
		Disabled   bool   `json:"disabled"`
		Permission int    `json:"permission"`
		Otp        bool   `json:"otp"`
	} `json:"data"`
}

type DirInfoRespModel struct {
	CommonResp
	Data *DirContentModel `json:"data"`
}

type DirContentModel struct {
	ID       string `json:"id"`
	Path     string `json:"path"`
	Name     string `json:"name"`
	Size     int64  `json:"size"`
	IsDir    bool   `json:"is_dir"`
	Modified string `json:"modified"`
	Created  string `json:"created"`
	Sign     string `json:"sign"`
	Thumb    string `json:"thumb"`
	Type     int    `json:"type"`
}
