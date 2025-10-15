package service

import (
	"context"
	"encoding/json"
	"sort"
	"strconv"
	"strings"
	"time"

	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gcache"
	"github.com/gogf/gf/v2/os/gtime"
	"github.com/shichen437/gowlive/internal/app/stream/dao"
	"github.com/shichen437/gowlive/internal/app/stream/model/do"
	"github.com/shichen437/gowlive/internal/app/stream/model/entity"
	"github.com/shichen437/gowlive/internal/pkg/consts"
	"github.com/shichen437/gowlive/internal/pkg/lives"
	"github.com/shichen437/gowlive/internal/pkg/utils"
)

var cacheSectors = gcache.New()

type LiveDuration struct {
	StartTime string
	EndTime   string
}

type arcSector struct {
	startMin int
	endMin   int
}

func GenLiveSessionById(ctx context.Context, liveId int) *lives.LiveSession {
	liveSession := &lives.LiveSession{
		Id: liveId,
	}
	var liveConfig *lives.LiveConfig
	err := dao.LiveManage.Ctx(ctx).WherePri(liveId).Scan(&liveConfig)
	if err != nil || liveConfig == nil {
		g.Log().Error(ctx, "获取直播配置失败", err)
		return nil
	}
	liveSession.Config = *liveConfig
	var liveState *lives.LiveState
	err = dao.LiveRoomInfo.Ctx(ctx).Where(dao.LiveRoomInfo.Columns().LiveId, liveId).Scan(&liveState)
	if err != nil || liveState == nil {
		g.Log().Error(ctx, "获取直播房间信息失败", err)
		return nil
	}
	liveSession.State = *liveState
	return liveSession
}

func GenLiveSessionsByIds(ctx context.Context, liveIds []int) []*lives.LiveSession {
	liveSessions := make([]*lives.LiveSession, 0, len(liveIds))
	var liveConfigs []*lives.LiveConfig
	err := dao.LiveManage.Ctx(ctx).WhereIn(dao.LiveManage.Columns().Id, liveIds).Scan(&liveConfigs)
	if err != nil || liveConfigs == nil || len(liveConfigs) == 0 {
		g.Log().Error(ctx, "获取直播配置失败", err)
		return nil
	}

	var liveRoomInfos []*lives.LiveState
	err = dao.LiveRoomInfo.Ctx(ctx).WhereIn(dao.LiveRoomInfo.Columns().LiveId, liveIds).Scan(&liveRoomInfos)
	if err != nil || liveRoomInfos == nil || len(liveRoomInfos) == 0 {
		g.Log().Error(ctx, "获取直播房间信息失败", err)
		return nil
	}
	liveStateMap := make(map[int]*lives.LiveState)
	for _, liveRoomInfo := range liveRoomInfos {
		liveStateMap[liveRoomInfo.LiveId] = liveRoomInfo
	}
	for _, liveConfig := range liveConfigs {
		liveState, ok := liveStateMap[liveConfig.Id]
		if !ok {
			g.Log().Error(ctx, "live state not found, liveId: %d", liveConfig.Id)
			continue
		}
		liveSessions = append(liveSessions, &lives.LiveSession{
			Id:     liveConfig.Id,
			Config: *liveConfig,
			State:  *liveState,
		})
	}
	return liveSessions
}

func UpdateRoomInfo(ctx context.Context, liveSession *lives.LiveSession) {
	state := liveSession.State
	dao.LiveRoomInfo.Ctx(ctx).Where(dao.LiveRoomInfo.Columns().LiveId, liveSession.Id).
		Update(do.LiveRoomInfo{
			RoomName:  state.RoomName,
			Anchor:    state.Anchor,
			UpdatedAt: utils.Now(),
		})
}

func GenIntelligentInterval(ctx context.Context, liveId int) int {
	list := getCircularSectorsWithCache(ctx, liveId)
	if len(list) <= 0 {
		return consts.DefaultInterval
	}
	if isOutsideLiveDuration(list) {
		return consts.MaxInterval
	}
	return consts.DefaultInterval
}

func getCircularSectorsWithCache(ctx context.Context, liveId int) []*LiveDuration {
	key := consts.DurationCacheKeyPrefix + strconv.Itoa(liveId)
	if v, err := cacheSectors.Get(ctx, key); err == nil && v != nil {
		if s := v.String(); s != "" {
			var sectors []*LiveDuration
			if json.Unmarshal([]byte(s), &sectors) == nil {
				return sectors
			}
		}
	}

	sectors := computeCircularSectors(ctx, liveId)
	if sectors == nil {
		sectors = make([]*LiveDuration, 0)
	}

	if b, err := json.Marshal(sectors); err == nil {
		_ = cacheSectors.Set(ctx, key, string(b), consts.DurationExpired)
	}

	return sectors
}

