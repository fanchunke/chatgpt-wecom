package message

const (
	TxTextType              = "text"
	TxImageType             = "image"
	TxVoiceType             = "voice"
	TxVideoType             = "video"
	TxFileType              = "file"
	TxTextCardType          = "textcard"
	TxNewsType              = "news"
	TxMpNewsType            = "mpnews"
	TxMarkdownType          = "markdown"
	TxMiniProgramNoticeType = "miniprogram_notice"
	TxTemplateCardType      = "template_card"
)

type txMessage interface {
	TxTextMessage | TxImageMessage | TxVoiceMessage | TxVideoMessage | TxFileMessage | TxTextCardMessage | TxNewsMessage | TxMarkdownMessage | TxMiniProgramNoticeMessage | TxTextNoticeTemplateCardMessage | TxNewsNoticeTemplateCardMessage | TxButtonInteractionTemplateCardMessage | TxVoteInteractionTemplateCardMessage | TxMultipleInteractionTemplateCardMessage
}

// type TxMessage[T txMessage] struct {
// 	MsgType string `json:"msgtype"`
// 	AgentId int64  `json:"agentid"`
// 	Data    T
// }

// func (tx *TxMessage[T]) Marshal() map[string]interface{} {
// 	data, _ := json.Marshal(tx.Data)
// 	var payload map[string]interface{}
// 	_ = json.Unmarshal(data, &payload)
// 	payload["msgtype"] = tx.MsgType
// 	payload["agentid"] = tx.AgentId
// 	return payload
// }

type TxMessage interface {
	messageType() string
}

type TxMessageMetadata struct {
	ToUser  string `json:"touser"`
	ToParty string `json:"toparty"`
	ToTag   string `json:"totag"`
	MsgType string `json:"msgtype"`
	AgentId int64  `json:"agentid"`
}

func (m TxMessageMetadata) messageType() string {
	return m.MsgType
}

type TxTextMessage struct {
	TxMessageMetadata
	Text          Text  `json:"text"`
	EnableIDTrans int64 `json:"enable_id_trans"`
}

type Text struct {
	Content string `json:"content"`
}

type TxImageMessage struct {
	TxMessageMetadata
	Image                  Image `json:"image"`
	Safe                   int64 `json:"safe"`
	EnableDuplicateCheck   int64 `json:"enable_duplicate_check"`
	DuplicateCheckInterval int64 `json:"duplicate_check_interval"`
}

type Image struct {
	MediaID string `json:"media_id"`
}

type TxVoiceMessage struct {
	TxMessageMetadata
	Voice                  Voice `json:"voice"`
	EnableDuplicateCheck   int64 `json:"enable_duplicate_check"`
	DuplicateCheckInterval int64 `json:"duplicate_check_interval"`
}

type Voice struct {
	MediaID string `json:"media_id"`
}

type TxVideoMessage struct {
	TxMessageMetadata
	Video                  Video `json:"video"`
	Safe                   int64 `json:"safe"`
	EnableDuplicateCheck   int64 `json:"enable_duplicate_check"`
	DuplicateCheckInterval int64 `json:"duplicate_check_interval"`
}

type Video struct {
	MediaID     string `json:"media_id"`
	Title       string `json:"title"`
	Description string `json:"description"`
}

type TxFileMessage struct {
	TxMessageMetadata
	File                   File  `json:"file"`
	Safe                   int64 `json:"safe"`
	EnableDuplicateCheck   int64 `json:"enable_duplicate_check"`
	DuplicateCheckInterval int64 `json:"duplicate_check_interval"`
}

type File struct {
	MediaID string `json:"media_id"`
}

type TxTextCardMessage struct {
	TxMessageMetadata
	Textcard               Textcard `json:"textcard"`
	EnableIDTrans          int64    `json:"enable_id_trans"`
	EnableDuplicateCheck   int64    `json:"enable_duplicate_check"`
	DuplicateCheckInterval int64    `json:"duplicate_check_interval"`
}

type Textcard struct {
	Title       string `json:"title"`
	Description string `json:"description"`
	URL         string `json:"url"`
	Btntxt      string `json:"btntxt"`
}

type TxNewsMessage struct {
	TxMessageMetadata
	News                   News  `json:"news"`
	EnableIDTrans          int64 `json:"enable_id_trans"`
	EnableDuplicateCheck   int64 `json:"enable_duplicate_check"`
	DuplicateCheckInterval int64 `json:"duplicate_check_interval"`
}

