package axios

import (
	"bytes"
	"errors"
	"fmt"
	"io"
	"net/url"
	"sync"
	"time"

	"d1y.io/neovideo/config"
	"github.com/imroc/req/v3"
	"github.com/patrickmn/go-cache"
)

var sn = new(sync.Once)
var iq *iReq

var errFlag = []byte("err{0}")

func getCacheKey(m, u string, qs url.Values) string {
	return fmt.Sprintf("%s%s?=%s", m, u, qs.Encode())
}

type iReq struct {
	Request *req.Client
	cc      *cache.Cache
	sync.Once
}

func GetClient() *iReq {
	sn.Do(initInstance)
	return iq
}

func Request() *req.Request {
	return GetClient().Request.R()
}

func request(api string, qs ...map[string]string) (string, *req.Request, []byte, bool) {
	r := Request()
	if len(qs) == 1 {
		r.SetQueryParams(qs[0])
	}
	q := r.QueryParams.Encode()
	realURL := fmt.Sprintf("%s?=%s", api, q)
	if val, ok := iq.cc.Get(realURL); ok {
		if v, o := val.([]byte); o {
			return realURL, nil, v, true
		}
	}
	return realURL, r, nil, false
}

func handleResponse(resp *req.Response, err error, key string) ([]byte, error) {
	if err != nil {
		return nil, err
	}
	b, e := io.ReadAll(resp.Body)
	if e != nil {
		return nil, e
	}
	if bytes.Equal(b, errFlag) {
		return nil, errors.New("err{0}")
	}
	iq.cc.SetDefault(key, b)
	return b, nil
}

func Get(api string, qs ...map[string]string) ([]byte, error) {
	cacheKey, r, cacheValue, cacheOK := request(api, qs...)
	if cacheOK {
		return cacheValue, nil
	}
	resp, err := r.Get(api)
	return handleResponse(resp, err, cacheKey)
}

func Post(api string, qs ...map[string]string) ([]byte, error) {
	cacheKey, r, cacheValue, cacheOK := request(api, qs...)
	if cacheOK {
		return cacheValue, nil
	}
	resp, err := r.Post(api)
	return handleResponse(resp, err, cacheKey)
}

func initInstance() {
	iq = new(iReq)
	iq.cc = cache.New(42*time.Second, 60*time.Second)
	iq.Request = req.C().
		EnableInsecureSkipVerify().
		OnAfterResponse(func(client *req.Client, resp *req.Response) error {
			if resp.IsSuccessState() {
				key := getCacheKey(resp.Request.Method, resp.Request.RawURL, resp.Request.QueryParams)
				if _, ok := iq.cc.Get(key); !ok {
					iq.cc.SetDefault(key, resp.Bytes())
				}
			}
			return nil
		})
	// FIXME: 单元测试的话这里由于没有初始化会 panic
	if config.Get().IsDev() {
		iq.Request.EnableDebugLog()
	}
}
