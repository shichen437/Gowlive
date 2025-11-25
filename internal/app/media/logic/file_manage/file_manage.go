package logic

import (
	"context"
	"fmt"
	"io/fs"
	"mime"
	"net/http"
	"os"
	"path/filepath"
	"strings"

	"github.com/gogf/gf/v2/errors/gerror"
	"github.com/gogf/gf/v2/net/ghttp"
	v1 "github.com/shichen437/gowlive/api/v1/media"
	"github.com/shichen437/gowlive/internal/app/media/model"
	"github.com/shichen437/gowlive/internal/app/media/service"
	"github.com/shichen437/gowlive/internal/pkg/utils"
)

type (
	sFileManage struct{}
)

func init() {
	service.RegisterFileManage(New())
}

func New() service.IFileManage {
	return &sFileManage{}
}

func (s *sFileManage) List(ctx context.Context, req *v1.GetFileListReq) (res *v1.GetFileListRes, err error) {
	res = &v1.GetFileListRes{}
	absPath, err := utils.FileAbsPath(req.Path, "")
	if err != nil {
		return nil, gerror.New("获取系统目录路径失败")
	}
	files, err := os.ReadDir(absPath)
	if err != nil {
		return nil, gerror.New("获取目录列表失败")
	}
	if len(files) == 0 {
		return
	}
	var list []*model.FileInfo
	for _, file := range files {
		info, err := file.Info()
		if err != nil || isHiddenFile(info) {
			continue
		}
		if req.Filename != "" && !matchPattern(file.Name(), req.Filename) {
			continue
		}
		list = append(list, &model.FileInfo{
			Filename:     file.Name(),
			IsFolder:     file.IsDir(),
			Size:         info.Size(),
			LastModified: info.ModTime().Local().UnixMilli(),
		})
	}
	res.Rows = list
	return
}

func (s *sFileManage) Delete(ctx context.Context, req *v1.DeleteFileReq) (res *v1.DeleteFileRes, err error) {
	res = &v1.DeleteFileRes{}
	if req.Path == "" && strings.Contains(req.Filename, "gowlive.db") {
		return
	}
	absPath, err := utils.FileAbsPath(req.Path, "")
	if err != nil {
		return nil, gerror.New("获取系统目录路径失败")
	}
	err = os.RemoveAll(filepath.Join(absPath, req.Filename))
	if err != nil {
		return nil, gerror.New("删除文件失败")
	}
	return
}

func (s *sFileManage) Empty(ctx context.Context, req *v1.GetEmptyFolderReq) (res *v1.GetEmptyFolderRes, err error) {
	res = &v1.GetEmptyFolderRes{}
	absPath, err := utils.FileAbsPath(req.Path, "")
	if err != nil {
		return nil, gerror.New("获取系统目录路径失败")
	}
	res.IsEmpty = !utils.HasAnyFile(absPath)
	return
}

func (s *sFileManage) AnchorFilePath(ctx context.Context, req *v1.GetAnchorFilePathReq) (res *v1.GetAnchorFilePathRes, err error) {
	res = &v1.GetAnchorFilePathRes{}
	if req.Anchor == "" || req.Platform == "" {
		return
	}
	cachePath := "stream/" + req.Platform
	absPath, err := utils.FileAbsPath(cachePath, "")
	if err != nil {
		return res, nil
	}
	res.Path = cachePath
	cachePath = res.Path + "/" + req.Anchor
	absPath, err = utils.FileAbsPath(cachePath, "")
	if err != nil {
		return res, nil
	}
	_, err = os.Stat(absPath)
	if err != nil {
		return res, nil
	}
	res.Path = cachePath
	return
}

