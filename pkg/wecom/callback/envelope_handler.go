package callback

import (
	"context"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"

	"github.com/fanchunke/chatgpt-wecom/pkg/wecom/cryptor"
	"github.com/fanchunke/chatgpt-wecom/pkg/wecom/envelope"
)

type EnvelopeHandler interface {
	OnIncomingEnvelope(ctx context.Context, ev envelope.Envelope) error
}

type envelopeHandler struct {
	token   string
	cryptor *cryptor.Cryptor
	ep      *envelope.Processor
	eh      EnvelopeHandler
}

func NewEnvelopeHandler(token, encodingAESKey string, eh EnvelopeHandler) (*envelopeHandler, error) {
	ep, err := envelope.NewProcessor(token, encodingAESKey)
	if err != nil {
		return nil, err
	}

	return &envelopeHandler{
		token:   token,
		cryptor: ep.Cryptor(),
		ep:      ep,
		eh:      eh,
	}, nil
}

func (h *envelopeHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		h.echoHandler(w, r)
	case http.MethodPost:
		h.eventHandler(w, r)
	default:
		w.WriteHeader(http.StatusNotImplemented)
	}
}

func (h *envelopeHandler) echoHandler(w http.ResponseWriter, r *http.Request) {
	url := r.URL

	if !envelope.VerifyHTTPRequestSignature(h.token, url, "") {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	echoStr := url.Query().Get("echostr")
	if echoStr == "" {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	payload, err := h.cryptor.Decrypt([]byte(echoStr))
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write(payload.Msg)
}

func (h *envelopeHandler) eventHandler(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	ev, err := h.ep.UnPackMsg(r.URL, body)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	// err = h.eh.OnIncomingEnvelope(r.Context(), ev)
	// if err != nil {
	// 	w.WriteHeader(http.StatusInternalServerError)
	// 	return
	// }

	go func() {
		defer func() {
			if err := recover(); err != nil {
				fmt.Printf("recover ok :%s\n", err)
			}
		}()

		err := h.eh.OnIncomingEnvelope(r.Context(), ev)
		log.Println(err)
	}()

	w.WriteHeader(http.StatusOK)
}
