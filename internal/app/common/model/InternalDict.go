package model

type InternalDict struct {
	DictType  string `json:"dictType"`
	DictLabel string `json:"dictLabel"`
	DictValue string `json:"dictValue"`
	DictSort  int    `json:"dictSort"`
}

var (
	LivePlatform = "live_platform"
	ChannelType  = "channel_type"
	m            = make(map[string][]InternalDict)
)

func init() {
	addDictData(LivePlatform, []InternalDict{
		{
			DictType:  LivePlatform,
			DictLabel: "抖音",
			DictValue: "douyin",
			DictSort:  1,
		},
		{
			DictType:  LivePlatform,
			DictLabel: "哔哩哔哩",
			DictValue: "bilibili",
			DictSort:  2,
		},
		{
			DictType:  LivePlatform,
			DictLabel: "YY",
			DictValue: "yy",
			DictSort:  3,
		},
		{
			DictType:  LivePlatform,
			DictLabel: "Bigo",
			DictValue: "bigo",
			DictSort:  4,
		},
	})
	addDictData(ChannelType, []InternalDict{
		{
			DictType:  ChannelType,
			DictLabel: "邮箱",
			DictValue: "email",
			DictSort:  1,
		},
		{
			DictType:  ChannelType,
			DictLabel: "Gotify",
			DictValue: "gotify",
			DictSort:  2,
		},
		{
			DictType:  ChannelType,
			DictLabel: "飞书",
			DictValue: "lark",
			DictSort:  3,
		},
		{
			DictType:  ChannelType,
			DictLabel: "钉钉",
			DictValue: "dingTalk",
			DictSort:  4,
		},
		{
			DictType:  ChannelType,
			DictLabel: "企业微信",
			DictValue: "weCom",
			DictSort:  5,
		},
	})
}

func GetDictDataByType(dictType string) *[]InternalDict {
	data, ok := m[dictType]
	if !ok {
		return &[]InternalDict{}
	}
	return &data
}

func addDictData(dictType string, dictData []InternalDict) {
	data, ok := m[dictType]
	if !ok {
		data = []InternalDict{}
	}
	m[dictType] = append(data, dictData...)
}
