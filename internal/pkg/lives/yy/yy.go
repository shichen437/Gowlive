package yy

import (
	"context"
	"fmt"
	"maps"
	"net/http/cookiejar"
	"net/url"
	"regexp"

	"github.com/gogf/gf/v2/encoding/gjson"
	"github.com/gogf/gf/v2/errors/gerror"
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gctx"
	"github.com/gogf/gf/v2/os/gtime"
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
	return &YY{
		Url:         url,
		Platform:    platform,
		UserAgent:   useragent,
		RespCookies: make(map[string]string),
	}, nil
}

func (l *YY) GetInfo() (info *lives.LiveState, err error) {
	info = &lives.LiveState{
		Platform: l.Platform,
	}
	ctx := gctx.GetInitCtx()
	flag := lives.GetBucketManager().Acquire(ctx, platform)
	if flag != nil {
		err = gerror.New("YY直播间获取令牌失败")
		return
	}
	headers := g.MapStrStr{
		"User-Agent":      l.UserAgent,
		"Referer":         referer,
		"Accept-Language": consts.CommonLang,
	}
	cookieMap := l.assembleCookieMap()
	body, err := l.requestWebApi(ctx, headers, cookieMap)
	if err != nil {
		metrics.GetIndicatorManager().Record(gctx.GetInitCtx(), platform, false, false)
		return
	}
	metrics.GetIndicatorManager().Record(gctx.GetInitCtx(), platform, true, false)
	anchorName, err := parseAnchorName(body)
	if err != nil {
		return nil, gerror.New("parse anchor name failed")
	}
	cid, err := parseCID(body)
	if err != nil {
		return nil, gerror.New("parse cid failed")
	}
	info.Anchor = anchorName
	detailResp, err := requestWebDetailApi(ctx, cid, headers, cookieMap)
	if err != nil {
		metrics.GetIndicatorManager().Record(gctx.GetInitCtx(), platform, false, false)
		return
	}
	metrics.GetIndicatorManager().Record(gctx.GetInitCtx(), platform, true, false)
	detailJson := gjson.New(detailResp)
	roomName := detailJson.Get("data.roomName").String()
	info.IsLive = false
	if roomName == "" {
		return
	}
	info.RoomName = roomName
	streamResp, err := requestWebStreamInfoApi(ctx, cid, headers, cookieMap)
	if err != nil {
		metrics.GetIndicatorManager().Record(gctx.GetInitCtx(), platform, false, true)
		return info, gerror.New("解析YY直播间流地址失败")
	}
	metrics.GetIndicatorManager().Record(gctx.GetInitCtx(), platform, true, true)
	err = l.parseStreamUrl(streamResp, info)
	if err != nil {
		return
	}
	return info, nil
}

func (l *YY) requestWebApi(ctx context.Context, headers, cookieMap map[string]string) (string, error) {
	c := g.Client()
	c.SetHeaderMap(headers)
	c.SetCookieMap(cookieMap)
	g.Log().Info(gctx.GetInitCtx(), "Get Room Web Request: "+l.Url.String())
	resp, err := c.Get(ctx, l.Url.String())
	if err != nil {
		return "", err
	}
	body, err := utils.Text(resp.Response)
	if err != nil {
		return "", err
	}
	return body, nil
}

func (l *YY) getHeadersForDownloader() map[string]string {
	return map[string]string{
		"User-Agent": consts.CommonAgent,
		"Referer":    l.Url.String(),
	}
}

func (l *YY) parseStreamUrl(body string, info *lives.LiveState) error {
	streamJson := gjson.New(body)
	sArr := streamJson.Get("avp_info_res.stream_line_addr").Array()
	if len(sArr) == 0 {
		return gerror.New("解析YY直播间流地址失败")
	}
	sm := gjson.New(sArr[0]).Map()
	var sUrl *url.URL
	var sUrlString string
	for _, v := range sm {
		realUrl := gjson.New(v).Get("cdn_info.url").String()
		if realUrl == "" {
			continue
		}
		sUrlString = realUrl
		break
	}
	sUrl, err := url.Parse(sUrlString)
	if err != nil || sUrlString == "" {
		return gerror.New("解析YY直播间流地址失败")
	}
	streamUrlInfos := make([]*lives.StreamUrlInfo, 0, 1)
	streamUrlInfo := &lives.StreamUrlInfo{
		Url:                  sUrl,
		HeadersForDownloader: l.getHeadersForDownloader(),
	}
	streamUrlInfos = append(streamUrlInfos, streamUrlInfo)
	info.StreamInfos = streamUrlInfos
	info.IsLive = true
	return nil
}

func requestWebStreamInfoApi(ctx context.Context, cid string, headers, cookieMap map[string]string) (string, error) {
	dataJSON := `{"head":{"seq":1701869217590,"appidstr":"0","bidstr":"121","cidstr":"` + cid + `","sidstr":"` + cid + `","uid64":0,"client_type":108,"client_ver":"5.17.0","stream_sys_ver":1,"app":"yylive_web","playersdk_ver":"5.17.0","thundersdk_ver":"0","streamsdk_ver":"5.17.0"},"client_attribute":{"client":"web","model":"web0","cpu":"","graphics_card":"","os":"chrome","osversion":"0","vsdk_version":"","app_identify":"","app_version":"","business":"","width":"1920","height":"1080","scale":"","client_type":8,"h265":0},"avp_parameter":{"version":1,"client_type":8,"service_type":0,"imsi":0,"send_time":1701869217,"line_seq":-1,"gear":4,"ssl":1,"stream_format":0}}`
	q := url.Values{}
	q.Set("uid", "0")
	q.Set("cid", cid)
	q.Set("sid", cid)
	q.Set("appid", "0")
	q.Set("sequence", "1701869217590")
	q.Set("encode", "json")
	streamAPI := streamApi + q.Encode()
	c := g.Client()
	c.SetHeaderMap(headers)
	c.SetCookieMap(cookieMap)
	resp, err := c.Post(ctx, streamAPI, dataJSON)
	if err != nil {
		return "", err
	}
	return utils.Text(resp.Response)
}

func requestWebDetailApi(ctx context.Context, cid string, headers, cookieMap map[string]string) (string, error) {
	detailParams := url.Values{}
	detailParams.Set("uid", "")
	detailParams.Set("sid", cid)
	detailParams.Set("ssid", cid)
	detailParams.Set("_", fmt.Sprintf("%d", gtime.Now().UnixMilli()))
	detailAPI := detailApi + detailParams.Encode()
	c := g.Client()
	c.SetHeaderMap(headers)
	c.SetCookieMap(cookieMap)
	resp, err := c.Get(ctx, detailAPI)
	if err != nil {
		return "", err
	}
	return utils.Text(resp.Response)
}

func parseAnchorName(body string) (string, error) {
	re := regexp.MustCompile(regexpAnchor)
	m := re.FindStringSubmatch(body)
	if len(m) < 2 {
		return "", gerror.New("anchor name not found")
	}
	return m[1], nil
}

func parseCID(body string) (string, error) {
	re := regexp.MustCompile(regexpCid)
	m := re.FindStringSubmatch(body)
	if len(m) < 2 {
		return "", gerror.New("cid (sid) not found")
	}
	return m[1], nil
}

func (l *YY) assembleCookieMap() map[string]string {
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
