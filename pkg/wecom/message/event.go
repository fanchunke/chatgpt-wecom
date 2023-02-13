package message

type RxEventMetadata struct {
	ToUserName   string `xml:"ToUserName"`
	FromUserName string `xml:"FromUserName"`
	CreateTime   int64  `xml:"CreateTime"`
	MsgType      string `xml:"MsgType"`
	Event        string `xml:"Event"`
	AgentId      string `xml:"AgentID"`
}

func (m *RxEventMetadata) kind() string {
	return m.Event
}

// RxSubscribeEvent: 成员关注及取消关注事件
// https://developer.work.weixin.qq.com/document/path/90240#%E6%88%90%E5%91%98%E5%85%B3%E6%B3%A8%E5%8F%8A%E5%8F%96%E6%B6%88%E5%85%B3%E6%B3%A8%E4%BA%8B%E4%BB%B6
type RxSubscribeEvent struct {
	RxEventMetadata
	EventKey string `xml:"EventKey"`
}

// RxUnSubscribeEvent: 成员关注及取消关注事件
// https://developer.work.weixin.qq.com/document/path/90240#%E6%88%90%E5%91%98%E5%85%B3%E6%B3%A8%E5%8F%8A%E5%8F%96%E6%B6%88%E5%85%B3%E6%B3%A8%E4%BA%8B%E4%BB%B6
type RxUnSubscribeEvent struct {
	RxEventMetadata
	EventKey string `xml:"EventKey"`
}

// RxEnterAgentEvent: 进入应用
// https://developer.work.weixin.qq.com/document/path/90240#%E8%BF%9B%E5%85%A5%E5%BA%94%E7%94%A8
type RxEnterAgentEvent struct {
	RxEventMetadata
	EventKey string `xml:"EventKey"`
}

// RxLocationEvent: 上报地理位置
// https://developer.work.weixin.qq.com/document/path/90240#%E4%B8%8A%E6%8A%A5%E5%9C%B0%E7%90%86%E4%BD%8D%E7%BD%AE
type RxLocationEvent struct {
	RxEventMetadata
	Latitude  string `xml:"Latitude"`
	Longitude string `xml:"Longitude"`
	Precision string `xml:"Precision"`
	AppType   string `xml:"AppType"`
}

type BatchJob struct {
	JobId   string `xml:"JobId"`
	JobType string `xml:"JobType"`
	ErrCode int    `xml:"ErrCode"`
	ErrMsg  string `xml:"ErrMsg"`
}

// RxBatchJobResultEvent: 异步任务完成事件推送
// https://developer.work.weixin.qq.com/document/path/90240#%E4%B8%8A%E6%8A%A5%E5%9C%B0%E7%90%86%E4%BD%8D%E7%BD%AE
type RxBatchJobResultEvent struct {
	RxEventMetadata
	BtachJob BatchJob `xml:"BatchJob"`
}

// RxChangeContactEvent: 通讯录变更事件
// https://developer.work.weixin.qq.com/document/path/90240#%E9%80%9A%E8%AE%AF%E5%BD%95%E5%8F%98%E6%9B%B4%E4%BA%8B%E4%BB%B6
type RxChangeContactEvent struct {
	RxEventMetadata
	ChangeType string `xml:"ChangeType"`
}

// RxCreatePartyEvent: 新增部门事件
// https://developer.work.weixin.qq.com/document/path/90240#%E6%96%B0%E5%A2%9E%E9%83%A8%E9%97%A8%E4%BA%8B%E4%BB%B6
type RxCreatePartyEvent struct {
	RxChangeContactEvent
	Id       string `xml:"Id"`
	Name     string `xml:"Name"`
	ParentId string `xml:"ParentId"`
	Order    string `xml:"Order"`
}

// RxUpdatePartyEvent: 更新部门事件
// https://developer.work.weixin.qq.com/document/path/90240#%E6%9B%B4%E6%96%B0%E9%83%A8%E9%97%A8%E4%BA%8B%E4%BB%B6
type RxUpdatePartyEvent struct {
	RxChangeContactEvent
	Id       string `xml:"Id"`
	Name     string `xml:"Name"`
	ParentId string `xml:"ParentId"`
}

