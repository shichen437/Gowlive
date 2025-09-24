package utils

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"os/exec"
	"runtime"
	"strings"

	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gctx"
)

var (
	defaultUserAgent = "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/134.0.0.0 Safari/537.36"
)

type FFmpegBuilder struct {
	inputs       []string
	inputOptions []string
	outputs      []string
	globalArgs   []string
	headers      map[string]string
	progress     bool
	ffmpegPath   string
}

func NewFFmpegBuilder() *FFmpegBuilder {
	path, err := GetDefaultFFmpegPath()
	if err != nil {
		g.Log().Error(gctx.New(), err)
		path = "ffmpeg"
	}
	return &FFmpegBuilder{
		headers:    make(map[string]string),
		ffmpegPath: path,
	}
}

func (b *FFmpegBuilder) Input(input string) *FFmpegBuilder {
	b.inputs = append(b.inputs, input)
	return b
}

func (b *FFmpegBuilder) InputOption(option string) *FFmpegBuilder {
	b.inputOptions = append(b.inputOptions, option)
	return b
}

func (b *FFmpegBuilder) Output(output string) *FFmpegBuilder {
	b.outputs = append(b.outputs, output)
	return b
}

func (b *FFmpegBuilder) Codec(codecType, codec string) *FFmpegBuilder {
	b.globalArgs = append(b.globalArgs, fmt.Sprintf("-c:%s", codecType), codec)
	return b
}

func (b *FFmpegBuilder) VideoCodec(codec string) *FFmpegBuilder {
	return b.Codec("v", codec)
}

func (b *FFmpegBuilder) AudioCodec(codec string) *FFmpegBuilder {
	return b.Codec("a", codec)
}

func (b *FFmpegBuilder) CopyCodec() *FFmpegBuilder {
	return b.VideoCodec("copy").AudioCodec("copy")
}

func (b *FFmpegBuilder) FastStart() *FFmpegBuilder {
	b.globalArgs = append(b.globalArgs, "-movflags", "+faststart")
	return b
}

func (b *FFmpegBuilder) Overwrite() *FFmpegBuilder {
	b.globalArgs = append(b.globalArgs, "-y")
	return b
}

func (b *FFmpegBuilder) AddHeader(key, value string) *FFmpegBuilder {
	b.headers[key] = value
	return b
}

func (b *FFmpegBuilder) AddDefaultUserAgent() *FFmpegBuilder {
	return b.AddUserAgent(defaultUserAgent)
}

func (b *FFmpegBuilder) AddUserAgent(userAgent string) *FFmpegBuilder {
	return b.AddHeader("User-Agent", userAgent)
}

func (b *FFmpegBuilder) AddReferer(referer string) *FFmpegBuilder {
	return b.AddHeader("Referer", referer)
}

func (b *FFmpegBuilder) AddCookie(cookie string) *FFmpegBuilder {
	return b.AddHeader("Cookie", cookie)
}

func (b *FFmpegBuilder) AddDefaultThreads() *FFmpegBuilder {
	cpuCount := runtime.NumCPU()
	threads := cpuCount / 2
	if threads < 1 {
		threads = 1
	} else if threads > 4 {
		threads = 4
	}
	return b.AddThreads(threads)
}

func (b *FFmpegBuilder) AddThreads(threads int) *FFmpegBuilder {
	b.globalArgs = append(b.globalArgs, "-threads", fmt.Sprintf("%d", threads))
	return b
}

func (b *FFmpegBuilder) AddCpuUsage(usage string) *FFmpegBuilder {
	b.globalArgs = append(b.globalArgs, "-preset:v", usage)
	return b
}

func (b *FFmpegBuilder) ShowProgress() *FFmpegBuilder {
	b.progress = true
	return b
}

func (b *FFmpegBuilder) AddArg(arg string) *FFmpegBuilder {
	b.globalArgs = append(b.globalArgs, arg)
	return b
}

func (b *FFmpegBuilder) AddArgs(args ...string) *FFmpegBuilder {
	b.globalArgs = append(b.globalArgs, args...)
	return b
}

func (b *FFmpegBuilder) buildArgs() []string {
	var args []string

	// 添加 HTTP 头
	if len(b.headers) > 0 {
		var headers []string
		for k, v := range b.headers {
			headers = append(headers, fmt.Sprintf("%s: %s", k, v))
		}
		headerStr := strings.Join(headers, "\r\n")
		args = append(args, "-headers", headerStr+"\r\n")
	}

	// 添加输入选项
	args = append(args, b.inputOptions...)

	// 添加输入文件
	for _, input := range b.inputs {
		args = append(args, "-i", input)
	}

	// 添加全局参数
	args = append(args, b.globalArgs...)

	// 添加进度参数
	if b.progress {
		args = append(args, "-progress", "pipe:1")
	}

	// 添加输出文件
	args = append(args, b.outputs...)

	return args
}

// Build 构建 FFmpeg 命令
func (b *FFmpegBuilder) Build() *exec.Cmd {
	args := b.buildArgs()
	return exec.Command(b.ffmpegPath, args...)
}

// BuildWithContext 构建带有上下文的 FFmpeg 命令
func (b *FFmpegBuilder) BuildWithContext(ctx context.Context) *exec.Cmd {
	args := b.buildArgs()
	return exec.CommandContext(ctx, b.ffmpegPath, args...)
}

// String 返回完整的命令字符串
func (b *FFmpegBuilder) String() string {
	args := b.buildArgs()
	var buffer bytes.Buffer
	buffer.WriteString(b.ffmpegPath)
	for _, arg := range args {
		buffer.WriteString(" ")
		if strings.Contains(arg, " ") {
			buffer.WriteString("\"")
			buffer.WriteString(arg)
			buffer.WriteString("\"")
		} else {
			buffer.WriteString(arg)
		}
	}
	return buffer.String()
}

// Execute 执行 FFmpeg 命令
func (b *FFmpegBuilder) Execute(ctx context.Context) ([]byte, error) {
	cmd := b.Build()
	g.Log().Debug(ctx, "执行FFmpeg命令:", b.String())
	output, err := cmd.CombinedOutput()

	if err != nil {
		lastOutput := output
		if len(output) > 500 {
			lastOutput = output[len(output)-500:]
		}
		g.Log().Errorf(ctx, "FFmpeg命令执行失败: %v, 输出: %s", err, string(lastOutput))
		return output, fmt.Errorf("FFmpeg执行失败: %w, 输出: %s", err, string(lastOutput))
	}

	return output, nil
}

// ExecuteWithProgress 执行 FFmpeg 命令并返回进度读取器
func (b *FFmpegBuilder) ExecuteWithProgress(ctx context.Context) (io.ReadCloser, error) {
	cmd := b.BuildWithContext(ctx)
	stdout, err := cmd.StdoutPipe()
	if err != nil {
		return nil, err
	}

	if err := cmd.Start(); err != nil {
		return nil, err
	}

	// 在后台等待命令完成
	go func() {
		_ = cmd.Wait()
	}()

	return stdout, nil
}
