package message

import (
	"encoding/xml"
	"fmt"
)

type MessageType string

const (
	TextType         MessageType = "text"
	ImageType        MessageType = "image"
	VoiceType        MessageType = "voice"
	VideoType        MessageType = "video"
	LocationType     MessageType = "location"
	LinkType         MessageType = "link"
	EventMessageType MessageType = "event"
)

type EventType string

const (
	SubscribeEventType          EventType = "subscribe"
	UnSubscribeEventType        EventType = "unsubscribe"
	EnterAgentEventType         EventType = "enter_agent"
	LocationEventType           EventType = "LOCATION"
	BatchJobResultEventType     EventType = "batch_job_result"
	ChangeContactEventType      EventType = "change_contact"
	ClickEventType              EventType = "click"
	ViewEventType               EventType = "view"
	ScanCodePushEventType       EventType = "scancode_push"
	ScanCodeWaitMsgEventType    EventType = "scancode_waitmsg"
	PicSysPhotoEventType        EventType = "pic_sysphoto"
	PicPhotoOrAlbumEventType    EventType = "pic_photo_or_album"
	PicWeixinEventType          EventType = "pic_weixin"
	LocationSelectionEventType  EventType = "location_select"
	OpenApprovalChangeEventType EventType = "open_approval_change"
	ShareAgentChangeEventType   EventType = "share_agent_change"
	ShareChainChangeEventType   EventType = "share_chain_change"
	TemplateCardEventType       EventType = "template_card_event"
	TemplateCardMenuEventType   EventType = "template_card_menu_event"
)

type ChangeType string

const (
	CreatePartyEventType ChangeType = "create_party"
	UpdatePartyEventType ChangeType = "update_party"
	DeletePartyEventType ChangeType = "delete_party"
	UpdateTagEventType   ChangeType = "update_tag"
)

type messageKind interface {
	kind() string
}

// RxMessage: 接收到的消息
type RxMessage struct {
	FromUserName string      `xml:"FromUserName"`
	CreateTime   int64       `xml:"CreateTime"`
	MsgId        uint64      `xml:"MsgId"`
	MsgType      MessageType `xml:"MsgType"`
	Event        EventType   `xml:"Event"`
	ChangeType   ChangeType  `xml:"ChangeType"`

	// 接收消息
	Text     *RxTextMessage
	Image    *RxImageMessage
	Voice    *RxVoiceMessage
	Video    *RxVideoMessage
	Location *RxLocationMessage
	Link     *RxLinkMessage

	// 接收事件
	SubscribeEvent          *RxSubscribeEvent
	UnSubscribeEvent        *RxUnSubscribeEvent
	EnterAgentEvent         *RxEnterAgentEvent
	LocationEvent           *RxLocationEvent
	BatchJobResultEvent     *RxBatchJobResultEvent
	CreatePartyEvent        *RxCreatePartyEvent
	UpdatePartyEvent        *RxUpdatePartyEvent
	DeletePartyEvent        *RxDeletePartyEvent
	UpdateTagEvent          *RxUpdateTagEvent
	ClickEvent              *RxClickEvent
	ViewEvent               *RxViewEvent
	ScanCodePushEvent       *RxScanCodePushEvent
	ScanCodeWaitMsgEvent    *RxScanCodeWaitMsgEvent
	PicSysPhotoEvent        *RxPicSysPhotoEvent
	PicPhotoOrAlbumEvent    *RxPicPhotoOrAlbumEvent
	PicWeixinEvent          *RxPicWeixinEvent
	LocationSelectionEvent  *RxLocationSelectionEvent
	OpenApprovalChangeEvent *RxOpenApprovalChangeEvent
	ShareAgentChangeEvent   *RxShareAgentChangeEvent
	ShareChainChangeEvent   *RxShareChainChangeEvent
	TemplateCardEvent       *RxTemplateCardEvent
	TemplateCardMenuEvent   *RxTemplateCardMenuEvent
}

type rxMetadata struct {
	MsgType    MessageType `xml:"MsgType"`
	Event      EventType   `xml:"Event"`
	ChangeType ChangeType  `xml:"ChangeType"`
}

