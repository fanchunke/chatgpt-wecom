package wecom

import (
	"encoding/json"
	"errors"
	"fmt"
	"time"

	"github.com/go-resty/resty/v2"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

const (
	defaultBaseURL = "https://qyapi.weixin.qq.com"
	defaultTimeout = time.Second * 5
)

type WeCom struct {
	rc     *resty.Client
	corpId string
	logger *zerolog.Logger
}

func New(baseURL, corpId string, timeout time.Duration) *WeCom {
	rc := resty.New()
	rc.SetBaseURL(baseURL)
	rc.SetTimeout(timeout)

	client := &WeCom{
		rc:     rc,
		corpId: corpId,
	}
	client.OnBeforeRequest(client.defaultRequestMiddleware())
	client.OnAfterResponse(client.defaultResponseMiddleware())
	client.OnError(client.defaultErrorHook())
	return client
}

func Default(corpId string) *WeCom {
	return New(defaultBaseURL, corpId, defaultTimeout)
}

func (w *WeCom) Logger() *zerolog.Logger {
	l := w.logger
	if l == nil {
		l = &log.Logger
	}
	return l
}

func (w *WeCom) SetLogger(logger *zerolog.Logger) {
	w.logger = logger
}

func (w *WeCom) defaultRequestMiddleware() resty.RequestMiddleware {
	return func(c *resty.Client, r *resty.Request) error {
		w.Logger().Debug().Msgf("[Request] URL: %s, Method: %s, Query: %s, Form: %s, Body: %v, Headers: %s", r.URL, r.Method, r.QueryParam, r.FormData, r.Body, r.Header)
		return nil
	}
}

func (w *WeCom) OnBeforeRequest(m resty.RequestMiddleware) {
	w.rc.OnBeforeRequest(m)
}

func (w *WeCom) defaultResponseMiddleware() resty.ResponseMiddleware {
	return func(c *resty.Client, r *resty.Response) error {
		w.Logger().Debug().Msgf("[Response] Status: %d, Body: %v, Headers; %s", r.StatusCode(), r, r.Header())
		// TODO: 展示具体的错误信息
		if r.IsError() {
			return errors.New("请求失败")
		}

		var resp RespCommon
		if err := json.Unmarshal(r.Body(), &resp); err != nil {
			return errors.New("响应不是合法的 JSON")
		}
		if resp.IsError() {
			return fmt.Errorf("响应失败，错误码：%d，错误原因：%s", resp.Errcode, resp.Errmsg)
		}
		return nil
	}
}

func (w *WeCom) OnAfterResponse(m resty.ResponseMiddleware) {
	w.rc.OnAfterResponse(m)
}

func (w *WeCom) defaultErrorHook() resty.ErrorHook {
	return func(r *resty.Request, err error) {
		w.Logger().Error().Err(err).Msgf("请求失败")
	}
}

func (w *WeCom) OnError(h resty.ErrorHook) {
	w.rc.OnError(h)
}

func (w *WeCom) AddRetryCondition(condition resty.RetryConditionFunc) {
	w.rc.AddRetryCondition(condition)
}

func (w *WeCom) WithApp(encodingAESKey, token string) *WeComApp {
	return &WeComApp{
		WeCom:          w,
		encodingAESKey: encodingAESKey,
		token:          token,
	}
}

func (w *WeCom) Healthz() (interface{}, error) {
	resp, err := w.rc.R().
		EnableTrace().
		Get("/healthz")
	return resp, err
}
