package wechat

import (
	"net/http"
)

type Context interface {
	Request() Request
	Text(content string) TextResponseMessage
	Image(mediaId string) ImageResponseMessage
	Voice(mediaId string) VoiceResponseMessage
	Video(video Video) VideoResponseMessage
	Music(music Music) MusicResponseMessage
	Article(articles ...ArticleItem) ArticleResponseMessage
}

type MessageHandler func(c Context) error

type Request interface {
	ToUserName() string
	FromUserName() string
	CreateTime() int
	MsgType() string
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
}

type context struct {
	dft *defaultRequestMessage
	r   *http.Request
	w   http.ResponseWriter
	wc  *Wechat
}

func newContext(w http.ResponseWriter, r *http.Request, wc *Wechat) (c context, err error) {
	c = context{
		r:  r,
		w:  w,
		wc: wc,
	}

	data, err := c.wc.body(r)
	if err != nil {
		return
	}

	dft := &defaultRequestMessage{}
	err = dft.Unmarshal(data)
	if err != nil {
		return
	}

	c.dft = dft

	return
}

func (c context) Request() Request {
	return c.dft
}

func (c context) Text(content string) TextResponseMessage {
	return NewTextResponseMessage(
		c.dft.FromUserName(),
		c.dft.ToUserName(),
		content,
	)
}

func (c context) Image(mediaId string) ImageResponseMessage {
	return NewImageResponseMessage(
		c.dft.FromUserName(),
		c.dft.ToUserName(),
		mediaId,
	)
}

func (c context) Voice(mediaId string) VoiceResponseMessage {
	return NewVoiceResponseMessage(
		c.dft.FromUserName(),
		c.dft.ToUserName(),
		mediaId,
	)
}

func (c context) Video(video Video) VideoResponseMessage {
	return NewVideoResponseMessage(
		c.dft.FromUserName(),
		c.dft.ToUserName(),
		video,
	)
}

func (c context) Music(music Music) MusicResponseMessage {
	return NewMusicResponseMessage(
		c.dft.FromUserName(),
		c.dft.ToUserName(),
		music,
	)
}

func (c context) Article(articles ...ArticleItem) ArticleResponseMessage {
	return NewArticleResponseMessage(
		c.dft.FromUserName(),
		c.dft.ToUserName(),
		articles...,
	)
}
