package bigo

import (
	"context"
	"maps"
	"net/http"
	"net/http/cookiejar"
	"net/url"
	"regexp"
	"strings"
	"time"

	"github.com/gogf/gf/v2/encoding/gjson"
	"github.com/gogf/gf/v2/errors/gerror"
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gctx"
	"github.com/shichen437/gowlive/internal/pkg/consts"
	"github.com/shichen437/gowlive/internal/pkg/lives"
	"github.com/shichen437/gowlive/internal/pkg/manager"
	"github.com/shichen437/gowlive/internal/pkg/metrics"
	"github.com/shichen437/gowlive/internal/pkg/utils"
)

func init() {
	lives.Register(domain, &builder{})
}

type builder struct{}

func (b *builder) Build(url *url.URL) (lives.LiveApi, error) {
	useragent, err := manager.GetUserAgentManager().GetUserAgent(platform)
	if err != nil {
		useragent = consts.CommonAgent
	}
	return &Bigo{
		Url:         url,
		Platform:    platform,
		UserAgent:   useragent,
		RespCookies: make(map[string]string),
	}, nil
}

func (l *Bigo) GetInfo() (info *lives.LiveState, err error) {
	info = &lives.LiveState{
		Platform: l.Platform,
	}
	uid := l.getUserId()
	if uid == "" {
		err = gerror.New("Bigo platform get uid failed")
		return
	}
	ctx := gctx.GetInitCtx()
	flag := lives.GetBucketManager().Acquire(ctx, platform)
	if flag != nil {
		err = gerror.New("Bigo platform get token failed")
		return
	}
	body, err := l.requestWebApi(ctx, uid)
	if err != nil {
		metrics.GetIndicatorManager().Record(gctx.GetInitCtx(), platform, false, true)
		return
	}
	if body == "" {
		metrics.GetIndicatorManager().Record(gctx.GetInitCtx(), platform, false, true)
		err = gerror.New("Bigo platform get info failed")
		return
	}
	metrics.GetIndicatorManager().Record(gctx.GetInitCtx(), platform, true, true)
	sir, err := parseJsonResp(body)
	if err != nil {
		return
	}
	info.Anchor = sir.Data.Nickname
	info.RoomName = sir.Data.RoomTopic
	l.parseStreamUrlInfo(sir, info)
	return info, nil
}

func (l *Bigo) requestWebApi(ctx context.Context, uid string) (string, error) {
	l.Proxy = manager.GetProxyManager().GetRandomProxy(platform)
	c := g.Client()
	if l.Proxy != "" {
		c.SetProxy(l.Proxy)
	}
	c.SetTimeout(time.Second * 10)
	headers := g.MapStrStr{
		"User-Agent":      l.UserAgent,
		"Referer":         referer,
		"Accept-Language": consts.CommonLang,
	}
	params := g.MapStrStr{
		"siteId": uid,
	}
	c.SetHeaderMap(headers)
	c.SetCookieMap(l.assembleCookieMap())
	g.Log().Info(gctx.GetInitCtx(), "Get Room Web Request: "+l.Url.String())
	resp, err := c.Post(ctx, studioInfoApi, params)
	if err != nil {
		g.Log().Error(gctx.GetInitCtx(), platform+" requestWebApi err info: ", err)
		return "", err
	}
	if resp.StatusCode != http.StatusOK {
		g.Log().Error(gctx.GetInitCtx(), platform+" requestWebApi response code: ", resp.StatusCode)
		return "", gerror.New("Bigo platform request failed")
	}
	body, err := utils.Text(resp.Response)
	if err != nil {
		return "", err
	}
	return body, nil
}

func (l *Bigo) getUserId() string {
	pattern := regexp.MustCompile(regexpUid)
	m := pattern.FindStringSubmatch(l.Url.String())
	if len(m) < 2 || strings.TrimSpace(m[1]) == "" {
		return ""
	}
	return m[1]
}

func (l *Bigo) parseStreamUrlInfo(sir *StudioInfoResp, info *lives.LiveState) {
	if sir.Data.Alive != 1 {
		return
	}
	sUrl, err := url.Parse(sir.Data.HlsSrc)
	if err != nil {
		return
	}
	streamUrlInfos := make([]*lives.StreamUrlInfo, 0, 1)
	streamUrlInfo := &lives.StreamUrlInfo{
		Url:                  sUrl,
		IsHls:                true,
		HeadersForDownloader: l.getHeadersForDownloader(),
	}
	streamUrlInfos = append(streamUrlInfos, streamUrlInfo)
	info.StreamInfos = streamUrlInfos
	info.IsLive = true
}

func (l *Bigo) assembleCookieMap() map[string]string {
	jar, _ := cookiejar.New(&cookiejar.Options{})
	cacheCookie := manager.GetCookieManager().Get(gctx.GetInitCtx(), l.Platform)
	jar.SetCookies(l.Url, utils.GetCookieList(cacheCookie))
	cookies := jar.Cookies(l.Url)
	cookieMap := make(map[string]string)
	maps.Copy(cookieMap, l.RespCookies)
	for _, c := range cookies {
		cookieMap[c.Name] = c.Value
	}
	return cookieMap
}

func (l *Bigo) getHeadersForDownloader() map[string]string {
	return map[string]string{
		"User-Agent": l.UserAgent,
		"Referer":    l.Url.String(),
		"Proxy":      l.Proxy,
	}
}

func parseJsonResp(body string) (*StudioInfoResp, error) {
	var sir *StudioInfoResp
	err := gjson.Unmarshal([]byte(body), &sir)
	if err != nil {
		return nil, gerror.New("Bigo platform parse response failed")
	}
	if sir.Code != 0 {
		err = gerror.New("Bigo platform get info failed")
		return nil, err
	}
	return sir, nil
}
