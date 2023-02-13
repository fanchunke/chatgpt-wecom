package callback

import (
	"context"
	"net/http"

	"github.com/fanchunke/chatgpt-wecom/pkg/wecom/envelope"
	"github.com/fanchunke/chatgpt-wecom/pkg/wecom/message"

	"github.com/rs/zerolog/log"
)

type RxMessageHandler interface {
	OnIncomingMessage(ctx context.Context, msg *message.RxMessage) error
}

type rxMessageHandler struct {
	handler RxMessageHandler
}

func (h *rxMessageHandler) OnIncomingEnvelope(ctx context.Context, ev envelope.Envelope) error {
	log.Ctx(ctx).Debug().Msgf("收到企业微信消息: %s", ev.Msg)
	msg, err := message.FromEnvelope(ev.Msg)
	if err != nil {
		return err
	}
	return h.handler.OnIncomingMessage(ctx, msg)
}

func NewHTTPHandler(token, encodingAESKey string, handler RxMessageHandler) (http.Handler, error) {
	h := &rxMessageHandler{handler: handler}
	eh, err := NewEnvelopeHandler(token, encodingAESKey, h)
	if err != nil {
		return nil, err
	}

	return eh, nil
}
