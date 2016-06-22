package wxofficial

import (
	"net/http"
)

const (
	textValue       = "text"
	imageValue      = "image"
	voiceValue      = "voice"
	videoValue      = "video"
	shortVideoValue = "shortvideo"
	locationValue   = "location"
	linkValue       = "link"

	eventValue            = "event"
	subscribeEventValue   = "subscribe"
	unsubscribeEventValue = "unsubscribe"
	scanEventValue        = "SCAN"
	locationEventValue    = "LOCATION"
	clickEventValue       = "CLICK"
	viewEventValue        = "VIEW"
)

type MsgType int

const (
	TextType MsgType = iota
	ImageType
	VioceType
	VideoType
	ShortVideoType
	LocationType
	LinkType

	SubscribeEventType
	UnsubscribeEventType
	ScanEventType
	ScanSubscribeEventType
	LocationEventType
	MenuViewEventType
	MenuClickEventType
)

type RequestMessage struct {
	ToUserName   string
	FromUserName string
	CreateTime   int
	MsgType      string
	Content      string
	MsgId        int64
	PicUrl       string
	MediaId      string
	Format       string
	Recognition  string
	ThumbMediaId string
	Location_X   float64
	Location_Y   float64
	Scale        int
	Label        string
	Title        string
	Description  string
	Url          string
	Event        string
	EventKey     string
	Ticket       string
	Latitude     float32
	Longitude    float32
	Precision    float32
}

type baseMessage interface {
	ToUserName() string
	FromUserName() string
	CreateTime() int
	MsgType() MsgType
}

type TextMessage interface {
	baseMessage
	MsgId() int64
	Content() string
}

type TextMessageHandler interface {
	Handle(t TextMessage) interface{}
}

type defaultTextMessageHandler struct{}

func (dt defaultTextMessageHandler) Handle(t TextMessage) interface{} {
	return NewTextResponseMessage(t.FromUserName(), t.ToUserName(), t.Content())
}

type ImageMessage interface {
	baseMessage
	MsgId() int64
	PicUrl() string
	MediaId() string
}

type ImageMessageHandler interface {
	Handle(i ImageMessage) interface{}
}
