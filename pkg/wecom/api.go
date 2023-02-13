package wecom

import (
	"bytes"
	"context"
	"fmt"
	"path"

	"github.com/fanchunke/chatgpt-wecom/pkg/wecom/errno"
	"github.com/fanchunke/chatgpt-wecom/pkg/wecom/message"

	"github.com/go-resty/resty/v2"
)

type RespCommon struct {
	Errcode int64  `json:"errcode"`
	Errmsg  string `json:"errmsg"`
}

func (r *RespCommon) IsError() bool {
	return r.Errcode != 0
}

func (r *RespCommon) IsAccessTokenError() bool {
	// https://developer.work.weixin.qq.com/document/path/90313
	codes := []int64{
		// 不合法的access_token
		40014,
		// 缺少access_token参数
		41001,
		// access_token已过期
		42001,
	}
	for _, code := range codes {
		if r.Errcode == code {
			return true
		}
	}
	return false
}

type RespAccessToken struct {
	RespCommon
	AccessToken string `json:"access_token"`
	ExpiresIn   int64  `json:"expires_in"`
}

func (app *WeComApp) getAccessToken(ctx context.Context) (*RespAccessToken, error) {
	params := map[string]string{"corpid": app.corpId, "corpsecret": app.corpSecret}
	resp, err := app.rc.R().
		SetQueryParams(params).
		SetResult(&RespAccessToken{}).
		Get("/cgi-bin/gettoken")
	if err != nil {
		return nil, err
	}

	if result, ok := resp.Result().(*RespAccessToken); ok {
		return result, nil
	}
	return nil, errno.ErrInvalidJson
}

type RespSendMessage struct {
	RespCommon
	Invaliduser    string `json:"invaliduser"`
	Invalidparty   string `json:"invalidparty"`
	Invalidtag     string `json:"invalidtag"`
	Unlicenseduser string `json:"unlicenseduser"`
	Msgid          string `json:"msgid"`
	ResponseCode   string `json:"response_code"`
}

func (app *WeComApp) SendMessage(ctx context.Context, msg message.TxMessage) (*RespSendMessage, error) {
	access_token, err := app.accessTokenManager.Get(ctx)
	if err != nil {
		return nil, err
	}
	resp, err := app.rc.R().
		AddRetryCondition(app.defaultRetryCondition()).
		SetQueryParam("access_token", access_token.Token).
		SetBody(msg).
		SetResult(&RespSendMessage{}).
		Post("/cgi-bin/message/send")
	if err != nil {
		return nil, err
	}

	if result, ok := resp.Result().(*RespSendMessage); ok {
		return result, nil
	}
	return nil, errno.ErrInvalidJson
}

func (app *WeComApp) SendTextMessage(ctx context.Context, msg message.TxTextMessage) (*RespSendMessage, error) {
	return app.SendMessage(ctx, msg)
}

func (app *WeComApp) SendImageMessage(ctx context.Context, msg message.TxImageMessage) (*RespSendMessage, error) {
	return app.SendMessage(ctx, msg)
}

func (app *WeComApp) SendVoiceMessage(ctx context.Context, msg message.TxVoiceMessage) (*RespSendMessage, error) {
	return app.SendMessage(ctx, msg)
}
func (app *WeComApp) SendVideoMessage(ctx context.Context, msg message.TxVideoMessage) (*RespSendMessage, error) {
	return app.SendMessage(ctx, msg)
}

func (app *WeComApp) SendFileMessage(ctx context.Context, msg message.TxFileMessage) (*RespSendMessage, error) {
	return app.SendMessage(ctx, msg)
}

func (app *WeComApp) SendTextCardMessage(ctx context.Context, msg message.TxTextCardMessage) (*RespSendMessage, error) {
	return app.SendMessage(ctx, msg)
}

func (app *WeComApp) SendNewsMessage(ctx context.Context, msg message.TxNewsMessage) (*RespSendMessage, error) {
	return app.SendMessage(ctx, msg)
}