// RxDeletePartyEvent: 删除部门事件
// https://developer.work.weixin.qq.com/document/path/90240#%E5%88%A0%E9%99%A4%E9%83%A8%E9%97%A8%E4%BA%8B%E4%BB%B6
type RxDeletePartyEvent struct {
	RxChangeContactEvent
	Id string `xml:"Id"`
}

// RxUpdateTagEvent: 标签成员变更事件
// https://developer.work.weixin.qq.com/document/path/90240#%E6%A0%87%E7%AD%BE%E6%88%90%E5%91%98%E5%8F%98%E6%9B%B4%E4%BA%8B%E4%BB%B6
type RxUpdateTagEvent struct {
	RxChangeContactEvent
	TagId         string `xml:"TagId"`
	AddUserItems  string `xml:"AddUserItems"`
	DelUserItems  string `xml:"DelUserItems"`
	AddPartyItems string `xml:"AddPartyItems"`
	DelPartyItems string `xml:"DelPartyItems"`
}

// RxClickEvent: 点击菜单拉取消息的事件推送
// https://developer.work.weixin.qq.com/document/path/90240#%E7%82%B9%E5%87%BB%E8%8F%9C%E5%8D%95%E6%8B%89%E5%8F%96%E6%B6%88%E6%81%AF%E7%9A%84%E4%BA%8B%E4%BB%B6%E6%8E%A8%E9%80%81
type RxClickEvent struct {
	RxEventMetadata
	EventKey string `xml:"EventKey"`
}

// RxViewEvent: 点击菜单跳转链接的事件推送
// https://developer.work.weixin.qq.com/document/path/90240#%E7%82%B9%E5%87%BB%E8%8F%9C%E5%8D%95%E8%B7%B3%E8%BD%AC%E9%93%BE%E6%8E%A5%E7%9A%84%E4%BA%8B%E4%BB%B6%E6%8E%A8%E9%80%81
type RxViewEvent struct {
	RxEventMetadata
	EventKey string `xml:"EventKey"`
}

type ScanCodeInfo struct {
	ScanType   string `xml:"ScanType"`
	ScanResult string `xml:"ScanResult"`
}

// RxScanCodePushEvent: 扫码推事件的事件推送
// https://developer.work.weixin.qq.com/document/path/90240#%E6%89%AB%E7%A0%81%E6%8E%A8%E4%BA%8B%E4%BB%B6%E7%9A%84%E4%BA%8B%E4%BB%B6%E6%8E%A8%E9%80%81
type RxScanCodePushEvent struct {
	RxEventMetadata
	ScanCodeInfo ScanCodeInfo `xml:"ScanCodeInfo"`
	EventKey     string       `xml:"EventKey"`
}

// RxScanCodeWaitMsgEvent: 扫码推事件且弹出“消息接收中”提示框的事件推送
// https://developer.work.weixin.qq.com/document/path/90240#%E6%89%AB%E7%A0%81%E6%8E%A8%E4%BA%8B%E4%BB%B6%E4%B8%94%E5%BC%B9%E5%87%BA%E2%80%9C%E6%B6%88%E6%81%AF%E6%8E%A5%E6%94%B6%E4%B8%AD%E2%80%9D%E6%8F%90%E7%A4%BA%E6%A1%86%E7%9A%84%E4%BA%8B%E4%BB%B6%E6%8E%A8%E9%80%81
type RxScanCodeWaitMsgEvent struct {
	RxEventMetadata
	ScanCodeInfo ScanCodeInfo `xml:"ScanCodeInfo"`
	EventKey     string       `xml:"EventKey"`
}

type SendPicsInfo struct {
	Count   string      `xml:"Count"`
	PicList SendPicList `xml:"PicList"`
}

type SendPicList struct {
	Item SendPicItem `xml:"item"`
}

type SendPicItem struct {
	PicMd5Sum string `xml:"PicMd5Sum"`
}