func FromEnvelope(body []byte) (*RxMessage, error) {
	// var metadata rxMetadata
	var err error
	// err = xml.Unmarshal(body, &metadata)
	// if err != nil {
	// 	return nil, err
	// }

	message := &RxMessage{}
	err = xml.Unmarshal(body, message)
	if err != nil {
		return nil, err
	}

	switch message.MsgType {
	case TextType:
		err = message.unmarshal(body, &RxTextMessage{})
	case ImageType:
		err = message.unmarshal(body, &RxImageMessage{})
	case VoiceType:
		err = message.unmarshal(body, &RxVoiceMessage{})
	case VideoType:
		err = message.unmarshal(body, &RxVideoMessage{})
	case LocationType:
		err = message.unmarshal(body, &RxLocationMessage{})
	case LinkType:
		err = message.unmarshal(body, &RxLinkMessage{})
	case EventMessageType:
		switch message.Event {
		case SubscribeEventType:
			err = message.unmarshal(body, &RxSubscribeEvent{})
		case UnSubscribeEventType:
			err = message.unmarshal(body, &RxUnSubscribeEvent{})
		case EnterAgentEventType:
			err = message.unmarshal(body, &RxEnterAgentEvent{})
		case LocationEventType:
			err = message.unmarshal(body, &RxLocationEvent{})
		case BatchJobResultEventType:
			err = message.unmarshal(body, &RxBatchJobResultEvent{})
		case ChangeContactEventType:
			switch message.ChangeType {
			case CreatePartyEventType:
				err = message.unmarshal(body, &RxCreatePartyEvent{})
			case UpdatePartyEventType:
				err = message.unmarshal(body, &RxUpdatePartyEvent{})
			case DeletePartyEventType:
				err = message.unmarshal(body, &RxDeletePartyEvent{})
			case UpdateTagEventType:
				err = message.unmarshal(body, &RxUpdateTagEvent{})
			default:
				err = fmt.Errorf("Unknown Wecom ChangeType: %s", message.ChangeType)
			}
		case ClickEventType:
			err = message.unmarshal(body, &RxClickEvent{})
		case ViewEventType:
			err = message.unmarshal(body, &RxViewEvent{})
		case ScanCodePushEventType:
			err = message.unmarshal(body, &RxScanCodePushEvent{})
		case ScanCodeWaitMsgEventType:
			err = message.unmarshal(body, &RxScanCodeWaitMsgEvent{})
		case PicSysPhotoEventType:
			err = message.unmarshal(body, &RxPicSysPhotoEvent{})
		case PicPhotoOrAlbumEventType:
			err = message.unmarshal(body, &RxPicPhotoOrAlbumEvent{})
		case PicWeixinEventType:
			err = message.unmarshal(body, &RxPicWeixinEvent{})
		case LocationSelectionEventType:
			err = message.unmarshal(body, &RxLocationSelectionEvent{})
		case OpenApprovalChangeEventType:
			err = message.unmarshal(body, &RxOpenApprovalChangeEvent{})
		case ShareAgentChangeEventType:
			err = message.unmarshal(body, &RxShareAgentChangeEvent{})
		case ShareChainChangeEventType:
			err = message.unmarshal(body, &RxShareChainChangeEvent{})
		case TemplateCardEventType:
			err = message.unmarshal(body, &RxTemplateCardEvent{})
		case TemplateCardMenuEventType:
			err = message.unmarshal(body, &RxTemplateCardMenuEvent{})
		default:
			err = fmt.Errorf("Unknown Wecom EventType: %s", message.Event)
		}
	default:
		err = fmt.Errorf("Unknown Wecom MsgType: %s", message.MsgType)
	}
	return message, err
}

func (rx *RxMessage) unmarshal(data []byte, v messageKind) error {
	err := xml.Unmarshal(data, v)
	if err != nil {
		return err
	}
	switch t := v.(type) {
	case *RxTextMessage:
		rx.Text = t
	case *RxImageMessage:
		rx.Image = t
	case *RxVoiceMessage:
		rx.Voice = t
	case *RxVideoMessage:
		rx.Video = t
	case *RxLocationMessage:
		rx.Location = t
	case *RxLinkMessage:
		rx.Link = t
	case *RxSubscribeEvent:
		rx.SubscribeEvent = t
	case *RxUnSubscribeEvent:
		rx.UnSubscribeEvent = t
	case *RxEnterAgentEvent:
		rx.EnterAgentEvent = t
	case *RxLocationEvent:
		rx.LocationEvent = t
	case *RxBatchJobResultEvent:
		rx.BatchJobResultEvent = t
	case *RxCreatePartyEvent:
		rx.CreatePartyEvent = t
	case *RxUpdatePartyEvent:
		rx.UpdatePartyEvent = t
	case *RxDeletePartyEvent:
		rx.DeletePartyEvent = t
	case *RxUpdateTagEvent:
		rx.UpdateTagEvent = t
	case *RxClickEvent:
		rx.ClickEvent = t
	case *RxViewEvent:
		rx.ViewEvent = t
	case *RxScanCodePushEvent:
		rx.ScanCodePushEvent = t
	case *RxScanCodeWaitMsgEvent:
		rx.ScanCodeWaitMsgEvent = t
	case *RxPicSysPhotoEvent:
		rx.PicSysPhotoEvent = t
	case *RxPicPhotoOrAlbumEvent:
		rx.PicPhotoOrAlbumEvent = t
	case *RxPicWeixinEvent:
		rx.PicWeixinEvent = t
	case *RxLocationSelectionEvent:
		rx.LocationSelectionEvent = t
	case *RxOpenApprovalChangeEvent:
		rx.OpenApprovalChangeEvent = t
	case *RxShareAgentChangeEvent:
		rx.ShareAgentChangeEvent = t
	case *RxShareChainChangeEvent:
		rx.ShareChainChangeEvent = t
	case *RxTemplateCardEvent:
		rx.TemplateCardEvent = t
	case *RxTemplateCardMenuEvent:
		rx.TemplateCardMenuEvent = t
	default:
	}
	return nil
}