type News struct {
	Articles []Article `json:"articles"`
}

type Article struct {
	Title       string `json:"title"`
	Description string `json:"description"`
	URL         string `json:"url"`
	Picurl      string `json:"picurl"`
	Appid       string `json:"appid"`
	Pagepath    string `json:"pagepath"`
}

type TxMarkdownMessage struct {
	TxMessageMetadata
	Markdown               Markdown `json:"markdown"`
	EnableDuplicateCheck   int64    `json:"enable_duplicate_check"`
	DuplicateCheckInterval int64    `json:"duplicate_check_interval"`
}

type Markdown struct {
	Content string `json:"content"`
}

type TxMiniProgramNoticeMessage struct {
	TxMessageMetadata
	MiniprogramNotice      MiniprogramNotice `json:"miniprogram_notice"`
	EnableIDTrans          int64             `json:"enable_id_trans"`
	EnableDuplicateCheck   int64             `json:"enable_duplicate_check"`
	DuplicateCheckInterval int64             `json:"duplicate_check_interval"`
}

type MiniprogramNotice struct {
	Appid             string        `json:"appid"`
	Page              string        `json:"page"`
	Title             string        `json:"title"`
	Description       string        `json:"description"`
	EmphasisFirstItem bool          `json:"emphasis_first_item"`
	ContentItem       []ContentItem `json:"content_item"`
}

type ContentItem struct {
	Key   string `json:"key"`
	Value string `json:"value"`
}

type TxTextNoticeTemplateCardMessage struct {
	TxMessageMetadata
	TemplateCard           TextNoticeTemplateCard `json:"template_card"`
	EnableIDTrans          int64                  `json:"enable_id_trans"`
	EnableDuplicateCheck   int64                  `json:"enable_duplicate_check"`
	DuplicateCheckInterval int64                  `json:"duplicate_check_interval"`
}

type TextNoticeTemplateCard struct {
	CardType              string                            `json:"card_type"`
	Source                TextNoticeSource                  `json:"source"`
	ActionMenu            TextNoticeActionMenu              `json:"action_menu"`
	TaskID                string                            `json:"task_id"`
	MainTitle             TextNoticeEmphasisContent         `json:"main_title"`
	QuoteArea             TextNoticeQuoteArea               `json:"quote_area"`
	EmphasisContent       TextNoticeEmphasisContent         `json:"emphasis_content"`
	SubTitleText          string                            `json:"sub_title_text"`
	HorizontalContentList []TextNoticeHorizontalContentList `json:"horizontal_content_list"`
	JumpList              []TextNoticeCardAction            `json:"jump_list"`
	CardAction            TextNoticeCardAction              `json:"card_action"`
}

type TextNoticeActionMenu struct {
	Desc       string                 `json:"desc"`
	ActionList []TextNoticeActionList `json:"action_list"`
}

type TextNoticeActionList struct {
	Text string `json:"text"`
	Key  string `json:"key"`
}

type TextNoticeCardAction struct {
	Type     int64   `json:"type"`
	URL      *string `json:"url,omitempty"`
	Appid    *string `json:"appid,omitempty"`
	Pagepath *string `json:"pagepath,omitempty"`
	Title    *string `json:"title,omitempty"`
}

type TextNoticeEmphasisContent struct {
	Title string `json:"title"`
	Desc  string `json:"desc"`
}

type TextNoticeHorizontalContentList struct {
	Keyname string  `json:"keyname"`
	Value   string  `json:"value"`
	Type    *int64  `json:"type,omitempty"`
	URL     *string `json:"url,omitempty"`
	MediaID *string `json:"media_id,omitempty"`
	Userid  *string `json:"userid,omitempty"`
}

type TextNoticeQuoteArea struct {
	Type      int64  `json:"type"`
	URL       string `json:"url"`
	Title     string `json:"title"`
	QuoteText string `json:"quote_text"`
}

type TextNoticeSource struct {
	IconURL   string `json:"icon_url"`
	Desc      string `json:"desc"`
	DescColor int64  `json:"desc_color"`
}

type TxNewsNoticeTemplateCardMessage struct {
	TxMessageMetadata
	TemplateCard           NewsNoticeTemplateCard `json:"template_card"`
	EnableIDTrans          int64                  `json:"enable_id_trans"`
	EnableDuplicateCheck   int64                  `json:"enable_duplicate_check"`
	DuplicateCheckInterval int64                  `json:"duplicate_check_interval"`
}

