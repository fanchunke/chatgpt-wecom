package message

type RxMessageMetadata struct {
	ToUserName   string      `xml:"ToUserName"`
	FromUserName string      `xml:"FromUserName"`
	CreateTime   uint32      `xml:"CreateTime"`
	MsgType      MessageType `xml:"MsgType"`
	MsgId        uint64      `xml:"MsgId"`
	AgentId      int64       `xml:"AgentID"`
}

func (m *RxMessageMetadata) kind() string {
	return string(m.MsgType)
}

// RxTextMessage: 文本消息
// https://developer.work.weixin.qq.com/document/path/90239#%E6%96%87%E6%9C%AC%E6%B6%88%E6%81%AF
type RxTextMessage struct {
	RxMessageMetadata
	Content string `xml:"Content"`
}

// RxImageMessage: 图片消息
// https://developer.work.weixin.qq.com/document/path/90239#%E5%9B%BE%E7%89%87%E6%B6%88%E6%81%AF
type RxImageMessage struct {
	RxMessageMetadata
	PicUrl  string `xml:"PicUrl"`
	MediaId string `xml:"MediaId"`
}

// RxVoiceMessage: 语音消息
// https://developer.work.weixin.qq.com/document/path/90239#%E8%AF%AD%E9%9F%B3%E6%B6%88%E6%81%AF
type RxVoiceMessage struct {
	RxMessageMetadata
	MediaId string `xml:"MediaId"`
	Format  string `xml:"Format"`
}

// RxVideoMessage: 视频消息
// https://developer.work.weixin.qq.com/document/path/90239#%E8%A7%86%E9%A2%91%E6%B6%88%E6%81%AF
type RxVideoMessage struct {
	RxMessageMetadata
	MediaId      string `xml:"MediaId"`
	ThumbMediaId string `xml:"ThumbMediaId"`
}

// RxLocationMessage: 位置消息
// https://developer.work.weixin.qq.com/document/path/90239#%E4%BD%8D%E7%BD%AE%E6%B6%88%E6%81%AF
type RxLocationMessage struct {
	RxMessageMetadata
	Location_X string `xml:"Location_X"`
	Location_Y string `xml:"Location_Y"`
	Scale      string `xml:"Scale"`
	Label      string `xml:"Label"`
	AppType    string `xml:"AppType"`
}

// RxLinkMessage: 链接消息
// https://developer.work.weixin.qq.com/document/path/90239#%E9%93%BE%E6%8E%A5%E6%B6%88%E6%81%AF
type RxLinkMessage struct {
	RxMessageMetadata
	Title       string `xml:"Title"`
	Description string `xml:"Description"`
	Url         string `xml:"Url"`
	PicUrl      string `xml:"PicUrl"`
}
