package bilibili

import (
	"fmt"
	"net/http/cookiejar"
	"net/url"
	"strings"

	"github.com/gogf/gf/v2/errors/gerror"
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gctx"
	"github.com/shichen437/gowlive/internal/pkg/lives"
	"github.com/shichen437/gowlive/internal/pkg/manager"
	"github.com/shichen437/gowlive/internal/pkg/metrics"
	"github.com/shichen437/gowlive/internal/pkg/utils"
	"github.com/tidwall/gjson"
)

func init() {
	lives.Register(domain, &builder{})
}

type builder struct{}

func (b *builder) Build(url *url.URL) (lives.LiveApi, error) {
	return &Bilibili{
		Url:         url,
		Platform:    platform,
		RespCookies: make(map[string]string),
	}, nil
}

type Bilibili struct {
	Url         *url.URL
	Platform    string
	RoomID      string
	RespCookies map[string]string
}

func (l *Bilibili) GetInfo() (info *lives.LiveState, err error) {
	if l.RoomID == "" && l.parseRoomID() != nil {
		return nil, gerror.New(l.Platform + "获取房间ID失败")
	}
	flag := lives.GetBucketManager().Acquire(gctx.GetInitCtx(), platform)
	if flag != nil {
		err = gerror.New("B站直播间获取令牌失败")
		return
	}
	info, err = l.getRoomInfo()
	if err != nil {
		metrics.GetIndicatorManager().Record(gctx.GetInitCtx(), platform, false, false)
		return nil, err
	}
	metrics.GetIndicatorManager().Record(gctx.GetInitCtx(), platform, true, false)
	err = l.getUserInfo(info)
	if err != nil {
		metrics.GetIndicatorManager().Record(gctx.GetInitCtx(), platform, false, false)
		return nil, err
	}
	metrics.GetIndicatorManager().Record(gctx.GetInitCtx(), platform, true, false)
	if info.IsLive {
		streamInfos, err := l.getStreamInfo()
		if err != nil || len(streamInfos) == 0 {
			metrics.GetIndicatorManager().Record(gctx.GetInitCtx(), platform, false, true)
			return nil, gerror.New(l.Platform + "获取直播流数据失败")
		}
		metrics.GetIndicatorManager().Record(gctx.GetInitCtx(), platform, true, true)
		info.StreamInfos = streamInfos
	}
	return info, nil
}

func (l *Bilibili) getStreamInfo() (infos []*lives.StreamUrlInfo, err error) {
	c := g.Client()
	c.SetAgent(userAgent)
	c.SetCookieMap(l.assembleCookieMap())
	payload := fmt.Sprintf(`?room_id=%s&protocol=0,1&format=0,1,2&codec=0,1&qn=10000&platform=web&ptype=8&dolby=5&panorama=1`, l.RoomID)
	resp, err := c.Get(gctx.GetInitCtx(), liveApiUrlv2+payload)
	if err != nil || resp.StatusCode != 200 {
		return nil, gerror.New(l.Platform + "获取直播信息失败")
	}
	defer resp.Close()
	body := resp.ReadAllString()
	jsonData := gjson.Parse(body)
	baseUrl := jsonData.Get(parseAddr + ".base_url").String()
	jsonArr := jsonData.Get(parseAddr + ".url_info").Array()
	streamUrlInfos := make([]*lives.StreamUrlInfo, 0, 10)
	for _, v := range jsonArr {
		hosts := gjson.Get(v.String(), "host").String()
		queries := gjson.Get(v.String(), "extra").String()
		streamUrl, err := url.Parse(hosts + baseUrl + queries)
		if err != nil {
			continue
		}
		streamUrlInfos = append(streamUrlInfos, &lives.StreamUrlInfo{
			Url:                  streamUrl,
			HeadersForDownloader: l.getHeadersForDownloader(),
		})
	}
	return streamUrlInfos, nil
}

func (l *Bilibili) getUserInfo(info *lives.LiveState) error {
	c := g.Client()
	c.SetAgent(userAgent)
	c.SetCookieMap(l.assembleCookieMap())
	resp, err := c.Get(gctx.GetInitCtx(), userApiUrl, g.Map{
		"roomid": l.RoomID,
	})
	if err != nil || resp.StatusCode != 200 {
		return gerror.New(l.Platform + "获取房间信息失败")
	}
	defer resp.Close()
	body := resp.ReadAllString()
	jsonData := gjson.Parse(body)
	if jsonData.String() == "" || jsonData.Get("code").Int() != 0 {
		return gerror.New(l.Platform + "获取房间信息数据失败")
	}
	info.Anchor = jsonData.Get("data.info.uname").String()
	return nil
}

func (l *Bilibili) getRoomInfo() (*lives.LiveState, error) {
	c := g.Client()
	c.SetAgent(userAgent)
	c.SetCookieMap(l.assembleCookieMap())
	resp, err := c.Get(gctx.GetInitCtx(), roomApiUrl, g.Map{
		"room_id": l.RoomID,
		"from":    "room",
	})
	if err != nil || resp.StatusCode != 200 {
		return nil, gerror.New(l.Platform + "获取房间信息失败")
	}
	defer resp.Close()
	body := resp.ReadAllString()
	jsonData := gjson.Parse(body)
	if jsonData.String() == "" || jsonData.Get("code").Int() != 0 {
		return nil, gerror.New(l.Platform + "获取房间信息数据失败")
	}
	info := &lives.LiveState{
		Platform: l.Platform,
		RoomName: jsonData.Get("data.title").String(),
		IsLive:   jsonData.Get("data.live_status").Int() == 1,
	}
	return info, nil
}

func (l *Bilibili) parseRoomID() error {
	paths := strings.Split(l.Url.Path, "/")
	if len(paths) < 2 {
		return gerror.New(l.Platform + "无效链接")
	}
	c := g.Client()
	c.SetAgent(userAgent)
	c.SetCookieMap(l.assembleCookieMap())
	resp, err := c.Get(gctx.GetInitCtx(), roomInitUrl, g.Map{
		"id": paths[1],
	})
	if err != nil || resp.StatusCode != 200 {
		return gerror.New(l.Platform + "获取房间信息失败")
	}
	body, err := utils.Text(resp.Response)
	fmt.Println("Response Body:", body)
	jsonData := gjson.Parse(body)
	if jsonData.String() == "" || jsonData.Get("code").Int() != 0 {
		return gerror.New(l.Platform + "获取房间信息数据失败")
	}
	l.RoomID = jsonData.Get("data.room_id").String()
	return nil
}

func (l *Bilibili) assembleCookieMap() map[string]string {
	jar, _ := cookiejar.New(&cookiejar.Options{})
	cacheCookie := manager.GetCookieManager().Get(gctx.GetInitCtx(), l.Platform)
	jar.SetCookies(l.Url, utils.GetCookieList(cacheCookie))
	cookies := jar.Cookies(l.Url)
	cookieMap := make(map[string]string)
	for k, v := range l.RespCookies {
		cookieMap[k] = v
	}
	for _, c := range cookies {
		cookieMap[c.Name] = c.Value
	}
	return cookieMap
}

func (l *Bilibili) getHeadersForDownloader() map[string]string {
	return map[string]string{
		"User-Agent": userAgent,
		"Referer":    l.Url.String(),
	}
}
