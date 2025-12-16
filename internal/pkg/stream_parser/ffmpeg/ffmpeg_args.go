package ffmpeg_parser

import (
	"net/url"
	"strings"

	"github.com/gogf/gf/v2/util/gconv"
	"github.com/shichen437/gowlive/internal/pkg/utils"
)

func (p *Parser) buildArgs(ffUserAgent, file string, sUrl *url.URL, headers map[string]string) []string {
	referer, exists := headers["Referer"]
	if !exists {
		referer = p.referer
	}
	args := p.basicArgs(ffUserAgent, referer, sUrl)
	for k, v := range headers {
		if k == "User-Agent" || k == "Referer" {
			continue
		}
		args = append(args, "-headers", k+": "+v)
	}
	switch strings.ToLower(p.format) {
	case "flv":
		args = p.flvArgs(file, args)
	case "mp4":
		args = p.mp4Args(file, args)
	case "mkv":
		args = p.mkvArgs(file, args)
	case "ts":
		args = p.tsArgs(file, args)
	default:
		args = p.mp3Args(file, args)
	}
	return args
}

func (p *Parser) basicArgs(ffUserAgent, referer string, sUrl *url.URL) []string {
	return []string{
		"-nostats",
		"-progress", "-",
		"-y", "-re",
		"-fflags", "+genpts+discardcorrupt",
		"-err_detect", "ignore_err",
		"-reconnect", "1",
		"-reconnect_streamed", "1",
		"-reconnect_delay_max", "5",
		"-user_agent", ffUserAgent,
		"-referer", referer,
		"-rw_timeout", p.timeoutInUs,
		"-i", sUrl.String(),
	}
}

func (p *Parser) flvArgs(file string, args []string) []string {
	args = append(args, "-c", "copy")
	if gconv.Int(p.st) > 0 {
		template := utils.BuildSegmentTemplate(file, ".flv")
		args = append(args,
			"-f", "segment",
			"-segment_time", p.st,
			"-reset_timestamps", "1",
			template,
		)
	} else {
		args = append(args, "-f", "flv", utils.EnsureSuffix(file, ".flv"))
	}
	return args
}

func (p *Parser) mp4Args(file string, args []string) []string {
	args = append(args, "-c", "copy", "-bsf:a", "aac_adtstoasc")
	if gconv.Int(p.st) > 0 {
		template := utils.BuildSegmentTemplate(file, ".mp4")
		args = append(args,
			"-f", "segment",
			"-segment_time", p.st,
			"-reset_timestamps", "1",
			"-segment_format_options", "movflags=+faststart",
			template,
		)
	} else {
		args = append(args, "-movflags", "+faststart", "-f", "mp4", utils.EnsureSuffix(file, ".mp4"))
	}
	return args
}

func (p *Parser) mkvArgs(file string, args []string) []string {
	if p.fr {
		args = append(args, "-vf", "scale=1920:1080:force_original_aspect_ratio=decrease,pad=1920:1080:(ow-iw)/2:(oh-ih)/2:black")
		args = append(args, "-c:v", "libx264", "-preset", "fast", "-crf", "23")
		args = append(args, "-c:a", "aac", "-b:a", "128k")
	} else {
		args = append(args, "-c", "copy")
	}

	if gconv.Int(p.st) > 0 {
		template := utils.BuildSegmentTemplate(file, ".mkv")
		args = append(args,
			"-f", "segment",
			"-segment_time", p.st,
			"-reset_timestamps", "1",
			template,
		)
	} else {
		args = append(args, "-f", "matroska", utils.EnsureSuffix(file, ".mkv"))
	}
	return args
}

func (p *Parser) tsArgs(file string, args []string) []string {
	args = append(args, "-c", "copy")

	if gconv.Int(p.st) > 0 {
		template := utils.BuildSegmentTemplate(file, ".ts")
		args = append(args,
			"-f", "segment",
			"-segment_time", p.st,
			"-reset_timestamps", "1",
			template,
		)
	} else {
		args = append(args, "-f", "mpegts", utils.EnsureSuffix(file, ".ts"))
	}
	return args
}

func (p *Parser) mp3Args(file string, args []string) []string {
	args = append(args, "-vn", "-c:a", "libmp3lame", "-b:a", "192k")
	if gconv.Int(p.st) > 0 {
		template := utils.BuildSegmentTemplate(file, ".mp3")
		args = append(args,
			"-f", "segment",
			"-segment_time", p.st,
			"-reset_timestamps", "1",
			template,
		)
	} else {
		args = append(args, "-f", "mp3", utils.EnsureSuffix(file, ".mp3"))
	}
	return args
}
