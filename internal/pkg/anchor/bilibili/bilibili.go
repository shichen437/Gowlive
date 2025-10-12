package bilibili

import (
	"context"
	"net/url"
	"strings"

	"github.com/gogf/gf/v2/errors/gerror"
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gctx"
	"github.com/shichen437/gowlive/internal/pkg/anchor"
	"github.com/shichen437/gowlive/internal/pkg/manager"
	"github.com/shichen437/gowlive/internal/pkg/utils"
	"github.com/tidwall/gjson"
)

func init() {
	anchor.Register(domain, &builder{})
}

type builder struct{}

func (b *builder) Build(u *url.URL) (anchor.AnchorApi, error) {
	return &Bilibili{Url: u, Platform: platform}, nil
}

func (d *Bilibili) ParseAnchorInfo(ctx context.Context) (info *anchor.AnchorInfo, err error) {
	var mid string
	if strings.Contains(d.Url.String(), domain) {
		mid = strings.TrimPrefix(utils.FindFirstMatch(d.Url.String(), `space\.bilibili\.com/(\d+)`), "space.bilibili.com/")
	}
	if mid == "" {
		return nil, gerror.New("B站主页链接格式错误")
	}
	info, err = d.getUserStatInfo(ctx, mid)
	if err != nil {
		return nil, err
	}
	return info, nil
}

func (d *Bilibili) getUserStatInfo(ctx context.Context, mid string) (*anchor.AnchorInfo, error) {
	c := g.Client()
	c.SetHeaderMap(bilibiliHeaders)
	cookie := manager.GetCookieManager().Get(gctx.GetInitCtx(), d.Platform)
	c.SetCookieMap(utils.GetCookieMap(d.Platform, cookie, d.Url.String()))
	resp, err := c.Get(ctx, userProfileInfoUrl+mid)
	if err != nil {
		return nil, gerror.New("请求用户信息失败")
	}
	defer resp.Close()
	body, err := utils.Text(resp.Response)
	if err != nil {
		return nil, gerror.New("请求用户信息失败")
	}
	jsonData := gjson.Parse(body)
	if jsonData.Get("code").Int() != 0 {
		return nil, gerror.New("请求用户信息失败")
	}
	info := jsonData.Get("data")
	userInfo := &anchor.AnchorInfo{
		UniqueId:       mid,
		Platform:       platform,
		AnchorName:     info.Get("card.name").String(),
		Signature:      info.Get("card.sign").String(),
		FollowerCount:  int(info.Get("follower").Int()),
		FollowingCount: int(info.Get("card.friend").Int()),
		LikeCount:      int(info.Get("like_num").Int()),
		VideoCount:     int(info.Get("archive_count").Int()),
	}
	return userInfo, nil
}
