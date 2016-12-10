package wechat

import (
	"crypto/sha1"
	"encoding/hex"
	"io"
	"io/ioutil"
	"net/http"
	"sort"
	"strings"

	"github.com/gomydodo/wxencrypter"
)

type (
	WechatType int

	Wechat struct {
		Token          string
		AppID          string
		EncodingAesKey string
		Secret         string
		Type           WechatType

		securityMode bool
		encrypter    *wxencrypter.Encrypter
		router       *Router

		WechatErrorHandler WechatErrorHandler
		defaultHandler     Handler
		middleware         []Middleware
	}

	Middleware func(Handler) Handler

	Handler func(Context) error

	WechatErrorHandler func(error, Context) error
)

const (
	_ WechatType = iota
	SubscriptionsType
	ServiceType
)

type ErrorHandler func(c Context, err error) error

func New(token, appID, encodingAesKey, secret string, wechatType WechatType) (w *Wechat, err error) {

	w = &Wechat{
		Token:  token,
		AppID:  appID,
		Secret: secret,
		router: NewRouter(),
		Type:   wechatType,
	}

	err = w.SetEncodingAesKey(encodingAesKey)
	if err != nil {
		return
	}

	return
}

func (w *Wechat) SetEncodingAesKey(encodingAesKey string) (err error) {
	if encodingAesKey == "" {
		return
	}

	var encrypter *wxencrypter.Encrypter
	encrypter, err = wxencrypter.NewEncrypter(w.Token, encodingAesKey, w.AppID)

	if err != nil {
		return
	}

	w.EncodingAesKey = encodingAesKey
	w.securityMode = true
	w.encrypter = encrypter

	return
}

func (w *Wechat) SetDefaultHandler(h Handler) {
	w.defaultHandler = h
}

func (w *Wechat) DefaultHandler() Handler {
	if w.defaultHandler != nil {
		return w.defaultHandler
	}

	return func(c Context) error {
		return c.Response().String("success")
	}
}

func (w *Wechat) Signature(timestamp, nonce string) string {
	arr := []string{w.Token, timestamp, nonce}
	sort.Strings(arr)

	joinStr := strings.Join(arr, "")

	h := sha1.New()
	io.WriteString(h, joinStr)
	b := h.Sum(nil)
	return hex.EncodeToString(b)
}

func (w *Wechat) checkSignature(timestamp, nonce, signature string) bool {
	_signed := w.Signature(timestamp, nonce)
	return _signed == signature
}

func (w *Wechat) ServeHTTP(rw http.ResponseWriter, r *http.Request) {
	timestamp := r.URL.Query().Get("timestamp")
	nonce := r.URL.Query().Get("nonce")
	signature := r.URL.Query().Get("signature")
	if !w.checkSignature(timestamp, nonce, signature) {
		rw.WriteHeader(400)
		return
	}

	if r.Method == "GET" {

		echostr := r.URL.Query().Get("echostr")
		rw.Write([]byte(echostr))
		return

	} else if r.Method == "POST" {
		c := newContext(rw, r, w)

		err := c.parse()
		if err != nil {
			if h := w.WechatErrorHandler; h != nil {
				h(err, c)
			}
			return
		}

		w.router.Find(c)

		h := c.Handler()
		if h == nil {
			h = w.DefaultHandler()
		}
		for i := len(w.middleware) - 1; i >= 0; i-- {
			h = w.middleware[i](h)
		}

		if err := h(c); err != nil {
			if h := w.WechatErrorHandler; h != nil {
				h(err, c)
			}
			return
		}
	}

	return
}

func (w *Wechat) body(r *http.Request) (data []byte, err error) {
	data, err = ioutil.ReadAll(r.Body)

	if err != nil {
		return
	}

	if w.securityMode {

		timestamp := r.URL.Query().Get("timestamp")
		nonce := r.URL.Query().Get("nonce")

		msgSignature := r.URL.Query().Get("msg_signature")
		data, err = w.Decrypt(msgSignature, timestamp, nonce, data)
		if err != nil {
			return
		}
	}

	return
}

func (w *Wechat) Decrypt(msgSignature, timestamp, nonce string, data []byte) (d []byte, err error) {
	return w.encrypter.Decrypt(msgSignature, timestamp, nonce, data)
}

func (w *Wechat) Encrypt(d []byte) (b []byte, err error) {
	return w.encrypter.Encrypt(d)
}

func (w *Wechat) Use(m ...Middleware) {
	w.middleware = append(w.middleware, m...)
}

func (w *Wechat) add(msgType MsgType, key string, h Handler) {
	w.router.Add(msgType, key, Route{
		MsgType: TextType,
		Key:     key,
		Handler: h,
	})
}

func (w *Wechat) Text(h Handler) {
	w.add(TextType, "", h)
}

func (w *Wechat) Image(h Handler) {
	w.add(ImageType, "", h)
}

func (w *Wechat) Voice(h Handler) {
	w.add(VoiceType, "", h)
}

func (w *Wechat) ShortVideo(h Handler) {
	w.add(ShortVideoType, "", h)
}

func (w *Wechat) Location(h Handler) {
	w.add(LocationType, "", h)
}

func (w *Wechat) Link(h Handler) {
	w.add(LinkType, "", h)
}

func (w *Wechat) SubscribeEvent(h Handler) {
	w.add(SubscribeEventType, "", h)
}

func (w *Wechat) UnsubscribeEvent(h Handler) {
	w.add(UnsubscribeEventType, "", h)
}

func (w *Wechat) ScanSubscribeEvent(h Handler) {
	w.add(ScanSubscribeEventType, "", h)
}

func (w *Wechat) ScanEvent(h Handler) {
	w.add(ScanEventType, "", h)
}

func (w *Wechat) LocationEvent(h Handler) {
	w.add(LocationEventType, "", h)
}

func (w *Wechat) MenuViewEvent(h Handler) {
	w.add(MenuViewEventType, "", h)
}

func (w *Wechat) MenuClickEvent(h Handler) {
	w.add(MenuClickEventType, "", h)
}

func (w *Wechat) ScancodePushEvent(h Handler) {
	w.add(ScancodePushEventType, "", h)
}

func (w *Wechat) ScancodeWaitmsgEvent(h Handler) {
	w.add(ScancodeWaitmsgEventType, "", h)
}

func (w *Wechat) PicSysphotoEvent(h Handler) {
	w.add(PicSysphotoEventType, "", h)
}

func (w *Wechat) PicPhotoOrAlbumEvent(h Handler) {
	w.add(PicPhotoOrAlbumEventType, "", h)
}

func (w *Wechat) PicWeixinEvent(h Handler) {
	w.add(PicWeixinEventType, "", h)
}

func (w *Wechat) LocationSelectEven(h Handler) {
	w.add(LocationSelectEvenType, "", h)
}

func (w *Wechat) TemplateSendJobFinishEvent(h Handler) {
	w.add(TemplateSendJobFinishEventType, "", h)
}
