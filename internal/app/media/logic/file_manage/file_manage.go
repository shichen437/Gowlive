package logic

import (
	"context"
	"io/fs"
	"os"
	"path/filepath"
	"strings"

	"github.com/gogf/gf/v2/errors/gerror"
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
	base, err := filepath.Abs(utils.DATA_PATH)
	if err != nil {
		return nil, gerror.New("获取系统目录信息失败")
	}
	absPath, err := filepath.Abs(filepath.Join(base, req.Path))
	if err != nil || !strings.HasPrefix(absPath, base) {
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
	if req.Filename == "" || strings.Contains(req.Filename, "gowlive.db") {
		return
	}
	base, err := filepath.Abs(utils.DATA_PATH)
	if err != nil {
		return nil, gerror.New("获取系统目录信息失败")
	}
	absPath, err := filepath.Abs(filepath.Join(base, req.Path))
	if err != nil || !strings.HasPrefix(absPath, base) {
		return nil, gerror.New("获取系统目录路径失败")
	}
	err = os.RemoveAll(filepath.Join(absPath, req.Filename))
	if err != nil {
		return nil, gerror.New("删除文件失败")
	}
	return
}

func (s *sFileManage) Play(ctx context.Context, req *v1.GetFilePlayReq) (res *v1.GetFilePlayRes, err error) {
	return
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
