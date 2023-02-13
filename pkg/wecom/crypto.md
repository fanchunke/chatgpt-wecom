package wecom

import (
	"bytes"
	"crypto/aes"
	"crypto/cipher"
	"crypto/sha1"
	"encoding/base64"
	"encoding/binary"
	"fmt"
	"math/rand"
	"sort"
)

const (
	letterBytes           = "0123456789abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"
	PKCS7PaddingBlockSize = 32
)

const (
	ValidateSignatureError int = -40001
	ParseXmlError          int = -40002
	ComputeSignatureError  int = -40003
	IllegalAesKey          int = -40004
	ValidateCorpidError    int = -40005
	EncryptAESError        int = -40006
	DecryptAESError        int = -40007
	IllegalBuffer          int = -40008
	EncodeBase64Error      int = -40009
	DecodeBase64Error      int = -40010
	GenXmlError            int = -40010
	ParseJsonError         int = -40012
	GenJsonError           int = -40013
	IllegalProtocolType    int = -40014
)

type CryptError struct {
	ErrCode int
	ErrMsg  string
}

func (e *CryptError) Error() string {
	return fmt.Sprintf("ErrCode: %d, ErrMsg: %s", e.ErrCode, e.ErrMsg)
}

func NewCryptError(errCode int, errMsg string) *CryptError {
	return &CryptError{ErrCode: errCode, ErrMsg: errMsg}
}

type ProtocolType int

const (
	XmlType ProtocolType = 1
)

type ProtocolProcessor interface {
	parse(msg []byte) (*RxMessage, *CryptError)
	serialize(msg *ReplyMessage) ([]byte, *CryptError)
}

type Cryptor struct {
	token             string
	encodingAESKey    string
	receiverId        string
	protocolProcessor ProtocolProcessor
}

func NewCryptor(receiverId, encodingAESKey, token string, protocolType ProtocolType) *Cryptor {
	var p ProtocolProcessor
	if protocolType != XmlType {
		panic("unsupport protocal")
	} else {
		p = new(XmlProcessor)
	}
	return &Cryptor{
		token:             token,
		encodingAESKey:    encodingAESKey + "=",
		receiverId:        receiverId,
		protocolProcessor: p,
	}
}

func (c *Cryptor) pkcs7Unpadding(paddingText []byte, blockSize int) ([]byte, *CryptError) {
	l := len(paddingText)
	if nil == paddingText || l == 0 {
		return nil, NewCryptError(DecryptAESError, "PKCS7Unpadding text nil or zero")
	}
	if l%blockSize != 0 {
		return nil, NewCryptError(DecryptAESError, "PKCS7Unpadding text not a multiple of the block size")
	}
	n := int(paddingText[l-1])
	return paddingText[:l-n], nil
}

func (c *Cryptor) pkcs7Padding(plainText string, blockSize int) []byte {
	padding := blockSize - (len(plainText) % blockSize)
	padtext := bytes.Repeat([]byte{byte(padding)}, padding)

	var buffer bytes.Buffer
	buffer.WriteString(plainText)
	buffer.Write(padtext)
	return buffer.Bytes()
}

func (c *Cryptor) cbcEncrypter(plainText string) ([]byte, *CryptError) {
	// AESKey: EncodingAESKey Base64 解码
	aesKey, err := base64.StdEncoding.DecodeString(c.encodingAESKey)
	if err != nil {
		return nil, NewCryptError(DecodeBase64Error, err.Error())
	}

	block, err := aes.NewCipher(aesKey)
	if err != nil {
		return nil, NewCryptError(EncryptAESError, err.Error())
	}

	// IV初始向量大小为16字节，取AESKey前16字节
	iv := aesKey[:aes.BlockSize]
	mode := cipher.NewCBCEncrypter(block, iv)

	padMsg := c.pkcs7Padding(plainText, PKCS7PaddingBlockSize)
	cipherText := make([]byte, len(padMsg))
	mode.CryptBlocks(cipherText, padMsg)
	encryptedMsg := make([]byte, base64.StdEncoding.EncodedLen(len(cipherText)))
	base64.StdEncoding.Encode(encryptedMsg, cipherText)
	return encryptedMsg, nil
}

