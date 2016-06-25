package wechat

import (
	"encoding/xml"
	"strings"
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
	UnknownType MsgType = iota
	TextType
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

type requestMessage struct {
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

type defaultRequestMessage struct {
	rm requestMessage
}

func (dft *defaultRequestMessage) Unmarshal(d []byte) (err error) {
	err = xml.Unmarshal(d, &dft.rm)
	return
}

func (dft *defaultRequestMessage) ToUserName() string {
	return dft.rm.ToUserName
}

func (dft *defaultRequestMessage) FromUserName() string {
	return dft.rm.FromUserName
}

func (dft *defaultRequestMessage) CreateTime() int {
	return dft.rm.CreateTime
}

func (dft *defaultRequestMessage) MsgType() MsgType {
	switch dft.rm.MsgType {
	case textValue:
		return TextType
	case imageValue:
		return ImageType
	case voiceValue:
		return VioceType
	case videoValue:
		return VideoType
	case shortVideoValue:
		return ShortVideoType
	case locationValue:
		return LocationType
	case linkValue:
		return LinkType
	case eventValue:
		if evt := dft.Event(); evt == unsubscribeEventValue {
			return UnsubscribeEventType
		} else if evt == subscribeEventValue {
			if strings.HasPrefix(dft.EventKey(), "qrscene_") {
				return ScanSubscribeEventType
			} else {
				return SubscribeEventType
			}
		} else if evt == scanEventValue {
			return ScanEventType
		} else if evt == locationEventValue {
			return LocationEventType
		} else if evt == clickEventValue {
			return MenuClickEventType
		} else if evt == viewEventValue {
			return MenuViewEventType
		}
	}
	return UnknownType
}

func (dft *defaultRequestMessage) Content() string {
	return dft.rm.Content
}

func (dft *defaultRequestMessage) MsgId() int64 {
	return dft.rm.MsgId
}

func (dft *defaultRequestMessage) PicUrl() string {
	return dft.rm.PicUrl
}

func (dft *defaultRequestMessage) MediaId() string {
	return dft.rm.MediaId
}

func (dft *defaultRequestMessage) Format() string {
	return dft.rm.Format
}

func (dft *defaultRequestMessage) Recognition() string {
	return dft.rm.Recognition
}

func (dft *defaultRequestMessage) ThumbMediaId() string {
	return dft.rm.ThumbMediaId
}

func (dft *defaultRequestMessage) LocationX() float64 {
	return dft.rm.Location_X
}

func (dft *defaultRequestMessage) LocationY() float64 {
	return dft.rm.Location_Y
}

func (dft *defaultRequestMessage) Scale() int {
	return dft.rm.Scale
}

func (dft *defaultRequestMessage) Label() string {
	return dft.rm.Label
}

func (dft *defaultRequestMessage) Title() string {
	return dft.rm.Title
}

func (dft *defaultRequestMessage) Description() string {
	return dft.rm.Description
}

func (dft *defaultRequestMessage) Url() string {
	return dft.rm.Url
}

func (dft *defaultRequestMessage) Event() string {
	return dft.rm.Event
}

func (dft *defaultRequestMessage) EventKey() string {
	return dft.rm.EventKey
}

func (dft *defaultRequestMessage) Ticket() string {
	return dft.rm.Ticket
}

func (dft *defaultRequestMessage) Latitude() float32 {
	return dft.rm.Latitude
}

func (dft *defaultRequestMessage) Longitude() float32 {
	return dft.rm.Longitude
}

func (dft *defaultRequestMessage) Precision() float32 {
	return dft.rm.Precision
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

type ImageMessage interface {
	baseMessage
	MsgId() int64
	PicUrl() string
	MediaId() string
}

type VoiceMessage interface {
	baseMessage
	MediaId() string
	Recognition() string
	Format() string
	MsgId() int64
}

type VideoMessage interface {
	baseMessage
	MediaId() string
	ThumbMediaId() string
	MsgId() int64
}

type ShortVideoMessage interface {
	baseMessage
	MediaId() string
	ThumbMediaId() string
	MsgId() int64
}

type LocationMessage interface {
	baseMessage
	LocationX()
	LocationY()
	Scale() int
	Label() string
	MsgId() int64
}

type LinkMessage interface {
	baseMessage
	Title() string
	Description() string
	Url() string
	MsgId() int64
}

type SubscribeEventMessage interface {
	baseMessage
	Event() string
}

type UnsubscribeEventMessage interface {
	baseMessage
	Event() string
}

type ScanSubscribeEventMessage interface {
	baseMessage
	Event() string
	EventKey() string
	Ticket() string
}

type ScanEventMessage interface {
	baseMessage
	Event() string
	EventKey() string
	Ticket() string
}

type LocationEventMessage interface {
	baseMessage
	Event() string
	Latitude()
	Longitude() float32
	Precision() float32
}

type MenuViewEventMessage interface {
	baseMessage
	Event() string
	EventKey() string
}

type MenuClickEventMessage interface {
	baseMessage
	Event() string
	EventKey() string
}
