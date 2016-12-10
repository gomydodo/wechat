package wechat

import (
	"bytes"
	"encoding/xml"
	"fmt"
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"
	"time"

	"github.com/gomydodo/wxencrypter"
)

func genRequest() (req *http.Request, err error) {
	timestamp := fmt.Sprintf("%d", time.Now().Unix())
	nonce := "nonce"
	signature := w.Signature(timestamp, nonce)

	p, err := wxencrypter.NewPrpcrypt(w.EncodingAesKey)
	if err != nil {
		return
	}

	ret, err := p.Encrypt(w.AppID, []byte(textMessage))
	if err != nil {
		return
	}

	msgSignature := wxencrypter.Sha1(w.Token, timestamp, nonce, ret)

	xmlData := wxencrypter.EncryptedRequestXML{
		ToUserName: "testname",
		Encrypt:    ret,
	}

	data, err := xml.Marshal(xmlData)
	if err != nil {
		return
	}

	v := url.Values{}
	v.Add("timestamp", timestamp)
	v.Add("nonce", nonce)
	v.Add("signature", signature)
	v.Add("msg_signature", msgSignature)

	req = httptest.NewRequest("POST", "http://1.2.3.4?"+v.Encode(), bytes.NewReader(data))
	return
}

func TestRoute(t *testing.T) {
	dft := &defaultRequestMessage{}
	err := dft.Unmarshal([]byte(textMessage))
	if err != nil {
		t.Fatal(err)
	}

	route := NewRouter()

	route.Add(dft.MsgType(), "", Route{
		MsgType: dft.MsgType(),
		Key:     "",
		Handler: func(c Context) error { return c.Response().Text("test") },
	})

	req, err := genRequest()
	if err != nil {
		t.Fatal(err)
	}

	rw := httptest.NewRecorder()

	handler := func(rw http.ResponseWriter, r *http.Request) {
		c := newContext(rw, req, w)
		err := c.parse()

		if err != nil {
			t.Fatal(err)
		}

		route.Find(c)

		c.Handler()(c)
	}

	handler(rw, req)
	if rw.Body.Len() < 0 {
		t.Fatal()
	}
}
