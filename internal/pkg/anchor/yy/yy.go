package yy

import (
	"context"
	"net/http"
	"net/http/cookiejar"
	"net/url"
	"regexp"
	"strings"

	"github.com/gogf/gf/v2/encoding/gjson"
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
	return &YY{Url: u, Platform: platform}, nil
}

func (d *YY) ParseAnchorInfo(ctx context.Context) (info *anchor.AnchorInfo, err error) {
	if !strings.Contains(d.Url.String(), "/u/") {
		return nil, gerror.New("YY用户主页 URL 格式错误")
	}
	secId := strings.TrimPrefix(utils.FindFirstMatch(d.Url.Path, `^/?u/([^/?#]+)`), "/u/")
	if secId == "" {
		return nil, gerror.New("获取用户 secId 失败")
	}
	cookieMap := d.assembleCookieMap()
	headersMap := g.MapStrStr{
		"User-Agent":      consts.CommonAgent,
		"Referer":         domain,
		"Accept-Language": consts.CommonLang,
	}
	body, err := d.getWebHomepageResp(ctx, headersMap, cookieMap)
	if err != nil {
		return nil, err
	}
	dataStr := extractPagePropsByScripts(body)
	if dataStr == "" {
		return nil, gerror.New("解析主页信息失败")
	}
	dataJson := gjson.New(dataStr)
	userInfo := &anchor.AnchorInfo{
		Platform:       platform,
		UniqueId:       dataJson.Get("userBaseInfo.uid").String(),
		AnchorName:     dataJson.Get("userBaseInfo.topData.baseUserInfo.anchorName").String(),
		Signature:      dataJson.Get("userBaseInfo.topData.baseUserInfo.sign").String(),
		FollowerCount:  int(dataJson.Get("userBaseInfo.topData.attentionCountInfo.fansCount").Int()),
		FollowingCount: int(dataJson.Get("userBaseInfo.topData.attentionCountInfo.attentionCount").Int()),
		LikeCount:      0,
		VideoCount:     0,
	}
	videoId := getVideoCount(ctx, secId, headersMap, cookieMap)
	if videoId != 0 {
		userInfo.VideoCount = int(videoId)
	}
	return userInfo, nil
}

func (d *YY) getWebHomepageResp(ctx context.Context, hMap, cMap map[string]string) (string, error) {
	c := g.Client()
	c.SetAgent(consts.CommonAgent)
	c.SetHeaderMap(hMap)
	c.SetCookieMap(cMap)
	resp, err := c.Get(ctx, d.Url.String())
	if err != nil || resp.StatusCode != http.StatusOK {
		return "", gerror.New("获取YY用户主页失败")
	}
	body, err := utils.Text(resp.Response)
	if err != nil {
		return "", err
	}
	return body, nil
}

func getVideoCount(ctx context.Context, uid string, hMap, cMap map[string]string) int {
	c := g.Client()
	c.SetAgent(consts.CommonAgent)
	c.SetHeaderMap(hMap)
	c.SetCookieMap(cMap)
	resp, err := c.Get(ctx, video_data_url+uid)
	if err != nil || resp.StatusCode != http.StatusOK {
		return 0
	}
	body, err := utils.Text(resp.Response)
	if err != nil {
		return 0
	}
	return gjson.New(body).Get("data.duanpaiPage.totalCount").Int()
}

func (d *YY) assembleCookieMap() map[string]string {
	jar, _ := cookiejar.New(&cookiejar.Options{})
	cacheCookie := manager.GetCookieManager().Get(gctx.GetInitCtx(), d.Platform)
	jar.SetCookies(d.Url, utils.GetCookieList(cacheCookie))
	cookies := jar.Cookies(d.Url)
	cookieMap := make(map[string]string)
	for _, c := range cookies {
		cookieMap[c.Name] = c.Value
	}
	return cookieMap
}

func extractPagePropsByScripts(body string) string {
	props := regexp.MustCompile(regexpProps)
	result := props.FindStringSubmatch(body)
	if len(result) < 2 {
		return ""
	}
	return result[1]
}
