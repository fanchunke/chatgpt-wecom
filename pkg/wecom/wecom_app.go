package wecom

import (
	"encoding/json"
	"net/http"
	"sync"
	"time"

	"github.com/fanchunke/chatgpt-wecom/pkg/wecom/callback"

	"github.com/go-resty/resty/v2"
	"github.com/rs/zerolog/log"
)

type WeComApp struct {
	*WeCom
	// 应用 AgentId
	agentId int64
	// 应用 Secret
	corpSecret string
	// 接收消息 EncodingAESKey
	encodingAESKey string
	// 接收消息 Token
	token string
	// accessToken
	accessToken *token
	// accessToken Manager
	accessTokenManager TokenManager
}

func NewApp(baseURL, corpId string, agentId int64, corpSecret, encodingAESKey, appToken string, timeout time.Duration) *WeComApp {
	wecom := New(baseURL, corpId, timeout)
	app := &WeComApp{
		WeCom:          wecom,
		encodingAESKey: encodingAESKey,
		token:          appToken,
		agentId:        agentId,
		corpSecret:     corpSecret,
		accessToken:    &token{mutex: &sync.RWMutex{}},
	}
	app.accessTokenManager = &token{
		mutex:        &sync.RWMutex{},
		getTokenFunc: app.getToken,
	}
	// app.AddRetryCondition(app.defaultRetryCondition())
	app.rc.SetRetryCount(1)
	return app
}

func (app *WeComApp) RxMessageHandler(h callback.RxMessageHandler) (http.Handler, error) {
	return callback.NewHTTPHandler(app.token, app.encodingAESKey, h)
}

func (app *WeComApp) defaultRetryCondition() resty.RetryConditionFunc {
	return func(r *resty.Response, err error) bool {
		var resp RespCommon
		if err := json.Unmarshal(r.Body(), &resp); err != nil {
			return false
		}
		if resp.IsAccessTokenError() {
			ctx := r.Request.Context()
			log.Ctx(ctx).Info().Msg("Access token 错误，开始重新获取后重试")
			token, err := app.accessTokenManager.Refresh(ctx)
			if err != nil {
				log.Ctx(ctx).Error().Err(err).Msg("Refresh Access Token 错误")
				return false
			}
			// r.Request.QueryParam.Add("access_token", token)
			r.Request.SetQueryParam("access_token", token.Token)
			return true
		}
		return false
	}
}

func (app *WeComApp) SetAccessTokenManager(m TokenManager) {
	app.accessTokenManager = m
}

func (app *WeComApp) GetAgentId() int64 {
	return app.agentId
}
