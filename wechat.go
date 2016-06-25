package wechat

import (
	"crypto/sha1"
	"encoding/hex"
	"github.com/gomydodo/wxencrypter"
	"io"
	"io/ioutil"
	"net/http"
	"sort"
	"strings"
)

type Wechat struct {
	Token          string
	AppID          string
	EncodingAesKey string
	securityMode   bool
	encrypter      *wxencrypter.Encrypter
	lastError      error
	handle         MessageHandler
}

func New(token, appID, encodingAesKey string, securityMode bool) (w *Wechat, err error) {

	var encrypter *wxencrypter.Encrypter

	if securityMode {

		encrypter, err = wxencrypter.NewEncrypter(token, encodingAesKey, appID)

		if err != nil {
			return
		}

	}

	w = &Wechat{
		Token:          token,
		AppID:          appID,
		EncodingAesKey: encodingAesKey,
		securityMode:   securityMode,
		encrypter:      encrypter,
	}
	return
}

func (w *Wechat) Use(h MessageHandler) {
	w.handle = h
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
			return
		}

		err = w.handle(c)
		if err != nil {
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
