package utils

import (
	"github.com/gogf/gf/v2/net/gclient"
)

func SafeStatusCode(resp *gclient.Response) (code int) {
	defer func() {
		if rec := recover(); rec != nil {
			code = 0
		}
	}()
	code = resp.StatusCode
	return code
}

func SafeResponseBody(resp *gclient.Response) (body string, err error) {
	defer func() {
		if rec := recover(); rec != nil {
			body = ""
			err = nil
		}
	}()
	body, err = Text(resp.Response)
	return body, err
}
