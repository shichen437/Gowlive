package douyin

import (
	"context"
	"encoding/json"
	"net/http/cookiejar"
	"net/url"
	"strings"

	"github.com/gogf/gf/v2/errors/gerror"
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gctx"
	"github.com/shichen437/gowlive/internal/pkg/anchor"
	"github.com/shichen437/gowlive/internal/pkg/consts"
	"github.com/shichen437/gowlive/internal/pkg/manager"
	"github.com/shichen437/gowlive/internal/pkg/utils"
)

func init() {
	anchor.Register(domain, &builder{})
}

type builder struct{}

func (b *builder) Build(u *url.URL) (anchor.AnchorApi, error) {
	return &Douyin{Url: u, Platform: platform}, nil
}

func (d *Douyin) ParseAnchorInfo(ctx context.Context) (info *anchor.AnchorInfo, err error) {
	if !strings.Contains(d.Url.String(), "/user/") {
		return nil, gerror.New("抖音用户主页 URL 格式错误")
	}
	secId := strings.TrimPrefix(utils.FindFirstMatch(d.Url.Path, `^/?user/([^/?#]+)`), "/user/")
	if secId == "" {
		return nil, gerror.New("获取用户 secId 失败")
	}
	info = &anchor.AnchorInfo{
		Platform: d.Platform,
	}
	c := g.Client()
	c.SetAgent(consts.CommonAgent)
	c.SetCookieMap(d.assembleCookieMap())
	resp, err := c.Get(ctx, user_info_url+secId)
	if err != nil {
		return
	}
	body, err := utils.Text(resp.Response)
	if err != nil {
		return
	}
	var respBody AnchorInfoResp
	err = json.Unmarshal([]byte(body), &respBody)
	if err != nil {
		return
	}
	info.UniqueId = respBody.UserInfo.UniqueID
	info.AnchorName = respBody.UserInfo.Nickname
	info.Signature = respBody.UserInfo.Signature
	info.FollowingCount = respBody.UserInfo.FollowingCount
	info.FollowerCount = respBody.UserInfo.FollowerCount
	info.VideoCount = respBody.UserInfo.VideoCount
	info.LikeCount = utils.ParseChineseNumberToInt(respBody.UserInfo.LikeCount)
	return info, nil
}

func (d *Douyin) assembleCookieMap() map[string]string {
	jar, _ := cookiejar.New(&cookiejar.Options{})
	c := manager.GetCookieManager().Get(gctx.GetInitCtx(), d.Platform)
	jar.SetCookies(d.Url, utils.GetCookieList(c))
	cookies := jar.Cookies(d.Url)
	cookieMap := make(map[string]string)
	cookieMap["__ac_nonce"] = utils.GenRandomString(21, randomCookieChars)
	for _, c := range cookies {
		cookieMap[c.Name] = c.Value
	}
	return cookieMap
}
