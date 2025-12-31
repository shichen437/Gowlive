package openlist

import (
	"bufio"
	"context"
	"io"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	"github.com/gogf/gf/v2/encoding/gjson"
	"github.com/gogf/gf/v2/errors/gerror"
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/net/gclient"
	"github.com/gogf/gf/v2/os/gcron"
	"github.com/gogf/gf/v2/os/gctx"
	"github.com/shichen437/gowlive/internal/pkg/utils"
)

func Login() error {
	ctx := gctx.GetInitCtx()
	c := g.Client()
	c.SetContentType("application/json")
	resp, err := c.Post(ctx, utils.OPENLIST_DOMAIN+loginUrl, g.Map{
		"username": utils.OPENLIST_USERNAME,
		"password": utils.OPENLIST_PASSWORD,
		"otp_code": utils.OPENLIST_CODE,
	})
	_status = 400
	if err != nil || resp.StatusCode != http.StatusOK {
		err = gerror.Newf("failed to login, err : %v, code: %d", err, resp.StatusCode)
		return err
	}
	body, err := utils.Text(resp.Response)
	if err != nil {
		return err
	}
	var loginModel *LoginRespModel
	err = gjson.Unmarshal([]byte(body), &loginModel)
	if err != nil {
		return gerror.New("Openlist login parse response failed")
	}
	if loginModel.Data.Token == "" {
		return gerror.New("Openlist login error: " + loginModel.Message)
	}
	_token = loginModel.Data.Token
	_status = 200

	if job := gcron.Search("openlist-login-retry"); job != nil {
		g.Log().Info(ctx, "Openlist login successful, removing scheduled retry job.")
		gcron.Remove("openlist-login-retry")
	}

	return nil
}

func Username() (string, error) {
	resp, err := withRetry(func(ctx context.Context) (*gclient.Response, error) {
		return getClientWithToken().Get(ctx, utils.OPENLIST_DOMAIN+userInfoUrl)
	})
	if err != nil {
		return "", err
	}
	defer resp.Close()

	if resp.StatusCode != http.StatusOK {
		return "", gerror.Newf("failed to get username, code: %d", resp.StatusCode)
	}

	body, err := utils.Text(resp.Response)
	if err != nil {
		return "", err
	}
	var userInfoModel *UserInfoRespModel
	if err := gjson.Unmarshal([]byte(body), &userInfoModel); err != nil {
		return "", gerror.New("Openlist get username parse response failed")
	}
	if userInfoModel.Data.ID == "" {
		return "", gerror.New("Openlist get username error: " + userInfoModel.Message)
	}
	return userInfoModel.Data.Username, nil
}

func DirInfo(path string) (*DirContentModel, error) {
	resp, err := withRetry(func(ctx context.Context) (*gclient.Response, error) {
		return getClientWithToken().Post(ctx, utils.OPENLIST_DOMAIN+dirInfoUrl, map[string]any{
			"path":     path,
			"password": "",
		})
	})
	if err != nil {
		return nil, err
	}
	defer resp.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, gerror.Newf("failed to list dir, code: %d", resp.StatusCode)
	}

	body, err := utils.Text(resp.Response)
	if err != nil {
		return nil, err
	}
	var dirInfo *DirInfoRespModel
	if err := gjson.Unmarshal([]byte(body), &dirInfo); err != nil {
		return nil, gerror.New("Openlist list dir parse response failed")
	}
	if dirInfo.Code != 200 {
		return nil, gerror.New("Openlist list dir error: " + dirInfo.Message)
	}
	if dirInfo.Data == nil {
		return nil, nil
	}
	return dirInfo.Data, nil
}

func Upload(root, path, filename string) error {
	localFilePath := filepath.Join(path, filename)
	file, err := os.Open(localFilePath)
	if err != nil {
		return gerror.Wrapf(err, "failed to open file: %s", localFilePath)
	}
	defer file.Close()
	dirs := strings.Split(path, "stream/")
	var uploadPath string
	if len(dirs) > 1 {
		uploadPath = filepath.ToSlash(filepath.Join(root, "/stream/", dirs[1]))
	} else {
		uploadPath = filepath.ToSlash(filepath.Join(root, path))
	}

	mkdir(uploadPath)

	fileInfo, err := file.Stat()
	if err != nil {
		return gerror.Wrapf(err, "failed to get file info: %s", localFilePath)
	}
	bufReader := bufio.NewReaderSize(file, 10*1024*1024)

	filename = utils.ReplaceColonWithDash(filename)

	req, err := http.NewRequest("PUT", utils.OPENLIST_DOMAIN+uploadFileStreamUrl, bufReader)
	if err != nil {
		return err
	}
	req.Header.Add("Authorization", _token)
	req.Header.Add("File-Path", "/"+uploadPath+"/"+filename)
	req.Header.Add("As-Task", "false")
	req.Header.Add("Content-Length", strconv.FormatInt(fileInfo.Size(), 10))
	req.Header.Add("Content-Type", "application/octet-stream")

	c := g.Client()
	c.SetTimeout(time.Hour * 24 * 30)
	resp, err := c.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return gerror.Newf("failed to upload stream, code: %d", resp.StatusCode)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return err
	}
	var uploadResult *CommonResp
	if err := gjson.Unmarshal([]byte(body), &uploadResult); err != nil {
		return gerror.New("Openlist upload stream parse response failed")
	}
	if uploadResult.Code != 200 {
		return gerror.New("Openlist upload stream error: " + uploadResult.Message)
	}
	return nil
}

func mkdir(path string) error {
	resp, err := withRetry(func(ctx context.Context) (*gclient.Response, error) {
		return getClientWithToken().Post(ctx, utils.OPENLIST_DOMAIN+mkdirUrl, map[string]any{
			"path": "/" + path,
		})
	})
	if err != nil {
		return err
	}
	defer resp.Close()
	if resp.StatusCode != http.StatusOK {
		return gerror.Newf("failed to create dir, code: %d", resp.StatusCode)
	}

	body, err := utils.Text(resp.Response)
	if err != nil {
		return err
	}
	var uploadResult *CommonResp
	if err := gjson.Unmarshal([]byte(body), &uploadResult); err != nil {
		return gerror.New("Openlist create dir parse response failed")
	}
	if uploadResult.Code != 200 {
		return gerror.New("Openlist create dir error: " + uploadResult.Message)
	}
	return nil
}

func Logout() {
	getClientWithToken().Get(gctx.GetInitCtx(), logoutUrl)
	_token = ""
	_status = 0
}

func GetStatus() int {
	return _status
}

func CheckLoginParams() bool {
	if utils.OPENLIST_DOMAIN == "" || utils.OPENLIST_USERNAME == "" || utils.OPENLIST_PASSWORD == "" {
		_status = 0
		return false
	}
	_, err := url.Parse(utils.OPENLIST_DOMAIN)
	if err != nil {
		return false
	}
	return true
}
