package wechat

import (
	"crypto/sha1"
	"encoding/hex"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"sort"
	"strings"

	"github.com/gomydodo/wxencrypter"
)

type Wechat struct {
	Token          string
	AppID          string
	EncodingAesKey string
	Secret         string

	securityMode bool
	encrypter    *wxencrypter.Encrypter
	router       *Router

	defaultHandler Handler

	lastError error
}

type ErrorHandler func(c Context, err error) error

func New(token, appID, encodingAesKey, secret string) (w *Wechat, err error) {

	w = &Wechat{
		Token:  token,
		AppID:  appID,
		Secret: secret,
		router: NewRouter(),
	}

	if encodingAesKey != "" {
		err = w.SetEncodingAesKey(encodingAesKey)
		if err != nil {
			return
		}
	}

	return
}

func (w *Wechat) SetEncodingAesKey(encodingAesKey string) (err error) {
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
		return nil
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
		c, err := newContext(rw, r, w)

		if err != nil {
			rw.WriteHeader(501)
			log.Println("wechat newContext error: ", err)
			return
		}

		w.router.Find(c)

		if h := c.Handler(); h != nil {
			err = h(c)
		} else {
			err = w.DefaultHandler()(c)
		}

		if err != nil {
			log.Println("wechat handle error: ", err)
			rw.WriteHeader(502)
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

func (w *Wechat) MenuClickEventh(h Handler) {
	w.add(MenuClickEventType, "", h)
}
