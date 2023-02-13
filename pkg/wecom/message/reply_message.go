package message

import "encoding/xml"

// ReplyMessage: 被动回复消息
type ReplyMessage struct {
	XMLName   xml.Name `xml:"xml"`
	Encrypt   string   `xml:"Encrypt"`
	Signature string   `xml:"MsgSignature"`
	Timestamp string   `xml:"TimeStamp"`
	Nonce     string   `xml:"Nonce"`
}

// func NewReplyMessage(encrypt, signature, timestamp, nonce string) *ReplyMessage {
// 	return &ReplyMessage{
// 		Encrypt:   string{Value: encrypt},
// 		Signature: string{Value: signature},
// 		Timestamp: timestamp,
// 		Nonce:     string{Value: nonce}}
// }

type ReplyMessageMetadata struct {
	ToUserName   string `xml:"ToUserName"`
	FromUserName string `xml:"FromUserName"`
	CreateTime   int64  `xml:"CreateTime"`
	MsgType      string `xml:"MsgType"`
}

// ReplyTextMessage: 被动回复文本消息
// https://developer.work.weixin.qq.com/document/path/90241#%E6%96%87%E6%9C%AC%E6%B6%88%E6%81%AF
type ReplyTextMessage struct {
	ReplyMessageMetadata
	Content string `xml:"Content"`
}

type Media struct {
	MediaId string `xml:"MediaId"`
}

// ReplyImageMessage: 被动回复图片消息
// https://developer.work.weixin.qq.com/document/path/90241#%E5%9B%BE%E7%89%87%E6%B6%88%E6%81%AF
type ReplyImageMessage struct {
	ReplyMessageMetadata
	Image Media `xml:"Image"`
}

// ReplyVoiceMessage: 被动回复语音消息
// https://developer.work.weixin.qq.com/document/path/90241#%E8%AF%AD%E9%9F%B3%E6%B6%88%E6%81%AF
type ReplyVoiceMessage struct {
	ReplyMessageMetadata
	Voice Media `xml:"Voice"`
}

type VideoMedia struct {
	Media
	Title       string `xml:"Title"`
	Description string `xml:"Description"`
}

type ReplyVideoMessage struct {
	ReplyMessageMetadata
	Video VideoMedia `xml:"Video"`
}
