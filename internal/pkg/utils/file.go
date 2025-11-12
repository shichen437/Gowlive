package utils

import (
	"context"
	"path/filepath"

	"github.com/gogf/gf/v2/errors/gerror"
)

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
	if err != nil || stderr != "" {
		return gerror.New("QuickCheckFile Failed")
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
