package utils

import (
	"errors"
	"os/exec"
	"strconv"
	"strings"
	"time"

	"github.com/shirou/gopsutil/v3/disk"
)

func GetDiskUsage() int {
	diskInfo, err := disk.Usage("/")
	if err != nil {
		return 1
	}
	return int(diskInfo.UsedPercent)
}

func GetDiskFreeGBInt() int {
	usage, err := disk.Usage("/")
	if err != nil {
		return -1
	}
	const GB = 1024 * 1024 * 1024
	return int(usage.Free / GB)
}

func GetDefaultFFmpegPath() (string, error) {
	path, err := exec.LookPath("ffmpeg")
	if errors.Is(err, exec.ErrDot) {
		path, err = exec.LookPath("./ffmpeg")
	}
	return path, err
}

func GetDefaultFFprobePath() (string, error) {
	path, err := exec.LookPath("ffprobe")
	if errors.Is(err, exec.ErrDot) {
		path, err = exec.LookPath("./ffprobe")
	}
	return path, err
}

func GetOutputPath() string {
	return STREAM_PATH
}

func GetDownloadPath() string {
	return STREAM_PATH + "/download"
}

func GetTempDownloadPath() string {
	return "./resources/temp/download"
}

func IsTimeRange(st, et string) bool {
	startTime, endTime, ok := validAndGetTime(st, et)
	if !ok {
		return false
	}
	// 获取当前时间
	now := time.Now()
	currentTime := time.Date(now.Year(), now.Month(), now.Day(), now.Hour(), now.Minute(), 0, 0, now.Location())
	if startTime.After(endTime) {
		// 如果结束时间在开始时间之前，则认为是跨天
		return (currentTime.After(startTime) && currentTime.Before(endTime.AddDate(0, 0, 1))) ||
			(currentTime.Before(startTime) && currentTime.Before(endTime))
	} else {
		// 正常情况
		return currentTime.After(startTime) && currentTime.Before(endTime)
	}
}

func validAndGetTime(st string, et string) (time.Time, time.Time, bool) {
	if st == "" || et == "" {
		return time.Time{}, time.Time{}, false
	}
	sArr := strings.Split(st, ":")
	eArr := strings.Split(et, ":")
	if len(sArr) != 2 || len(eArr) != 2 {
		return time.Time{}, time.Time{}, false
	}
	sh, err1 := strconv.Atoi(sArr[0])
	sm, err2 := strconv.Atoi(sArr[1])
	eh, err3 := strconv.Atoi(eArr[0])
	em, err4 := strconv.Atoi(eArr[1])
	if err1 != nil || err2 != nil || err3 != nil || err4 != nil {
		return time.Time{}, time.Time{}, false
	}
	now := time.Now()
	startTime := time.Date(now.Year(), now.Month(), now.Day(), sh, sm, 0, 0, now.Location())
	endTime := time.Date(now.Year(), now.Month(), now.Day(), eh, em, 0, 0, now.Location())
	return startTime, endTime, true
}
