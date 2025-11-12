package utils

import (
	"bytes"
	"context"
	"os/exec"

	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gctx"
)

type FFprobeBuilder struct {
	ffprobePath string
	globalArgs  []string
	inputPath   string
}

func NewFFprobeBuilder() *FFprobeBuilder {
	path, err := GetDefaultFFprobePath()
	if err != nil {
		g.Log().Error(gctx.New(), err)
		path = "ffprobe"
	}
	return &FFprobeBuilder{
		ffprobePath: path,
	}
}

func (b *FFprobeBuilder) Input(absPath string) *FFprobeBuilder {
	b.inputPath = absPath
	return b
}

func (b *FFprobeBuilder) AddArg(arg string) *FFprobeBuilder {
	b.globalArgs = append(b.globalArgs, arg)
	return b
}

func (b *FFprobeBuilder) AddArgs(args ...string) *FFprobeBuilder {
	b.globalArgs = append(b.globalArgs, args...)
	return b
}

func (b *FFprobeBuilder) BuildArgs() []string {
	var args []string
	args = append(args, b.globalArgs...)
	args = append(args, b.inputPath)
	return args
}

func (b *FFprobeBuilder) Build(ctx context.Context) *exec.Cmd {
	return exec.CommandContext(ctx, b.ffprobePath, b.BuildArgs()...)
}

func (b *FFprobeBuilder) Execute(ctx context.Context) (string, string, error) {
	var stdout, stderr bytes.Buffer
	cmd := b.Build(ctx)
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr
	err := cmd.Run()
	return stdout.String(), stderr.String(), err
}
