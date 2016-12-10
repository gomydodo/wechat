package wechat

import (
	"encoding/xml"
	"log"
	"testing"

	"github.com/gomydodo/wxencrypter"
)

const (
	token          = "pamtest"
	appID          = "wxb11529c136998cb6"
	secret         = "xxxx"
	encodingAesKey = "abcdefghijklmnopqrstuvwxyz0123456789ABCDEFG"
)

var w *Wechat

func init() {
	var err error
	w, err = New(token, appID, encodingAesKey, secret, SubscriptionsType)
	if err != nil {
		log.Fatal("new wechat error :", err)
	}
}

func TestEncrypt(t *testing.T) {
	b, _ := xml.Marshal(NewTextResponseMessage("to", "from", "content"))

	b, err := w.Encrypt(b)
	if err != nil {
		t.Fatal("wechat encrypt error: ", err)
	}

	e, err := wxencrypter.ParseResponseXML(b)
	if err != nil {
		t.Fatal("parse xml error :", err)
	}

	_, err = w.Decrypt(e.MsgSignature, e.TimeStamp, e.Nonce, b)
	if err != nil {
		t.Fatal("decrypt data error: ", err)
	}

}

func TestNewTextResponseMessage(t *testing.T) {
	tm := NewTextResponseMessage("to", "from", "content")

	b, err := xml.Marshal(tm)
	if err != nil {
		t.Fatal("marshal text response message error: ", err)
	}

	dft := &defaultRequestMessage{}

	err = dft.Unmarshal(b)
	if err != nil {
		t.Fatal("unmarshal request message error: ", err)
	}

	if dft.ToUserName() != "to" {
		t.Fatalf("parse tousername erorr, should be %s, actual %s", "to", dft.ToUserName())
	}
}