// RxPicSysPhotoEvent: 弹出系统拍照发图的事件推送
// https://developer.work.weixin.qq.com/document/path/90240#%E5%BC%B9%E5%87%BA%E7%B3%BB%E7%BB%9F%E6%8B%8D%E7%85%A7%E5%8F%91%E5%9B%BE%E7%9A%84%E4%BA%8B%E4%BB%B6%E6%8E%A8%E9%80%81
type RxPicSysPhotoEvent struct {
	RxEventMetadata
	EventKey     string       `xml:"EventKey"`
	SendPicsInfo SendPicsInfo `xml:"SendPicsInfo"`
}

// RxPicPhotoOrAlbum: 弹出拍照或者相册发图的事件推送
// https://developer.work.weixin.qq.com/document/path/90240#%E5%BC%B9%E5%87%BA%E6%8B%8D%E7%85%A7%E6%88%96%E8%80%85%E7%9B%B8%E5%86%8C%E5%8F%91%E5%9B%BE%E7%9A%84%E4%BA%8B%E4%BB%B6%E6%8E%A8%E9%80%81
type RxPicPhotoOrAlbumEvent struct {
	RxEventMetadata
	EventKey     string       `xml:"EventKey"`
	SendPicsInfo SendPicsInfo `xml:"SendPicsInfo"`
}

// RxPicWeixinEvent: 弹出微信相册发图器的事件推送
// https://developer.work.weixin.qq.com/document/path/90240#%E5%BC%B9%E5%87%BA%E5%BE%AE%E4%BF%A1%E7%9B%B8%E5%86%8C%E5%8F%91%E5%9B%BE%E5%99%A8%E7%9A%84%E4%BA%8B%E4%BB%B6%E6%8E%A8%E9%80%81
type RxPicWeixinEvent struct {
	RxEventMetadata
	EventKey     string       `xml:"EventKey"`
	SendPicsInfo SendPicsInfo `xml:"SendPicsInfo"`
}

type SendLocationInfo struct {
	Location_X string `xml:"Location_X"`
	Location_Y string `xml:"Location_Y"`
	Scale      string `xml:"Scale"`
	Label      string `xml:"Label"`
}

// RxLocationSelectionEvent: 弹出地理位置选择器的事件推送
// https://developer.work.weixin.qq.com/document/path/90240#%E5%BC%B9%E5%87%BA%E5%9C%B0%E7%90%86%E4%BD%8D%E7%BD%AE%E9%80%89%E6%8B%A9%E5%99%A8%E7%9A%84%E4%BA%8B%E4%BB%B6%E6%8E%A8%E9%80%81
type RxLocationSelectionEvent struct {
	RxEventMetadata
	EventKey         string           `xml:"EventKey"`
	AppType          string           `xml:"AppType"`
	SendLocationInfo SendLocationInfo `xml:"SendLocationInfo"`
}

type ApprovalInfo struct {
	ApprovalInfo ApprovalInfoClass `xml:"ApprovalInfo"`
}

type ApprovalInfoClass struct {
	ThirdNo        string        `xml:"ThirdNo"`
	OpenSPName     string        `xml:"OpenSpName"`
	OpenTemplateID string        `xml:"OpenTemplateId"`
	OpenSPStatus   string        `xml:"OpenSpStatus"`
	ApplyTime      string        `xml:"ApplyTime"`
	ApplyUserName  string        `xml:"ApplyUserName"`
	ApplyUserID    string        `xml:"ApplyUserId"`
	ApplyUserParty string        `xml:"ApplyUserParty"`
	ApplyUserImage string        `xml:"ApplyUserImage"`
	ApprovalNodes  ApprovalNodes `xml:"ApprovalNodes"`
	NotifyNodes    NotifyNodes   `xml:"NotifyNodes"`
	Approverstep   string        `xml:"approverstep"`
}

type ApprovalNodes struct {
	ApprovalNode ApprovalNode `xml:"ApprovalNode"`
}