func (s *sFileManage) Play(ctx context.Context, req *v1.GetFilePlayReq) (res *v1.GetFilePlayRes, err error) {
	r := ghttp.RequestFromCtx(ctx)
	if r == nil {
		return nil, gerror.New("invalid request context")
	}
	res = &v1.GetFilePlayRes{}

	setupCORS(r)
	if r.Method == http.MethodOptions {
		r.Response.WriteHeader(http.StatusNoContent)
		return res, nil
	}

	path := strings.TrimSpace(req.Path)
	if path == "" {
		writeErrorPlain(r, http.StatusBadRequest, "文件路径不能为空")
		return
	}

	abs, errAbsJoin := utils.FileAbsPath(path, "")
	if errAbsJoin != nil {
		writeErrorPlain(r, http.StatusBadRequest, "非法文件路径")
		return
	}

	fi, statErr := os.Stat(abs)
	if statErr != nil {
		if os.IsNotExist(statErr) {
			writeErrorPlain(r, http.StatusNotFound, "文件不存在")
		} else {
			writeErrorPlain(r, http.StatusInternalServerError, "无法访问文件")
		}
		return
	}
	if fi.IsDir() {
		writeErrorPlain(r, http.StatusBadRequest, "路径为目录，无法播放")
		return
	}

	ctype := detectContentType(abs)
	if !isSupportedMedia(ctype) {
		writeErrorPlain(r, http.StatusUnsupportedMediaType, fmt.Sprintf("不支持的媒体类型: %s", ctype))
		return
	}

	f, openErr := os.Open(abs)
	if openErr != nil {
		writeErrorPlain(r, http.StatusInternalServerError, "文件打开失败")
		return
	}
	defer f.Close()

	r.Response.Header().Set("Content-Type", ctype)
	r.Response.Header().Set("Accept-Ranges", "bytes")

	r.Response.Header().Set("Cache-Control", "no-store")

	modTime := fi.ModTime()
	if !modTime.IsZero() {
		r.Response.Header().Set("Last-Modified", modTime.UTC().Format(http.TimeFormat))
	}

	http.ServeContent(r.Response.Writer, r.Request, fi.Name(), modTime, f)

	return res, nil
}

func isHiddenFile(file fs.FileInfo) bool {
	if file.IsDir() {
		return false
	}
	return strings.HasPrefix(file.Name(), ".")
}

func matchPattern(filename, pattern string) bool {
	return strings.Contains(filename, pattern)
}

func detectContentType(path string) string {
	ext := strings.ToLower(filepath.Ext(path))
	switch ext {
	case ".mp4":
		return "video/mp4"
	case ".mp3":
		return "audio/mpeg"
	case ".flv":
		return "video/x-flv"
	case ".mkv":
		return "video/x-matroska"
	case ".ts":
		return "video/mp2t"
	default:
		ctype := mime.TypeByExtension(ext)
		if ctype == "" {
			ctype = "application/octet-stream"
		}
		return ctype
	}
}

func isSupportedMedia(ctype string) bool {
	if strings.HasPrefix(ctype, "video/mp4") {
		return true
	}
	if strings.HasPrefix(ctype, "audio/mpeg") {
		return true
	}
	if strings.HasPrefix(ctype, "video/x-matroska") {
		return true
	}
	if strings.HasPrefix(ctype, "video/mp2t") {
		return true
	}
	if ctype == "video/x-flv" || strings.HasPrefix(ctype, "video/flv") {
		return true
	}
	return false
}

func writeErrorPlain(r *ghttp.Request, code int, msg string) {
	r.Response.ClearBuffer()
	r.Response.WriteHeader(code)
	_, _ = r.Response.Writer.Write([]byte(msg))
	r.Exit()
}

func setupCORS(r *ghttp.Request) {
	origin := r.Header.Get("Origin")
	if origin == "" {
		return
	}
	h := r.Response.Header()
	h.Set("Vary", "Origin")
	h.Set("Access-Control-Allow-Origin", origin)
	h.Set("Access-Control-Allow-Credentials", "true")
	h.Set("Access-Control-Allow-Headers", "Authorization, Range, Content-Type")
	h.Set("Access-Control-Allow-Methods", "GET, POST, OPTIONS, HEAD, CONNECT, TRACE")
	h.Set("Access-Control-Expose-Headers", "Content-Range, Accept-Ranges, Content-Length, Last-Modified")
}
