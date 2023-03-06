package api

import (
	"context"
	"fmt"
	"strings"

	config "github.com/fanchunke/chatgpt-wecom/conf"
	"github.com/fanchunke/chatgpt-wecom/pkg/wecom"
	"github.com/fanchunke/chatgpt-wecom/pkg/wecom/message"
	"github.com/fanchunke/xgpt3"

	"github.com/rs/zerolog/log"
	openai "github.com/sashabaranov/go-openai"
)

type versionType string

const (
	callbackVersionV1 versionType = "v1"
	callbackVersionV2 versionType = "v2"
)

type callbackHandler struct {
	cfg         *config.Config
	xgpt3Client *xgpt3.Client
	wecomClient *wecom.WeComApp
	version     versionType
}

func NewCallbackHandler(cfg *config.Config, xgpt3Client *xgpt3.Client, wecomClient *wecom.WeComApp, version versionType) *callbackHandler {
	return &callbackHandler{
		cfg:         cfg,
		xgpt3Client: xgpt3Client,
		wecomClient: wecomClient,
		version:     version,
	}
}

func (h *callbackHandler) OnIncomingMessage(ctx context.Context, msg *message.RxMessage) error {
	var reply string
	var err error
	if msg.Text != nil {
		// 判断是否需要重启会话
		content := msg.Text.Content
		closeSession := h.cfg.Conversation.CloseSessionFlag == content

		// 获取回复
		if !closeSession {
			var handler func(ctx context.Context, agentId int64, userId string, content string) (string, error)
			if h.version == callbackVersionV1 {
				handler = h.getOpenAICompletion
			} else {
				handler = h.getOpenAIChatCompletion
			}

			reply, err = handler(context.Background(), msg.AgentId, msg.FromUserName, content)
			if err != nil {
				log.Error().Err(err).Msgf("Get GPT Response error: %v", err)
				return err
			}
		} else {
			if err := h.xgpt3Client.CloseConversation(context.Background(), msg.FromUserName); err != nil {
				log.Error().Err(err).Msgf("Close Conversation error: %v", err)
				return err
			}
			reply = h.cfg.Conversation.CloseSessionReply
		}
	} else if msg.EnterAgentEvent != nil {
		if h.cfg.Conversation.EnableEnterEvent {
			reply = h.cfg.Conversation.EnterEventReply
		}
	} else {
		return fmt.Errorf("UnSupported MsgType: %v", msg.MsgType)
	}

	// 发送回复
	if reply == "" {
		log.Debug().Msg("Reply is empty")
		return nil
	}
	if err := h.sendTextMessage(context.Background(), msg.AgentId, msg.FromUserName, reply); err != nil {
		log.Error().Err(err).Msgf("Send Wecom Response error: %v", err)
		return err
	}

	return nil
}

func (h *callbackHandler) getOpenAICompletion(ctx context.Context, agentId int64, userId, content string) (string, error) {
	// 获取 GPT 回复
	req := openai.CompletionRequest{
		Model:           openai.GPT3TextDavinci003,
		MaxTokens:       1500,
		Prompt:          content,
		TopP:            1,
		Temperature:     0.9,
		PresencePenalty: 0.6,
		User:            userId,
	}
	resp, err := h.xgpt3Client.CreateConversationCompletionWithChannel(ctx, req, fmt.Sprintf("%d", agentId))
	if err != nil {
		return "", fmt.Errorf("CreateCompletion failed: %w", err)
	}

	if len(resp.Choices) == 0 {
		return "", fmt.Errorf("Empty GPT Choices")
	}

	// 发送回复给用户
	reply := strings.TrimSpace(resp.Choices[0].Text)
	return reply, nil
}

func (h *callbackHandler) sendTextMessage(ctx context.Context, agentId int64, userId, content string) error {
	log.Info().Msgf("[AgentId: %d] [UserId: %s] Start Send Wecom Response: %s", agentId, userId, string(content))
	_, err := h.wecomClient.SendTextMessage(ctx, message.TxTextMessage{
		Text: message.Text{
			Content: content,
		},
		TxMessageMetadata: message.TxMessageMetadata{
			ToUser:  userId,
			AgentId: agentId,
			MsgType: message.TxTextType,
		},
	})

	if err != nil {
		return fmt.Errorf("Send Wecom Message failed: %w", err)
	}

	return nil
}

func (h *callbackHandler) getOpenAIChatCompletion(ctx context.Context, agentId int64, userId, content string) (string, error) {
	// 获取 GPT 回复
	req := openai.ChatCompletionRequest{
		Model:     openai.GPT3Dot5Turbo,
		MaxTokens: 1500,
		Messages: []openai.ChatCompletionMessage{
			{
				Role:    openai.ChatMessageRoleUser,
				Content: content,
			},
		},
		TopP:            1,
		Temperature:     0.9,
		PresencePenalty: 0.6,
		User:            userId,
	}
	resp, err := h.xgpt3Client.CreateChatCompletionWithChannel(ctx, req, fmt.Sprintf("%d", agentId))
	if err != nil {
		return "", fmt.Errorf("CreateCompletion failed: %w", err)
	}

	if len(resp.Choices) == 0 {
		return "", fmt.Errorf("Empty GPT Choices")
	}

	// 发送回复给用户
	reply := strings.TrimSpace(resp.Choices[0].Message.Content)
	return reply, nil
}
