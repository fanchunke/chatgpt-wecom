package cryptor

import (
	"bytes"
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"encoding/base64"
	"encoding/binary"
	"io"

	"github.com/fanchunke/chatgpt-wecom/pkg/wecom/errno"
)

const (
	PKCS7PaddingBlockSize = 32
)

type PlainPayload struct {
	random    []byte
	msgLen    uint32
	Msg       []byte
	ReceiveID []byte
}

type Cryptor struct {
	aesKey        []byte
	entropySource io.Reader
}

func NewCryptor(encodingAESKey string) (*Cryptor, error) {
	aesKey, err := base64.StdEncoding.DecodeString(encodingAESKey + "=")
	if err != nil {
		return nil, err
	}

	obj := &Cryptor{aesKey: aesKey, entropySource: rand.Reader}
	return obj, nil
}

func (c *Cryptor) Decrypt(base64Msg []byte) (PlainPayload, error) {
	msg, err := base64.StdEncoding.DecodeString(string(base64Msg))
	if err != nil {
		return PlainPayload{}, err
	}
	if len(msg) < aes.BlockSize {
		return PlainPayload{}, errno.ErrInvalidEncryptMsgSize
	}

	if len(msg)%aes.BlockSize != 0 {
		return PlainPayload{}, errno.ErrInvalidEncryptMsgSize
	}

	block, err := aes.NewCipher(c.aesKey)
	if err != nil {
		return PlainPayload{}, err
	}

	iv := c.aesKey[:aes.BlockSize]
	mode := cipher.NewCBCDecrypter(block, iv)

	padMsg := make([]byte, len(msg))
	mode.CryptBlocks(padMsg, msg)

	buf, err := c.pkcs7Unpadding(padMsg, PKCS7PaddingBlockSize)
	if err != nil {
		return PlainPayload{}, err
	}
	return c.parsePlainPayload(buf)
}

func (c *Cryptor) Encrypt(payload *PlainPayload) (string, error) {
	resultMsgLen := 16 + 4 + len(payload.Msg) + len(payload.ReceiveID)

	// allocate buffer
	buf := make([]byte, 16, resultMsgLen)

	// add random prefix
	_, err := io.ReadFull(c.entropySource, buf) // len(buf) == 16 at this moment
	if err != nil {
		return "", err
	}

	buf = buf[:cap(buf)] // grow to full capacity
	binary.BigEndian.PutUint32(buf[16:], uint32(len(payload.Msg)))
	copy(buf[20:], payload.Msg)
	copy(buf[20+len(payload.Msg):], payload.ReceiveID)
	buf = c.pkcs7Padding(buf, PKCS7PaddingBlockSize)

	block, err := aes.NewCipher(c.aesKey)
	if err != nil {
		return "", err
	}

	// IV初始向量大小为16字节，取AESKey前16字节
	iv := c.aesKey[:aes.BlockSize]
	mode := cipher.NewCBCEncrypter(block, iv)
	mode.CryptBlocks(buf, buf)
	return base64.StdEncoding.EncodeToString(buf), nil
}

func (c *Cryptor) pkcs7Unpadding(paddingText []byte, blockSize int) ([]byte, error) {
	l := len(paddingText)
	if nil == paddingText || l == 0 {
		return nil, errno.ErrInvalidPKCS7UnpaddingTextSize
	}
	if l%blockSize != 0 {
		return nil, errno.ErrInvalidPKCS7UnpaddingTextSize
	}
	n := int(paddingText[l-1])
	return paddingText[:l-n], nil
}

func (c *Cryptor) pkcs7Padding(plainText []byte, blockSize int) []byte {
	padding := blockSize - (len(plainText) % blockSize)
	padtext := bytes.Repeat([]byte{byte(padding)}, padding)

	var buffer bytes.Buffer
	buffer.Write(plainText)
	buffer.Write(padtext)
	return buffer.Bytes()
}

func (c *Cryptor) parsePlainPayload(plainText []byte) (PlainPayload, error) {
	textLen := uint32(len(plainText))
	if textLen < 20 {
		return PlainPayload{}, errno.ErrInvalidPlainPayloadSize
	}
	random := plainText[:16]
	msgLen := binary.BigEndian.Uint32(plainText[16:20])
	if textLen < (20 + msgLen) {
		return PlainPayload{}, errno.ErrInvalidPlainPayloadSize
	}

	msg := plainText[20 : 20+msgLen]
	receiveID := plainText[20+msgLen:]
	return PlainPayload{
		random:    random,
		msgLen:    msgLen,
		Msg:       msg,
		ReceiveID: receiveID,
	}, nil
}