type NewsNoticeTemplateCard struct {
	CardType              string                            `json:"card_type"`
	Source                NewsNoticeSource                  `json:"source"`
	ActionMenu            NewsNoticeActionMenu              `json:"action_menu"`
	TaskID                string                            `json:"task_id"`
	MainTitle             NewsNoticeMainTitle               `json:"main_title"`
	QuoteArea             NewsNoticeQuoteArea               `json:"quote_area"`
	ImageTextArea         NewsNoticeImageTextArea           `json:"image_text_area"`
	CardImage             NewsNoticeCardImage               `json:"card_image"`
	VerticalContentList   []NewsNoticeMainTitle             `json:"vertical_content_list"`
	HorizontalContentList []NewsNoticeHorizontalContentList `json:"horizontal_content_list"`
	JumpList              []NewsNoticeCardAction            `json:"jump_list"`
	CardAction            NewsNoticeCardAction              `json:"card_action"`
}

type NewsNoticeActionMenu struct {
	Desc       string                 `json:"desc"`
	ActionList []NewsNoticeActionList `json:"action_list"`
}

type NewsNoticeActionList struct {
	Text string `json:"text"`
	Key  string `json:"key"`
}

type NewsNoticeCardAction struct {
	Type     int64   `json:"type"`
	URL      *string `json:"url,omitempty"`
	Appid    *string `json:"appid,omitempty"`
	Pagepath *string `json:"pagepath,omitempty"`
	Title    *string `json:"title,omitempty"`
}

type NewsNoticeCardImage struct {
	URL         string  `json:"url"`
	AspectRatio float64 `json:"aspect_ratio"`
}

type NewsNoticeHorizontalContentList struct {
	Keyname string  `json:"keyname"`
	Value   string  `json:"value"`
	Type    *int64  `json:"type,omitempty"`
	URL     *string `json:"url,omitempty"`
	MediaID *string `json:"media_id,omitempty"`
	Userid  *string `json:"userid,omitempty"`
}

type NewsNoticeImageTextArea struct {
	Type     int64  `json:"type"`
	URL      string `json:"url"`
	Title    string `json:"title"`
	Desc     string `json:"desc"`
	ImageURL string `json:"image_url"`
}

type NewsNoticeMainTitle struct {
	Title string `json:"title"`
	Desc  string `json:"desc"`
}

type NewsNoticeQuoteArea struct {
	Type      int64  `json:"type"`
	URL       string `json:"url"`
	Title     string `json:"title"`
	QuoteText string `json:"quote_text"`
}

type NewsNoticeSource struct {
	IconURL   string `json:"icon_url"`
	Desc      string `json:"desc"`
	DescColor int64  `json:"desc_color"`
}

type TxButtonInteractionTemplateCardMessage struct {
	TxMessageMetadata
	TemplateCard           ButtonInteractionTemplateCard `json:"template_card"`
	EnableIDTrans          int64                         `json:"enable_id_trans"`
	EnableDuplicateCheck   int64                         `json:"enable_duplicate_check"`
	DuplicateCheckInterval int64                         `json:"duplicate_check_interval"`
}

type ButtonInteractionTemplateCard struct {
	CardType              string                                   `json:"card_type"`
	Source                ButtonInteractionSource                  `json:"source"`
	ActionMenu            ButtonInteractionActionMenu              `json:"action_menu"`
	MainTitle             ButtonInteractionMainTitle               `json:"main_title"`
	QuoteArea             ButtonInteractionQuoteArea               `json:"quote_area"`
	SubTitleText          string                                   `json:"sub_title_text"`
	HorizontalContentList []ButtonInteractionHorizontalContentList `json:"horizontal_content_list"`
	CardAction            ButtonInteractionCardAction              `json:"card_action"`
	TaskID                string                                   `json:"task_id"`
	ButtonSelection       ButtonInteractionButtonSelection         `json:"button_selection"`
	ButtonList            []ButtonInteractionButtonList            `json:"button_list"`
}

type ButtonInteractionActionMenu struct {
	Desc       string                        `json:"desc"`
	ActionList []ButtonInteractionActionList `json:"action_list"`
}

type ButtonInteractionActionList struct {
	Text string `json:"text"`
	Key  string `json:"key"`
}

type ButtonInteractionButtonList struct {
	Text  string `json:"text"`
	Style int64  `json:"style"`
	Key   string `json:"key"`
}

