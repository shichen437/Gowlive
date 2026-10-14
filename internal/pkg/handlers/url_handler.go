package handlers

import (
	"context"
	"encoding/json"
	"net/url"
	"strings"

	"github.com/gogf/gf/v2/errors/gerror"
	"github.com/gogf/gf/v2/frame/g"
	"github.com/shichen437/gowlive/internal/pkg/consts"
	"github.com/shichen437/gowlive/internal/pkg/utils"
)

var (
	douyin_live   = "https://live.douyin.com/"
	douyin_home   = "www.douyin.com"
	douyin_share  = "www.iesdouyin.com"
	douyin_mobile = "v.douyin.com"
	user_info_url = "https://www.douyin.com/web/api/v2/user/info/?sec_uid="
)

type TempAnchorInfoResp struct {
	UserInfo struct {
		UniqueID string `json:"unique_id"`
	} `json:"user_info"`
}

func HandleRedirectUrl(ctx context.Context, sUrl string) string {
	u, err := url.Parse(sUrl)
	if err != nil {
		return sUrl
	}
	hostname := strings.ToLower(u.Hostname())
	if hostname == douyin_home || hostname == douyin_mobile || hostname == douyin_share {
		realUrl, err := handleDouyinUrl(ctx, u, hostname)
		if err != nil {
			g.Log().Error(ctx, "handleDouyinUrl err:", err)
			return sUrl
		}
		return realUrl
	}
	return sUrl
}

func handleDouyinUrl(ctx context.Context, sUrl *url.URL, hostname string) (realUrl string, err error) {
	realUrl = ""
	if hostname == douyin_mobile {
		mUrl, err := utils.GetRedirectURL(ctx, sUrl.String())
		if err != nil {
			return "", gerror.Newf("获取重定向链接失败: %s", err)
		}
		mobileU, err := url.Parse(mUrl)
		if err != nil {
			return "", gerror.Newf("解析重定向链接失败: %s", err)
		}
		secUserId := ""
		if strings.ToLower(mobileU.Hostname()) == douyin_share {
			secUserId = extractSecUserIdFromPath(mobileU.Path)
		} else {
			secUserId = mobileU.Query().Get("sec_user_id")
		}
		if secUserId == "" {
			return "", gerror.New("获取用户唯一 ID 失败")
		}
		realUrl = "https://www.douyin.com/user/" + secUserId
	}
	if realUrl == "" {
		realUrl = sUrl.String()
	}
	return getDouyinLiveUrl(ctx, realUrl)
}

func getDouyinLiveUrl(ctx context.Context, sUrl string) (realUrl string, err error) {
	hUrl, err := url.Parse(sUrl)
	if err != nil {
		return "", gerror.Newf("解析真实链接失败: %s", err)
	}
	if !strings.Contains(hUrl.String(), "/user/") {
		return "", gerror.New("抖音用户主页 URL 格式错误")
	}
	secId := strings.TrimPrefix(utils.FindFirstMatch(hUrl.Path, `^/?user/([^/?#]+)`), "/user/")
	if secId == "" {
		return "", gerror.New("获取用户 secId 失败")
	}
	c := g.Client()
	c.SetAgent(consts.CommonAgent)
	resp, err := c.Get(ctx, user_info_url+secId)
	if err != nil {
		return "", gerror.Newf("请求用户信息失败: %s", err)
	}
	body, err := utils.Text(resp.Response)
	if err != nil {
		return "", gerror.Newf("格式化用户信息响应失败: %s", err)
	}
	var respBody TempAnchorInfoResp
	err = json.Unmarshal([]byte(body), &respBody)
	if err != nil {
		return "", gerror.Newf("转换用户数据失败: %s", err)
	}
	realUrl = douyin_live + respBody.UserInfo.UniqueID
	return
}

func extractSecUserIdFromPath(path string) string {
	_, after, found := strings.Cut(path, "/user/")
	if !found {
		return ""
	}
	return after
}
