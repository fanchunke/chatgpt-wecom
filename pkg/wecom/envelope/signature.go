package envelope

import (
	"bytes"
	"crypto/sha1"
	"fmt"
	"net/url"
	"sort"
)

const (
	msgSignatureKey = "msg_signature"
	timestampKey    = "timestamp"
	nonceKey        = "nonce"
)

type Signature interface {
	GetParam(key string) (string, bool)
	GetParams() ([]string, bool)
}

func Sign(params ...string) string {
	arr := make([]string, len(params))
	copy(arr, params)

	sort.Strings(arr)

	var buffer bytes.Buffer
	for _, value := range arr {
		buffer.WriteString(value)
	}

	sha := sha1.New()
	sha.Write(buffer.Bytes())
	return fmt.Sprintf("%x", sha.Sum(nil))
}

func VerifySignature(token string, s Signature) bool {
	msgSignature, ok := s.GetParam(msgSignatureKey)
	if !ok {
		return false
	}

	params, ok := s.GetParams()
	if !ok {
		return false
	}

	signature := Sign(append(params, token)...)
	return msgSignature == signature

}

func VerifyHTTPRequestSignature(token string, url *url.URL, body string) bool {
	warpped := &httpRequestWithSignature{
		url:  url,
		body: body,
	}
	return VerifySignature(token, warpped)
}

type httpRequestWithSignature struct {
	url  *url.URL
	body string
}

func (h *httpRequestWithSignature) GetParam(key string) (string, bool) {
	v := h.url.Query().Get(key)
	if len(v) == 0 {
		return "", false
	}
	return v, true
}

func (h *httpRequestWithSignature) GetParams() ([]string, bool) {
	params := make([]string, 0)
	for key, values := range h.url.Query() {
		if key == msgSignatureKey {
			continue
		}
		params = append(params, values...)
	}

	if len(h.body) > 0 {
		params = append(params, h.body)
	}
	return params, len(params) > 0
}
