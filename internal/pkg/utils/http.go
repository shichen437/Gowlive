package utils

import (
	"context"
	"errors"
	"io"
	"net/http"
	"net/http/cookiejar"
	"net/url"
	"strings"
	"time"

	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gctx"
)

func Text(r *http.Response) (string, error) {
	if r.Body == nil {
		return "", io.EOF
	}

	defer func() {
		if err := r.Body.Close(); err != nil {
			g.Log().Error(gctx.GetInitCtx(), "Error closing response body:", err)
		}
	}()

	b, err := io.ReadAll(r.Body)
	if err != nil {
		return "", err
	}
	return string(b), nil
}

func GetCookieList(cookie string) []*http.Cookie {
	cookiesList := make([]*http.Cookie, 0)
	if cookie == "" {
		return cookiesList
	}
	for cStr := range strings.SplitSeq(cookie, ";") {
		cArr := strings.SplitN(cStr, "=", 2)
		if len(cArr) != 2 {
			continue
		}
		cookiesList = append(cookiesList, &http.Cookie{
			Name:  strings.TrimSpace(cArr[0]),
			Value: strings.TrimSpace(cArr[1]),
		})
	}
	return cookiesList
}

func GetCookieMap(platform, cookie, refer string) map[string]string {
	jar, _ := cookiejar.New(&cookiejar.Options{})
	url, _ := url.Parse(refer)
	jar.SetCookies(url, GetCookieList(cookie))
	cookies := jar.Cookies(url)
	cookieMap := make(map[string]string)
	for _, c := range cookies {
		cookieMap[c.Name] = c.Value
	}
	return cookieMap
}

func UrlRemoveParams(raw string) string {
	head, _, found := strings.Cut(raw, "?")
	if found {
		return head
	}
	return raw
}

func GetRedirectURL(ctx context.Context, shortURL string) (string, error) {
	if shortURL == "" {
		return "", errors.New("empty url")
	}

	const maxRedirects = 5

	client := &http.Client{
		Timeout: 15 * time.Second,
		CheckRedirect: func(req *http.Request, via []*http.Request) error {
			return http.ErrUseLastResponse
		},
	}

	currentURL := shortURL
	for range maxRedirects {
		req, err := http.NewRequestWithContext(ctx, http.MethodGet, currentURL, nil)
		if err != nil {
			return "", err
		}

		resp, err := client.Do(req)
		if err != nil {
			return "", err
		}

		if resp.StatusCode >= 300 && resp.StatusCode < 400 {
			loc := resp.Header.Get("Location")
			if cerr := resp.Body.Close(); cerr != nil {
				g.Log().Warning(gctx.GetInitCtx(), "failed to close response body:", cerr)
			}
			if loc == "" {
				return "", errors.New("redirect response missing Location header")
			}
			currentURL = loc
			continue
		}

		if cerr := resp.Body.Close(); cerr != nil {
			g.Log().Warning(gctx.GetInitCtx(), "failed to close response body:", cerr)
		}
		return currentURL, nil
	}

	return "", errors.New("too many redirects")
}