func computeCircularSectors(ctx context.Context, liveId int) []*LiveDuration {
	hList := getLiveHistory(ctx, liveId)
	if len(hList) <= 0 {
		return nil
	}
	const gapMin = 30
	const maxSpanMin = 20 * 60

	sectors := make([]arcSector, 0, len(hList)*2)
	for _, h := range hList {
		if h.StartedAt == nil || h.EndedAt == nil {
			continue
		}
		if h.EndedAt.Time.Sub(h.StartedAt.Time) < 30*time.Minute {
			continue
		}
		sectors = append(sectors, sessionToArcSectors(h.StartedAt, h.EndedAt)...)
	}
	if len(sectors) == 0 {
		return nil
	}

	merged := mergeArcSectors(sectors, gapMin)

	// 超长检测
	for _, s := range merged {
		if sectorSpanMin(s) > maxSpanMin {
			return nil
		}
		// 边界防御
		if s.startMin < 0 || s.endMin > 1440 || s.startMin >= s.endMin {
			return nil
		}
	}

	out := make([]*LiveDuration, 0, len(merged))
	for _, s := range merged {
		out = append(out, &LiveDuration{
			StartTime: minuteOffsetToHHMM(s.startMin),
			EndTime:   minuteOffsetToHHMM(s.endMin),
		})
	}
	return out
}

func getLiveHistory(ctx context.Context, liveId int) []*entity.LiveHistory {
	m := dao.LiveHistory.Ctx(ctx).Where(dao.LiveHistory.Columns().LiveId, liveId).
		WhereGTE(dao.LiveHistory.Columns().Duration, 0.5)
	count, err := m.Count()
	if err != nil || count <= 3 {
		return nil
	}
	var hList []*entity.LiveHistory
	m.OrderDesc(dao.LiveHistory.Columns().Id).
		Limit(30).Scan(&hList)
	return hList
}

func sectorSpanMin(s arcSector) int {
	return s.endMin - s.startMin
}

func sessionToArcSectors(start *gtime.Time, end *gtime.Time) []arcSector {
	if start == nil || end == nil {
		return nil
	}
	st := start.Time
	ed := end.Time
	if !ed.After(st) {
		return nil
	}

	startMin := st.Hour()*60 + st.Minute()
	endMin := ed.Hour()*60 + ed.Minute()

	// 同日且顺序正常：单段
	if st.Year() == ed.Year() && st.Month() == ed.Month() && st.Day() == ed.Day() && endMin > startMin {
		return []arcSector{{startMin: startMin, endMin: endMin}}
	}

	// 跨天：两段
	res := make([]arcSector, 0, 2)
	if startMin < 1440 {
		res = append(res, arcSector{startMin: startMin, endMin: 1440})
	}
	if endMin > 0 {
		res = append(res, arcSector{startMin: 0, endMin: endMin})
	}
	return res
}

func mergeArcSectors(sectors []arcSector, gapMin int) []arcSector {
	if len(sectors) == 0 {
		return sectors
	}
	sort.Slice(sectors, func(i, j int) bool {
		if sectors[i].startMin == sectors[j].startMin {
			return sectors[i].endMin < sectors[j].endMin
		}
		return sectors[i].startMin < sectors[j].startMin
	})

	merged := make([]arcSector, 0, len(sectors))
	cur := sectors[0]
	for i := 1; i < len(sectors); i++ {
		s := sectors[i]
		if s.startMin <= cur.endMin {
			if s.endMin > cur.endMin {
				cur.endMin = s.endMin
			}
			continue
		}
		gap := s.startMin - cur.endMin
		if gap < gapMin {
			cur.endMin = s.endMin
		} else {
			merged = append(merged, cur)
			cur = s
		}
	}
	merged = append(merged, cur)
	return merged
}

func minuteOffsetToHHMM(min int) string {
	h := min / 60
	m := min % 60
	return two(h) + ":" + two(m)
}

func two(v int) string {
	if v < 10 {
		return "0" + strconv.Itoa(v)
	}
	return strconv.Itoa(v)
}

func parseHHMMToMinutes(hhmm string) (int, bool) {
	parts := strings.Split(hhmm, ":")
	if len(parts) != 2 {
		return 0, false
	}
	h, err1 := strconv.Atoi(parts[0])
	m, err2 := strconv.Atoi(parts[1])
	if err1 != nil || err2 != nil || h < 0 || h > 23 || m < 0 || m > 59 {
		return 0, false
	}
	return h*60 + m, true
}

// 是否不在预测直播时段内
func isOutsideLiveDuration(sectors []*LiveDuration) bool {

	if len(sectors) == 0 {
		return true
	}

	const windowMin = 20

	t := time.Now()
	// 当前时间的分钟偏移
	curMin := t.Hour()*60 + t.Minute()

	for _, s := range sectors {
		startMin, ok1 := parseHHMMToMinutes(s.StartTime)
		endMin, ok2 := parseHHMMToMinutes(s.EndTime)
		if !ok1 || !ok2 || startMin < 0 || endMin > 1440 || startMin >= endMin {
			continue
		}

		// 扩展前后20分钟，并裁剪到 [0,1440]
		winStart := startMin - windowMin
		if winStart < 0 {
			winStart = 0
		}
		winEnd := endMin + windowMin
		if winEnd > 1440 {
			winEnd = 1440
		}

		if curMin >= winStart && curMin <= winEnd {
			return false
		}
	}
	return true
}
