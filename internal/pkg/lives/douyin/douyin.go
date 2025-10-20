package douyin

import (
	"fmt"
	"maps"
	"net/http"
	"net/http/cookiejar"
	"net/url"
	"regexp"
	"sort"
	"strconv"
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
	return &Douyin{
		Url:         url,
		Platform:    platform,
		RespCookies: make(map[string]string),
	}, nil
}

func (l *Douyin) GetInfo() (info *lives.LiveState, err error) {
	info = &lives.LiveState{
		Platform: l.Platform,
	}
	flag := lives.GetBucketManager().Acquire(gctx.GetInitCtx(), platform)
	if flag != nil {
		err = gerror.New("抖音直播间获取令牌失败")
		return
	}
	body, err := l.getRoomWebPageResp()
	if err != nil {
		err = gerror.New("访问抖音直播间页面失败")
		return
	}
	json, err := getMainInfoLine(body)
	if err != nil {
		err = gerror.New("获取抖音直播间主信息失败")
		return
	}
	info.RoomName = json.Get("state.roomStore.roomInfo.room.title").String()
	info.Anchor = json.Get("state.roomStore.roomInfo.anchor.nickname").String()
	isStreaming := json.Get("state.roomStore.roomInfo.room.status_str").String() == "2"
	info.IsLive = isStreaming
	if !isStreaming {
		return
	}
	streamIdPath := "state.streamStore.streamData.H264_streamData.common.stream"
	streamId := json.Get(streamIdPath).String()
	streamUrlInfos, _ := getStreamInfo(body, streamId)
	info.StreamInfos = streamUrlInfos
	return
}

func (l *Douyin) getRoomWebPageResp() (body string, err error) {
	c := g.Client()
	cookieMap := l.assembleCookieMap()
	c.SetCookieMap(cookieMap)
	req, err := c.Get(gctx.GetInitCtx(), l.Url.String())
	g.Log().Info(gctx.GetInitCtx(), "Get Room Web Page: "+l.Url.String())
	if err != nil {
		metrics.GetIndicatorManager().Record(gctx.GetInitCtx(), platform, false, false)
		g.Log().Error(gctx.GetInitCtx(), err.Error())
		return
	}
	metrics.GetIndicatorManager().Record(gctx.GetInitCtx(), platform, true, false)
	cookieWithOdinTt := fmt.Sprintf("odin_tt=%s; %s", utils.GenRandomString(160, randomCookieChars), req.Header.Get("Cookie"))
	req.Header.Set("Cookie", cookieWithOdinTt)
	c2 := g.Client()
	resp, err := c2.Do(req.Request)
	if err != nil {
		metrics.GetIndicatorManager().Record(gctx.GetInitCtx(), platform, false, true)
		g.Log().Error(gctx.GetInitCtx(), err.Error())
		return
	}
	switch code := resp.StatusCode; code {
	case http.StatusOK:
		metrics.GetIndicatorManager().Record(gctx.GetInitCtx(), platform, true, true)
		body, err = utils.Text(resp)
		if err != nil {
			g.Log().Error(gctx.GetInitCtx(), err.Error())
			return
		}
		for _, cookie := range resp.Cookies() {
			l.RespCookies[cookie.Name] = cookie.Value
		}
		return
	default:
		metrics.GetIndicatorManager().Record(gctx.GetInitCtx(), platform, false, true)
		err = gerror.Newf(`http response error`)
		return
	}
}

func getMainInfoLine(body string) (json *gjson.Result, err error) {
	reg, err := regexp.Compile(mainInfoLineCatcherRegex)
	if err != nil {
		return
	}
	match := reg.FindAllStringSubmatch(body, -1)
	if match == nil {
		err = fmt.Errorf("0 match for mainInfoLineCatcherRegex: %s", mainInfoLineCatcherRegex)
		return
	}
	for _, item := range match {
		if len(item) < 2 {
			continue
		}
		mainInfoLine := item[1]

		// 获取房间信息
		mainJson := gjson.Parse(fmt.Sprintf(`"%s"`, mainInfoLine))
		if !mainJson.Exists() {
			continue
		}

		mainJson = gjson.Parse(mainJson.String()).Get("3")
		if !mainJson.Exists() {
			continue
		}

		if mainJson.Get("state.roomStore.roomInfo.room.status_str").Exists() {
			json = &mainJson
			return
		}
	}
	return nil, fmt.Errorf("MainInfoLine not found")
}

func getStreamInfo(body string, streamId string) (infos []*lives.StreamUrlInfo, err error) {
	streamUrlInfos := make([]*lives.StreamUrlInfo, 0, 10)
	reg2, _ := regexp.Compile(commonInfoLineCatcherRegex)
	match2 := reg2.FindAllStringSubmatch(body, -1)
	for _, item := range match2 {
		if len(item) < 2 {
			return
		}
		commonJson := gjson.Parse(gjson.Parse(fmt.Sprintf(`"%s"`, item[1])).String())
		if !commonJson.Exists() {
			return
		}
		if !commonJson.Get("common").Exists() {
			continue
		}
		commonStreamId := commonJson.Get("common.stream").String()
		if commonStreamId == "" {
			return
		}
		if commonStreamId != streamId {
			continue
		}
		commonJson.Get("data").ForEach(func(key, value gjson.Result) bool {
			flv := value.Get("main.flv").String()
			var Url *url.URL
			Url, err = url.Parse(flv)
			if err != nil {
				err = gerror.New("解析抖音直播间流地址失败")
				return false
			}
			paramsString := value.Get("main.sdk_params").String()
			paramsJson := gjson.Parse(paramsString)
			var description strings.Builder
			paramsJson.ForEach(func(key, value gjson.Result) bool {
				description.WriteString(key.String())
				description.WriteString(": ")
				description.WriteString(value.String())
				description.WriteString("\n")
				return true
			})
			Resolution := 0
			resolution := strings.Split(paramsJson.Get("resolution").String(), "x")
			if len(resolution) == 2 {
				x, err := strconv.Atoi(resolution[0])
				if err != nil {
					return true
				}
				y, err := strconv.Atoi(resolution[1])
				if err != nil {
					return true
				}
				Resolution = x * y
			}
			Vbitrate := int(paramsJson.Get("vbitrate").Int())
			streamUrlInfos = append(streamUrlInfos, &lives.StreamUrlInfo{
				Name:        key.String(),
				Description: description.String(),
				Url:         Url,
				Resolution:  Resolution,
				Vbitrate:    Vbitrate,
			})
			return true
		})
		sort.Slice(streamUrlInfos, func(i, j int) bool {
			if streamUrlInfos[i].Resolution != streamUrlInfos[j].Resolution {
				return streamUrlInfos[i].Resolution > streamUrlInfos[j].Resolution
			} else {
				return streamUrlInfos[i].Vbitrate > streamUrlInfos[j].Vbitrate
			}
		})
	}
	infos = streamUrlInfos
	return
}

func (l *Douyin) assembleCookieMap() map[string]string {
	jar, _ := cookiejar.New(&cookiejar.Options{})
	c := manager.GetCookieManager().Get(gctx.GetInitCtx(), l.Platform)
	jar.SetCookies(l.Url, utils.GetCookieList(c))
	cookies := jar.Cookies(l.Url)
	cookieMap := make(map[string]string)
	cookieMap["__ac_nonce"] = utils.GenRandomString(21, randomCookieChars)
	maps.Copy(cookieMap, l.RespCookies)
	for _, c := range cookies {
		cookieMap[c.Name] = c.Value
	}
	return cookieMap
}