type ApprovalNode struct {
	NodeStatus string            `xml:"NodeStatus"`
	NodeAttr   string            `xml:"NodeAttr"`
	NodeType   string            `xml:"NodeType"`
	Items      ApprovalNodeItems `xml:"Items"`
}

type ApprovalNodeItems struct {
	Item ApprovalNodeItem `xml:"Item"`
}

type ApprovalNodeItem struct {
	ItemName   string `xml:"ItemName"`
	ItemUserID string `xml:"ItemUserId"`
	ItemImage  string `xml:"ItemImage"`
	ItemStatus string `xml:"ItemStatus"`
	ItemOpTime string `xml:"ItemOpTime"`
}

type NotifyNodes struct {
	NotifyNode NotifyNode `xml:"NotifyNode"`
}

type NotifyNode struct {
	ItemName   string `xml:"ItemName"`
	ItemUserID string `xml:"ItemUserId"`
	ItemImage  string `xml:"ItemImage"`
}

// RxOpenApprovalChangeEvent: 审批状态通知事件
// https://developer.work.weixin.qq.com/document/path/90240#%E5%AE%A1%E6%89%B9%E7%8A%B6%E6%80%81%E9%80%9A%E7%9F%A5%E4%BA%8B%E4%BB%B6
type RxOpenApprovalChangeEvent struct {
	RxEventMetadata
	ApprovalInfo ApprovalInfo `xml:"ApprovalInfo"`
}

// RxShareAgentChangeEvent: 企业互联共享应用事件回调
// https://developer.work.weixin.qq.com/document/path/90240#%E4%BC%81%E4%B8%9A%E4%BA%92%E8%81%94%E5%85%B1%E4%BA%AB%E5%BA%94%E7%94%A8%E4%BA%8B%E4%BB%B6%E5%9B%9E%E8%B0%83
type RxShareAgentChangeEvent struct {
	RxEventMetadata
}

// RxShareChainChangeEvent: 上下游共享应用事件回调
// https://developer.work.weixin.qq.com/document/path/90240#%E4%B8%8A%E4%B8%8B%E6%B8%B8%E5%85%B1%E4%BA%AB%E5%BA%94%E7%94%A8%E4%BA%8B%E4%BB%B6%E5%9B%9E%E8%B0%83
type RxShareChainChangeEvent struct {
	RxEventMetadata
}

type SelectedItems struct {
	SelectedItem []SelectedItem `xml:"SelectedItem"`
}

type SelectedItem struct {
	QuestionKey string    `xml:"QuestionKey"`
	OptionIds   OptionIds `xml:"OptionIds"`
}

type OptionIds struct {
	OptionId []string `xml:"OptionId"`
}

// RxTemplateCardEvent: 模板卡片事件推送
// https://developer.work.weixin.qq.com/document/path/90240#%E6%A8%A1%E6%9D%BF%E5%8D%A1%E7%89%87%E4%BA%8B%E4%BB%B6%E6%8E%A8%E9%80%81
type RxTemplateCardEvent struct {
	RxEventMetadata
	EventKey      string        `xml:"EventKey"`
	TaskId        string        `xml:"TaskId"`
	CardType      string        `xml:"CardType"`
	ResponseCode  string        `xml:"ResponseCode"`
	SelectedItems SelectedItems `xml:"SelectedItems"`
}

// RxTemplateCardMenuEvent: 通用模板卡片右上角菜单事件推送
// https://developer.work.weixin.qq.com/document/path/90240#%E9%80%9A%E7%94%A8%E6%A8%A1%E6%9D%BF%E5%8D%A1%E7%89%87%E5%8F%B3%E4%B8%8A%E8%A7%92%E8%8F%9C%E5%8D%95%E4%BA%8B%E4%BB%B6%E6%8E%A8%E9%80%81
type RxTemplateCardMenuEvent struct {
	RxEventMetadata
	EventKey     string `xml:"EventKey"`
	TaskId       string `xml:"TaskId"`
	CardType     string `xml:"CardType"`
	ResponseCode string `xml:"ResponseCode"`
}
