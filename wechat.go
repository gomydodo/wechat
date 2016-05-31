package wxofficial

import (
	"crypto/sha1"
	"encoding/hex"
	"encoding/xml"
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
	encrypter      wxencrypter.Encrypter
	lastError      error
}

func New(token, appID, encodingAesKey string, securityMode bool) (w *Wechat, err error) {

	var encrypter wxencrypter.Encrypter

	if securityMode {

		encrypter, err = wxencrypter.NewEncrypter(token, encodingAesKey, appId)

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

func (w *Wechat) CheckSignature(timestamp, nonce, signature string) bool {
	arr := []string{w.Token, timestamp, nonce}
	sort.Strings(arr)

	joinStr := strings.Join(arr, "")

	h := sha1.New()
	io.WriteString(h, joinStr)
	b := h.Sum(nil)
	_signed := hex.EncodeToString(b)

	return _signed == signature
}

func (w *Wechat) ServeHTTP(rw http.ResponseWriter, r *http.Request) {
	timestamp := r.URL.Query().Get("timestamp")
	nonce := r.URL.Query().Get("nonce")
	signature := r.URL.Query().Get("signature")

	if !w.CheckSignature(imestamp, nonce, signature) {
		rw.WriteHeader(400)
		return
	}

	if r.Method == http.MethodGet {

		echostr := r.URL.Query().Get("echostr")
		rw.Write([]byte(echostr))
		return

	} else if r.Method == http.MethodPost {

		data, err := ioutil.ReadAll(r.Body)

		if err != nil {
			rw.WriteHeader(500)
			return
		}

		if w.securityMode {
			msgSignature := r.URL.Query().Get("msg_encrypt")
			data, err = w.Decrypt(msgSignature, timestamp, nonce, data)
			if err != nil {
				rw.WriteHeader(500)
				return
			}
		}

		rm, err := w.Unmarshal(data)
		if err != nil {
			rw.WriteHeader(500)
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

func (w *Wechat) Unmarshal(d []byte) (rm RequestMessage, err error) {
	err = xml.Unmarshal(d, &rm)
	return
}
