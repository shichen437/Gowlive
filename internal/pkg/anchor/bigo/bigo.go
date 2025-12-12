package bigo

import (
	"context"
	"net/http"
	"net/http/cookiejar"
	"net/url"
	"strings"
	"time"

	"github.com/gogf/gf/v2/encoding/gjson"
	"github.com/gogf/gf/v2/errors/gerror"
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gctx"
	"github.com/gogf/gf/v2/util/gconv"
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
	return &Bigo{Url: u, Platform: platform}, nil
}

func (d *Bigo) ParseAnchorInfo(ctx context.Context) (info *anchor.AnchorInfo, err error) {
	if !strings.Contains(d.Url.String(), "/user/") {
		return nil, gerror.New("Bigo user homepage URL format error")
	}
	uid := utils.FindFirstMatch(d.Url.Path, regexpUid)
	uidArr := strings.Split(uid, "/")
	if len(uidArr) > 0 {
		uid = uidArr[len(uidArr)-1]
	}
	if uid == "" {
		return nil, gerror.New("Bigo uid not found")
	}
	body, err := d.getWebRequestResp(ctx, uid)
	if err != nil {
		return nil, err
	}
	bupr, err := parseJsonResp(body)
	if err != nil {
		return nil, err
	}
	info = &anchor.AnchorInfo{
		Platform:       d.Platform,
		AnchorName:     bupr.Data.Nickname,
		UniqueId:       bupr.Data.SecId,
		Signature:      bupr.Data.Sign,
		FollowerCount:  bupr.Data.FollowerCount,
		FollowingCount: bupr.Data.FollowingCount,
		LikeCount:      0,
		VideoCount:     gconv.Int(bupr.Data.VideoCount),
	}
	return
}

func (d *Bigo) getWebRequestResp(ctx context.Context, uid string) (string, error) {
	c := g.Client()
	c.SetTimeout(time.Second * 10)
	c.SetAgent(consts.CommonAgent)
	c.SetContentType(consts.CommonJsonType)
	c.SetHeaderMap(g.MapStrStr{
		"User-Agent":      consts.CommonAgent,
		"Referer":         domain,
		"Accept-Language": consts.CommonLang,
	})
	c.SetCookieMap(d.assembleCookieMap())
	params := g.MapStrStr{
		"bigoId": uid,
		"lang":   "tw",
	}
	resp, err := c.Post(ctx, anchorApi, params)
	if err != nil || resp.StatusCode != http.StatusOK {
		return "", gerror.New("Bigo anchor api request failed")
	}
	body, err := utils.Text(resp.Response)
	if err != nil {
		return "", err
	}
	return body, nil
}

func (d *Bigo) assembleCookieMap() map[string]string {
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

func parseJsonResp(body string) (*BigoUserProfileResp, error) {
	var bupr *BigoUserProfileResp
	err := gjson.Unmarshal([]byte(body), &bupr)
	if err != nil {
		return nil, gerror.New("Bigo platform parse anchor response failed")
	}
	if bupr.Code != 0 {
		err = gerror.New("Bigo platform get anchor info failed")
		return nil, err
	}
	return bupr, nil
}