type ButtonInteractionButtonSelection struct {
	QuestionKey string                        `json:"question_key"`
	Title       string                        `json:"title"`
	OptionList  []ButtonInteractionOptionList `json:"option_list"`
	SelectedID  string                        `json:"selected_id"`
}

type ButtonInteractionOptionList struct {
	ID   string `json:"id"`
	Text string `json:"text"`
}

type ButtonInteractionCardAction struct {
	Type     int64  `json:"type"`
	URL      string `json:"url"`
	Appid    string `json:"appid"`
	Pagepath string `json:"pagepath"`
}

type ButtonInteractionHorizontalContentList struct {
	Keyname string  `json:"keyname"`
	Value   string  `json:"value"`
	Type    *int64  `json:"type,omitempty"`
	URL     *string `json:"url,omitempty"`
	MediaID *string `json:"media_id,omitempty"`
	Userid  *string `json:"userid,omitempty"`
}

type ButtonInteractionMainTitle struct {
	Title string `json:"title"`
	Desc  string `json:"desc"`
}

type ButtonInteractionQuoteArea struct {
	Type      int64  `json:"type"`
	URL       string `json:"url"`
	Title     string `json:"title"`
	QuoteText string `json:"quote_text"`
}

type ButtonInteractionSource struct {
	IconURL   string `json:"icon_url"`
	Desc      string `json:"desc"`
	DescColor int64  `json:"desc_color"`
}

type TxVoteInteractionTemplateCardMessage struct {
	TxMessageMetadata
	TemplateCard           VoteInteractionTemplateCard `json:"template_card"`
	EnableIDTrans          int64                       `json:"enable_id_trans"`
	EnableDuplicateCheck   int64                       `json:"enable_duplicate_check"`
	DuplicateCheckInterval int64                       `json:"duplicate_check_interval"`
}

type VoteInteractionTemplateCard struct {
	CardType     string                      `json:"card_type"`
	Source       VoteInteractionSource       `json:"source"`
	MainTitle    VoteInteractionMainTitle    `json:"main_title"`
	TaskID       string                      `json:"task_id"`
	Checkbox     VoteInteractionCheckbox     `json:"checkbox"`
	SubmitButton VoteInteractionSubmitButton `json:"submit_button"`
}

type VoteInteractionCheckbox struct {
	QuestionKey string                      `json:"question_key"`
	OptionList  []VoteInteractionOptionList `json:"option_list"`
	Mode        int64                       `json:"mode"`
}

type VoteInteractionOptionList struct {
	ID        string `json:"id"`
	Text      string `json:"text"`
	IsChecked bool   `json:"is_checked"`
}

type VoteInteractionMainTitle struct {
	Title string `json:"title"`
	Desc  string `json:"desc"`
}

type VoteInteractionSource struct {
	IconURL string `json:"icon_url"`
	Desc    string `json:"desc"`
}

type VoteInteractionSubmitButton struct {
	Text string `json:"text"`
	Key  string `json:"key"`
}

type TxMultipleInteractionTemplateCardMessage struct {
	TxMessageMetadata
	TemplateCard           MultipleInteractionTemplateCard `json:"template_card"`
	EnableIDTrans          int64                           `json:"enable_id_trans"`
	EnableDuplicateCheck   int64                           `json:"enable_duplicate_check"`
	DuplicateCheckInterval int64                           `json:"duplicate_check_interval"`
}

type MultipleInteractionTemplateCard struct {
	CardType     string                          `json:"card_type"`
	Source       MultipleInteractionSource       `json:"source"`
	MainTitle    MultipleInteractionMainTitle    `json:"main_title"`
	TaskID       string                          `json:"task_id"`
	SelectList   []MultipleInteractionSelectList `json:"select_list"`
	SubmitButton MultipleInteractionSubmitButton `json:"submit_button"`
}

type MultipleInteractionMainTitle struct {
	Title string `json:"title"`
	Desc  string `json:"desc"`
}

type MultipleInteractionSelectList struct {
	QuestionKey string                          `json:"question_key"`
	Title       string                          `json:"title"`
	SelectedID  string                          `json:"selected_id"`
	OptionList  []MultipleInteractionOptionList `json:"option_list"`
}

type MultipleInteractionOptionList struct {
	ID   string `json:"id"`
	Text string `json:"text"`
}

type MultipleInteractionSource struct {
	IconURL string `json:"icon_url"`
	Desc    string `json:"desc"`
}

type MultipleInteractionSubmitButton struct {
	Text string `json:"text"`
	Key  string `json:"key"`
}
