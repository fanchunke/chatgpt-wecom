package envelope

import (
	"encoding/xml"
	mathRand "math/rand"
	"net/url"
	"strconv"
	"time"

	"github.com/fanchunke/chatgpt-wecom/pkg/wecom/cryptor"
	"github.com/fanchunke/chatgpt-wecom/pkg/wecom/errno"
)

const (
	letterBytes = "0123456789abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"
)

type xmlRxEnvelope struct {
	ToUserName string `xml:"ToUserName"`
	AgentID    string `xml:"AgentID"`
	Encrypt    string `xml:"Encrypt"`
}

type cdataNode struct {
	CData string `xml:",cdata"`
}

type xmlTxEnvelope struct {
	XMLName      xml.Name  `xml:"xml"`
	Encrypt      cdataNode `xml:"Encrypt"`
	MsgSignature cdataNode `xml:"MsgSignature"`
	Timestamp    int64     `xml:"Timestamp"`
	Nonce        cdataNode `xml:"Nonce"`
}

type Envelope struct {
	ToUserName string
	AgentID    string
	Msg        []byte
	ReceiveID  []byte
}

type Processor struct {
	token   string
	cryptor *cryptor.Cryptor
}

func NewProcessor(token, encodingAESKey string) (*Processor, error) {
	obj := &Processor{
		token:   token,
		cryptor: nil,
	}

	c, err := cryptor.NewCryptor(encodingAESKey)
	if err != nil {
		return nil, err
	}

	obj.cryptor = c
	return obj, nil
}

func (p *Processor) Cryptor() *cryptor.Cryptor {
	return p.cryptor
}

func (p *Processor) UnPackMsg(
	url *url.URL,
	body []byte,
) (Envelope, error) {
	var x xmlRxEnvelope
	err := xml.Unmarshal(body, &x)
	if err != nil {
		return Envelope{}, err
	}

	if !VerifyHTTPRequestSignature(p.token, url, x.Encrypt) {
		return Envelope{}, errno.ErrInvalidSignature
	}

	msg, err := p.cryptor.Decrypt([]byte(x.Encrypt))
	if err != nil {
		return Envelope{}, err
	}

	return Envelope{
		ToUserName: x.ToUserName,
		AgentID:    x.AgentID,
		Msg:        msg.Msg,
		ReceiveID:  msg.ReceiveID,
	}, nil
}

func (p *Processor) PackMsg(msg []byte) ([]byte, error) {
	payload := &cryptor.PlainPayload{Msg: msg, ReceiveID: nil}
	encryptedMsg, err := p.cryptor.Encrypt(payload)
	if err != nil {
		return nil, err
	}

	ts := time.Now().Unix()
	nonce := randString(16)
	signature := Sign(p.token, strconv.FormatInt(ts, 10), nonce, encryptedMsg)
	envelope := xmlTxEnvelope{
		XMLName: xml.Name{},
		Encrypt: cdataNode{
			CData: encryptedMsg,
		},
		MsgSignature: cdataNode{
			CData: signature,
		},
		Timestamp: ts,
		Nonce: cdataNode{
			CData: nonce,
		},
	}

	if result, err := xml.Marshal(envelope); err != nil {
		return nil, err
	} else {
		return result, nil
	}
}

func randString(n int) string {
	b := make([]byte, n)
	for i := range b {
		b[i] = letterBytes[mathRand.Int63()%int64(len(letterBytes))]
	}
	return string(b)
}