func (app *WeComApp) SendMarkdownMessage(ctx context.Context, msg message.TxMarkdownMessage) (*RespSendMessage, error) {
	return app.SendMessage(ctx, msg)
}

func (app *WeComApp) SendMiniProgramNoticeMessage(ctx context.Context, msg message.TxMiniProgramNoticeMessage) (*RespSendMessage, error) {
	return app.SendMessage(ctx, msg)
}

func (app *WeComApp) SendTextNoticeTemplateCardMessage(ctx context.Context, msg message.TxTextNoticeTemplateCardMessage) (*RespSendMessage, error) {
	return app.SendMessage(ctx, msg)
}

func (app *WeComApp) SendNewsNoticeTemplateCardMessage(ctx context.Context, msg message.TxNewsNoticeTemplateCardMessage) (*RespSendMessage, error) {
	return app.SendMessage(ctx, msg)
}

func (app *WeComApp) SendButtonInteractionTemplateCardMessage(ctx context.Context, msg message.TxButtonInteractionTemplateCardMessage) (*RespSendMessage, error) {
	return app.SendMessage(ctx, msg)
}

func (app *WeComApp) SendVoteInteractionTemplateCardMessage(ctx context.Context, msg message.TxVoteInteractionTemplateCardMessage) (*RespSendMessage, error) {
	return app.SendMessage(ctx, msg)
}

func (app *WeComApp) SendMultipleInteractionTemplateCardMessage(ctx context.Context, msg message.TxMultipleInteractionTemplateCardMessage) (*RespSendMessage, error) {
	return app.SendMessage(ctx, msg)
}

type RespUoloadMedia struct {
	RespCommon
	Type      string `json:"type"`
	MediaId   string `json:"media_id"`
	CreatedAt string `json:"created_at"`
}

func (app *WeComApp) uploadMedia(ctx context.Context, data []byte, mediaType, filename string) (*RespUoloadMedia, error) {
	access_token, err := app.accessTokenManager.Get(ctx)
	if err != nil {
		return nil, err
	}
	resp, err := app.rc.R().
		AddRetryCondition(app.defaultRetryCondition()).
		SetQueryParams(map[string]string{
			"access_token": access_token.Token,
			"type":         mediaType,
		}).
		SetFileReader("media", filename, bytes.NewReader(data)).
		SetContentLength(true).
		SetResult(&RespUoloadMedia{}).
		Post("/cgi-bin/media/upload")
	if err != nil {
		return nil, err
	}

	if result, ok := resp.Result().(*RespUoloadMedia); ok {
		return result, nil
	}
	return nil, errno.ErrInvalidJson
}

func (app *WeComApp) UploadImageMedia(ctx context.Context, url string) (*RespUoloadMedia, error) {
	data, err := downloadMedia(url)
	if err != nil {
		return nil, err
	}
	return app.uploadMedia(ctx, data, "image", path.Base(url))
}

func (app *WeComApp) UploadVoiceMedia(ctx context.Context, url string) (*RespUoloadMedia, error) {
	data, err := downloadMedia(url)
	if err != nil {
		return nil, err
	}
	return app.uploadMedia(ctx, data, "voice", path.Base(url))
}

func (app *WeComApp) UploadVideoMedia(ctx context.Context, url string) (*RespUoloadMedia, error) {
	data, err := downloadMedia(url)
	if err != nil {
		return nil, err
	}
	return app.uploadMedia(ctx, data, "video", path.Base(url))
}

func (app *WeComApp) UploadFileMedia(ctx context.Context, url string) (*RespUoloadMedia, error) {
	data, err := downloadMedia(url)
	if err != nil {
		return nil, err
	}
	return app.uploadMedia(ctx, data, "file", path.Base(url))
}

func downloadMedia(url string) ([]byte, error) {
	resp, err := resty.New().R().Get(url)
	if err != nil {
		return nil, err
	}
	if resp.IsError() {
		return nil, fmt.Errorf("Download Media %s failed, Status Code: %d", url, resp.StatusCode())
	}
	return resp.Body(), nil
}
