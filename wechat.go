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

		WechatErrorHandler WechatErrorHandler
		defaultHandler     Handler
		middleware         []Middleware
	}

	Middleware func(Handler) Handler

	Handler func(Context) error

	WechatErrorHandler func(error, Context)
)

func New(token, appID, encodingAesKey, secret string) (w *Wechat, err error) {

	w = &Wechat{
		Token:  token,
		AppID:  appID,
		Secret: secret,
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
		c := newContext(rw, r, w)

		err := c.parse()
		if err != nil {
			if h := w.WechatErrorHandler; h != nil {
				h(err, c)
			}
			return
		}

		h := w.DefaultHandler()
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
