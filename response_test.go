package wechat

import (
	"encoding/xml"
	"testing"
)

func TestArticleResponseMessage(t *testing.T) {
	item := NewArticleItem("title", "description", "picUrl", "url")
	item1 := NewArticleItem("1", "1", "1", "1111111111")
	news, err := NewArticleResponseMessage("to", "from", item, item1)
	if err != nil {
		t.Fatal("news create error: ", err)
	}

	b, err := xml.Marshal(news)
	if err != nil {
		t.Fatal("marshal news error: ", err)
	}

	var news1 ArticleResponseMessage
	err = xml.Unmarshal(b, &news1)
	if err != nil {
		t.Fatal("unmarshal news error: ", err)
	}

	if news1.ArticleCount != len(news1.Articles) {
		t.Log(news1)
		t.Fatalf("count not equal to article's count, should be %d, actual is %d", news1.ArticleCount, len(news1.Articles))
	}
}

func TestImageResponseMessage(t *testing.T) {
	image := NewImageResponseMessage("to", "from", "mediaId")

	b, err := xml.Marshal(image)
	if err != nil {
		t.Fatal("marshal error: ", err)
	}

	t.Log(string(b))

	var image1 ImageResponseMessage
	err = xml.Unmarshal(b, &image1)
	if err != nil {
		t.Fatal("unmarshal error:", err)
	}

}