func (c *Cryptor) cbcDecrypter(base64EncryptedMsg string) ([]byte, *CryptError) {
	aesKey, err := base64.StdEncoding.DecodeString(c.encodingAESKey)
	if err != nil {
		return nil, NewCryptError(DecodeBase64Error, err.Error())
	}

	encryptedMsg, err := base64.StdEncoding.DecodeString(base64EncryptedMsg)
	if err != nil {
		return nil, NewCryptError(DecodeBase64Error, err.Error())
	}
	if len(encryptedMsg) < aes.BlockSize {
		return nil, NewCryptError(DecryptAESError, "encrypt_msg size is not valid")
	}

	if len(encryptedMsg)%aes.BlockSize != 0 {
		return nil, NewCryptError(DecryptAESError, "encrypt_msg not a multiple of the block size")
	}

	block, err := aes.NewCipher(aesKey)
	if err != nil {
		return nil, NewCryptError(DecryptAESError, err.Error())
	}

	iv := aesKey[:aes.BlockSize]
	mode := cipher.NewCBCDecrypter(block, iv)

	padMsg := make([]byte, len(encryptedMsg))
	mode.CryptBlocks(padMsg, encryptedMsg)

	return c.pkcs7Unpadding(padMsg, PKCS7PaddingBlockSize)
}

type PlainPayload struct {
	random     []byte
	msgLen     uint32
	Msg        []byte
	ReceiverId []byte
}

func (c *Cryptor) parsePlainPayload(plainText []byte) (*PlainPayload, *CryptError) {
	textLen := uint32(len(plainText))
	if textLen < 20 {
		return nil, NewCryptError(IllegalBuffer, "plain is to small 1")
	}
	random := plainText[:16]
	msgLen := binary.BigEndian.Uint32(plainText[16:20])
	if textLen < (20 + msgLen) {
		return nil, NewCryptError(IllegalBuffer, "plain is to small 2")
	}

	msg := plainText[20 : 20+msgLen]
	receiverId := plainText[20+msgLen:]
	return &PlainPayload{
		random:     random,
		msgLen:     msgLen,
		Msg:        msg,
		ReceiverId: receiverId,
	}, nil
}

func (c *Cryptor) sign(timestamp, nonce, data string) string {
	arr := []string{c.token, timestamp, nonce, data}
	sort.Strings(arr)
	var buffer bytes.Buffer
	for _, value := range arr {
		buffer.WriteString(value)
	}

	sha := sha1.New()
	sha.Write(buffer.Bytes())
	return fmt.Sprintf("%x", sha.Sum(nil))
}

func (c *Cryptor) VerifyURL(msgSign, timestamp, nonce, echoStr string) ([]byte, *CryptError) {
	signature := c.sign(timestamp, nonce, echoStr)
	if msgSign != signature {
		return nil, NewCryptError(ValidateSignatureError, "Invalid signature")
	}

	plaintext, err := c.cbcDecrypter(echoStr)
	if nil != err {
		return nil, err
	}

	msg, err := c.parsePlainPayload(plaintext)
	if nil != err {
		return nil, err
	}

	if len(msg.ReceiverId) > 0 && string(msg.ReceiverId) != c.receiverId {
		return nil, NewCryptError(ValidateCorpidError, "receiver_id is not equil")
	}

	return msg.Msg, nil
}

func (c *Cryptor) RandString(n int) string {
	b := make([]byte, n)
	for i := range b {
		b[i] = letterBytes[rand.Int63()%int64(len(letterBytes))]
	}
	return string(b)
}

func (c *Cryptor) Encrypt(timestamp, nonce string, plainText []byte) ([]byte, *CryptError) {
	var buffer bytes.Buffer
	// random
	buffer.WriteString(c.RandString(16))

	// msg_len
	msgLenBuf := make([]byte, 4)
	binary.BigEndian.PutUint32(msgLenBuf, uint32(len(string(plainText))))
	buffer.Write(msgLenBuf)
	// msg
	buffer.Write(plainText)
	// receiver_id
	buffer.WriteString(c.receiverId)

	b, err := c.cbcEncrypter(buffer.String())
	if err != nil {
		return nil, err
	}
	cipherText := string(b)
	signature := c.sign(timestamp, nonce, cipherText)
	msg := NewReplyMessage(cipherText, signature, timestamp, nonce)
	return c.protocolProcessor.serialize(msg)
}

func (c *Cryptor) Decrypt(msgSign, timestamp, nonce string, cipherText []byte) (*PlainPayload, *CryptError) {
	m, err := c.protocolProcessor.parse(cipherText)
	if err != nil {
		return nil, err
	}

	signature := c.sign(timestamp, nonce, m.Encrypt)
	if msgSign != signature {
		return nil, NewCryptError(ValidateSignatureError, "signature not equal")
	}

	plainText, err := c.cbcDecrypter(m.Encrypt)
	if err != nil {
		return nil, err
	}

	payload, err := c.parsePlainPayload(plainText)
	if err != nil {
		return nil, err
	}

	if len(c.receiverId) > 0 && c.receiverId != string(payload.ReceiverId) {
		return nil, NewCryptError(ValidateCorpidError, "receiver_id is not equal")
	}

	return payload, nil
}
