package utils

import (
	"context"
	"path/filepath"
	"strings"

	"github.com/gogf/gf/v2/errors/gerror"
	"github.com/gogf/gf/v2/frame/g"
)

var fatalKeywords = []string{
	"no frame",
	"invalid data found",
	"error while decoding",
	"could not find codec parameters",
	"moov atom not found",
	"end of file",
	"truncated",
	"failed to read",
	"malformed",
}

func FileAbsPath(path, filename string) (string, error) {
	base, err := filepath.Abs(DATA_PATH)
	if err != nil {
		return "", gerror.New("获取系统目录信息失败")
	}
	absPath, err := filepath.Abs(filepath.Join(base, path, filename))
	if err != nil {
		return "", gerror.New("获取文件绝对路径失败")
	}
	return absPath, nil
}

func QuickCheckFile(ctx context.Context, absPath string) error {
	fb := NewFFprobeBuilder()
	fb = fb.Input(absPath)
	args := []string{
		"-v", "error",
		"-hide_banner",
		"-print_format", "json",
		"-show_format",
		"-show_streams",
		"-read_intervals", "0%+1,23%+1,50%+1,78%+1,99%+1",
	}
	fb = fb.AddArgs(args...)
	_, stderr, err := fb.Execute(ctx)
	g.Log().Infof(ctx, "QuickCheckFile ffprobe exit error: %s", stderr)
	if err != nil {
		return gerror.Wrap(err, "QuickCheckFile ffprobe exit error")
	}
	if isFatalFFprobeLog(stderr) {
		return gerror.New("QuickCheckFile Failed: fatal probe error")
	}
	return nil
}

func CompletedCheckFile(ctx context.Context, absPath string) error {
	fb := NewFFmpegBuilder()
	fb = fb.AddArgs("-v", "error", "-xerror", "-hide_banner")
	fb = fb.Input(absPath)
	fb = fb.AddArgs("-map", "0")
	fb = fb.CopyCodec()
	fb = fb.AddArgs("-f", "null", "-")

	fb.Execute(ctx)

	return nil
}

func isFatalFFprobeLog(stderr string) bool {
	s := strings.ToLower(stderr)
	for _, kw := range fatalKeywords {
		if strings.Contains(s, strings.ToLower(kw)) {
			return true
		}
	}
	return false
}
