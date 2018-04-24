package wechat

import "net/http"

type Context interface {
	Wechat() *Wechat

	Request() Request

	IsText() bool
	IsImage() bool
	IsVoice() bool
	IsVideo() bool
	IsShortVideo() bool
	IsLocation() bool
	IsLink() bool
	IsEvent() bool

	IsSubscribeEvent() bool
	IsUnsubscribeEvent() bool
	IsScanEvent() bool
	IsScanSubscribeEvent() bool
	IsLocationEvent() bool
	IsMenuViewEvent() bool
	IsMenuClickEvent() bool

	IsScancodePushEvent() bool
	IsScancodeWaitmsgEvent() bool
	IsPicSysphotoEvent() bool
	IsPicPhotoOrAlbumEvent() bool
	IsPicWeixinEvent() bool
	IsLocationSelectEven() bool

	IsTemplateSendJobFinishEvent() bool

	Response() Response

	Success() error
	String(s string) error
	Bytes(b []byte) error
	XML(data interface{}) error
	Text(content string) error
	Image(mediaID string) error
	Voice(mediaID string) error
	Video(video Video) error
	Music(music Music) error
	Article(articles ...ArticleItem) error
}

type Request interface {
	ToUserName() string
	FromUserName() string
	CreateTime() int
	MsgType() MsgType
	EventType() EvtType
	Content() string
	MsgId() int64
	PicUrl() string
	MediaId() string
	Format() string
	Recognition() string
	ThumbMediaId() string
	LocationX() float64
	LocationY() float64
	Scale() int
	Label() string
	Title() string
	Description() string
	Url() string
	Event() string
	EventKey() string
	Ticket() string
	Latitude() float32
	Longitude() float32
	Precision() float32

	MenuId() int64
	ScanCodeInfo() ScanCodeInfo
	SendPicsInfo() SendPicsInfo
	SendLocationInfo() SendLocationInfo

	Status() string
}

type Response interface {
	Success() error
	String(s string) error
	Bytes(b []byte) error
	XML(data interface{}) error
	Text(content string) error
	Image(mediaID string) error
	Voice(mediaID string) error
	Video(video Video) error
	Music(music Music) error
	Article(articles ...ArticleItem) error
}

type context struct {
	drm *defaultRequestMessage
	r   *http.Request
	w   http.ResponseWriter
	wc  *Wechat
	dr  defaultResponse
}

func newContext(w http.ResponseWriter, r *http.Request, wc *Wechat) (c *context) {
	c = &context{
		r:   r,
		w:   w,
		wc:  wc,
		drm: &defaultRequestMessage{},
	}

	c.dr = defaultResponse{c: c, w: w}

	return
}

func (c *context) parse() (err error) {
	data, err := c.wc.body(c.r)
	if err != nil {
		return
	}

	err = c.drm.Unmarshal(data)
	return
}

func (c *context) Wechat() *Wechat {
	return c.wc
}

func (c *context) Request() Request {
	return c.drm
}

func (c *context) IsText() bool {
	return c.drm.MsgType() == TextType
}

func (c *context) IsImage() bool {
	return c.drm.MsgType() == ImageType
}

func (c *context) IsVoice() bool {
	return c.drm.MsgType() == VoiceType
}

func (c *context) IsVideo() bool {
	return c.drm.MsgType() == VideoType
}

func (c *context) IsShortVideo() bool {
	return c.drm.MsgType() == ShortVideoType
}

func (c *context) IsLocation() bool {
	return c.drm.MsgType() == LocationType
}

func (c *context) IsLink() bool {
	return c.drm.MsgType() == LinkType
}

func (c *context) IsEvent() bool {
	return c.drm.MsgType() == EventType
}

func (c *context) IsSubscribeEvent() bool {
	return c.drm.EventType() == SubscribeEventType || c.drm.EventType() == ScanSubscribeEventType
}

func (c *context) IsUnsubscribeEvent() bool {
	return c.drm.EventType() == UnsubscribeEventType
}

func (c *context) IsScanEvent() bool {
	return c.drm.EventType() == ScanEventType
}

func (c *context) IsScanSubscribeEvent() bool {
	return c.drm.EventType() == ScanSubscribeEventType
}

func (c *context) IsLocationEvent() bool {
	return c.drm.EventType() == LocationEventType
}

func (c *context) IsMenuViewEvent() bool {
	return c.drm.EventType() == MenuViewEventType
}

func (c *context) IsMenuClickEvent() bool {
	return c.drm.EventType() == MenuClickEventType
}

func (c *context) IsScancodePushEvent() bool {
	return c.drm.EventType() == ScancodePushEventType
}

func (c *context) IsScancodeWaitmsgEvent() bool {
	return c.drm.EventType() == ScancodeWaitmsgEventType
}

func (c *context) IsPicSysphotoEvent() bool {
	return c.drm.EventType() == PicSysphotoEventType
}

func (c *context) IsPicPhotoOrAlbumEvent() bool {
	return c.drm.EventType() == PicPhotoOrAlbumEventType
}

func (c *context) IsPicWeixinEvent() bool {
	return c.drm.EventType() == PicWeixinEventType
}

func (c *context) IsLocationSelectEven() bool {
	return c.drm.EventType() == LocationSelectEventType
}

func (c *context) IsTemplateSendJobFinishEvent() bool {
	return c.drm.EventType() == TemplateSendJobFinishEventType
}

func (c *context) Response() Response {
	return c.dr
}

func (c *context) Success() error {
	return c.dr.Success()
}

func (c *context) String(s string) error {
	return c.dr.String(s)
}

func (c *context) Bytes(b []byte) error {
	return c.dr.Bytes(b)
}

func (c *context) XML(data interface{}) error {
	return c.dr.XML(data)
}

func (c *context) Text(content string) error {
	return c.dr.Text(content)
}

func (c *context) Image(mediaID string) error {
	return c.dr.Image(mediaID)
}

func (c *context) Voice(mediaID string) error {
	return c.dr.Voice(mediaID)
}

func (c *context) Video(video Video) error {
	return c.dr.Video(video)
}

func (c *context) Music(music Music) error {
	return c.dr.Music(music)
}

func (c *context) Article(articles ...ArticleItem) error {
	return c.dr.Article(articles...)
}
