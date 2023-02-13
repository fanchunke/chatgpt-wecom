package api

import (
	"context"
	"fmt"
	"strings"

	"github.com/fanchunke/chatgpt-wecom/pkg/wecom"
	"github.com/fanchunke/chatgpt-wecom/pkg/wecom/message"
	"github.com/fanchunke/xgpt3"

	"github.com/rs/zerolog/log"
	gogpt "github.com/sashabaranov/go-gpt3"
)

type callbackHandler struct {
	xgpt3Client *xgpt3.Client
	wecomClient *wecom.WeComApp
}

func NewCallbackHandler(xgpt3Client *xgpt3.Client, wecomClient *wecom.WeComApp) *callbackHandler {
	return &callbackHandler{
		xgpt3Client: xgpt3Client,
		wecomClient: wecomClient,
	}
}

func (h *callbackHandler) OnIncomingMessage(ctx context.Context, msg *message.RxMessage) error {
	agentId, content, err := h.convertMessage(ctx, msg)
	if err != nil {
		log.Ctx(ctx).Error().Err(err).Msgf("Convert wecom msg error: %v", err)
		return err
	}

	go func() {
		defer func() {
			if err := recover(); err != nil {
				log.Error().Msgf("recovery from: %v", err)
			}
		}()

		if err := h.getAndSendGPTResponse(context.Background(), agentId, msg.FromUserName, content); err != nil {
			log.Error().Err(err).Msgf("Get GPT Response error: %v", err)
		}
	}()

	return nil
}

func (h *callbackHandler) convertMessage(ctx context.Context, msg *message.RxMessage) (int64, string, error) {
	// 文本消息
	if msg.Text != nil {
		return msg.Text.AgentId, msg.Text.Content, nil
	}

	return 0, "", fmt.Errorf("UnSupported MsgType: %v", msg.MsgType)
}

func (h *callbackHandler) getAndSendGPTResponse(ctx context.Context, agentId int64, userId, content string) error {
	// 获取 GPT 回复
	req := gogpt.CompletionRequest{
		Model:           gogpt.GPT3TextDavinci003,
		MaxTokens:       1500,
		Prompt:          content,
		TopP:            1,
		Temperature:     0.9,
		PresencePenalty: 0.6,
		User:            userId,
	}
	resp, err := h.xgpt3Client.CreateConversationCompletionWithChannel(ctx, req, fmt.Sprintf("%d", agentId))
	if err != nil {
		return fmt.Errorf("CreateCompletion failed: %w", err)
	}

	if len(resp.Choices) == 0 {
		return fmt.Errorf("Empty GPT Choices")
	}

	// 发送回复给用户
	reply := strings.TrimSpace(resp.Choices[0].Text)
	log.Info().Msgf("Start Send GPT Response: %s", string(reply))
	_, err = h.wecomClient.SendTextMessage(ctx, message.TxTextMessage{
		Text: message.Text{
			Content: reply,
		},
		TxMessageMetadata: message.TxMessageMetadata{
			ToUser:  userId,
			AgentId: agentId,
			MsgType: message.TxTextType,
		},
	})

	if err != nil {
		return fmt.Errorf("Send Lark Message failed: %w", err)
	}

	return nil
}
