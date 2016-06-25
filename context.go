package wechat

import (
	"encoding/xml"
	"net/http"
)

type Context interface {
	Request() Request
	Response(data interface{}) error
	Text(content string) error
	Image(mediaId string) error
	Voice(mediaId string) error
	Video(video Video) error
	Music(music Music) error
	Article(articles ...ArticleItem) error
}

type MessageHandler func(c Context) error

type Request interface {
	ToUserName() string
	FromUserName() string
	CreateTime() int
	MsgType() MsgType
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
		r:   r,
		w:   w,
		wc:  wc,
		dft: &defaultRequestMessage{},
	}

	data, err := c.wc.body(r)
	if err != nil {
		return
	}

	err = c.dft.Unmarshal(data)
	return
}

func (c context) Request() Request {
	return c.dft
}

func (c context) Response(data interface{}) (err error) {
	b, err := xml.Marshal(data)
	if err != nil {
		return
	}

	if c.wc.securityMode {
		b, err = c.wc.Encrypt(b)
		if err != nil {
			return
		}
	}

	_, err = c.w.Write(b)
	return
}

func (c context) Text(content string) error {
	return c.Response(NewTextResponseMessage(
		c.Request().FromUserName(),
		c.Request().ToUserName(),
		content,
	))
}

func (c context) Image(mediaId string) error {
	return c.Response(NewImageResponseMessage(
		c.Request().FromUserName(),
		c.Request().ToUserName(),
		mediaId,
	))
}

func (c context) Voice(mediaId string) error {
	return c.Response(NewVoiceResponseMessage(
		c.Request().FromUserName(),
		c.Request().ToUserName(),
		mediaId,
	))
}

func (c context) Video(video Video) error {
	return c.Response(NewVideoResponseMessage(
		c.Request().FromUserName(),
		c.Request().ToUserName(),
		video,
	))
}

func (c context) Music(music Music) error {
	return c.Response(NewMusicResponseMessage(
		c.Request().FromUserName(),
		c.Request().ToUserName(),
		music,
	))
}

func (c context) Article(articles ...ArticleItem) (err error) {
	article, err := NewArticleResponseMessage(
		c.Request().FromUserName(),
		c.Request().ToUserName(),
		articles...,
	)
	if err != nil {
		return
	}

	return c.Response(article)
}
