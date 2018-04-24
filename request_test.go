package wechat

import (
	"strings"
	"testing"
)

const (
	textMessage = `<xml>
 <ToUserName><![CDATA[toUser]]></ToUserName>
 <FromUserName><![CDATA[fromUser]]></FromUserName>
 <CreateTime>1348831860</CreateTime>
 <MsgType><![CDATA[text]]></MsgType>
 <Content><![CDATA[this is a test]]></Content>
 <MsgId>1234567890123456</MsgId>
 </xml>`

	imageMessage = `<xml>
 <ToUserName><![CDATA[toUser]]></ToUserName>
 <FromUserName><![CDATA[fromUser]]></FromUserName>
 <CreateTime>1348831860</CreateTime>
 <MsgType><![CDATA[image]]></MsgType>
 <PicUrl><![CDATA[this is a url]]></PicUrl>
 <MediaId><![CDATA[media_id]]></MediaId>
 <MsgId>1234567890123456</MsgId>
 </xml>`

	voiceMessage = `<xml>
<ToUserName><![CDATA[toUser]]></ToUserName>
<FromUserName><![CDATA[fromUser]]></FromUserName>
<CreateTime>1357290913</CreateTime>
<MsgType><![CDATA[voice]]></MsgType>
<MediaId><![CDATA[media_id]]></MediaId>
<Format><![CDATA[Format]]></Format>
<MsgId>1234567890123456</MsgId>
</xml>`

	voiceMessageWithRecongnition = `<xml>
<ToUserName><![CDATA[toUser]]></ToUserName>
<FromUserName><![CDATA[fromUser]]></FromUserName>
<CreateTime>1357290913</CreateTime>
<MsgType><![CDATA[voice]]></MsgType>
<MediaId><![CDATA[media_id]]></MediaId>
<Format><![CDATA[Format]]></Format>
<Recognition><![CDATA[腾讯微信团队]]></Recognition>
<MsgId>1234567890123456</MsgId>
</xml>`

	videoMessage = `<xml>
<ToUserName><![CDATA[toUser]]></ToUserName>
<FromUserName><![CDATA[fromUser]]></FromUserName>
<CreateTime>1357290913</CreateTime>
<MsgType><![CDATA[video]]></MsgType>
<MediaId><![CDATA[media_id]]></MediaId>
<ThumbMediaId><![CDATA[thumb_media_id]]></ThumbMediaId>
<MsgId>1234567890123456</MsgId>
</xml>`

	shortVedioMessage = `<xml>
<ToUserName><![CDATA[toUser]]></ToUserName>
<FromUserName><![CDATA[fromUser]]></FromUserName>
<CreateTime>1357290913</CreateTime>
<MsgType><![CDATA[shortvideo]]></MsgType>
<MediaId><![CDATA[media_id]]></MediaId>
<ThumbMediaId><![CDATA[thumb_media_id]]></ThumbMediaId>
<MsgId>1234567890123456</MsgId>
</xml>`

	locationMessage = `<xml>
<ToUserName><![CDATA[toUser]]></ToUserName>
<FromUserName><![CDATA[fromUser]]></FromUserName>
<CreateTime>1351776360</CreateTime>
<MsgType><![CDATA[location]]></MsgType>
<Location_X>23.134521</Location_X>
<Location_Y>113.358803</Location_Y>
<Scale>20</Scale>
<Label><![CDATA[位置信息]]></Label>
<MsgId>1234567890123456</MsgId>
</xml>`

	linkMessage = `<xml>
<ToUserName><![CDATA[toUser]]></ToUserName>
<FromUserName><![CDATA[fromUser]]></FromUserName>
<CreateTime>1351776360</CreateTime>
<MsgType><![CDATA[link]]></MsgType>
<Title><![CDATA[公众平台官网链接]]></Title>
<Description><![CDATA[公众平台官网链接]]></Description>
<Url><![CDATA[url]]></Url>
<MsgId>1234567890123456</MsgId>
</xml>`

	subscribeEventMessage = `<xml>
<ToUserName><![CDATA[toUser]]></ToUserName>
<FromUserName><![CDATA[FromUser]]></FromUserName>
<CreateTime>123456789</CreateTime>
<MsgType><![CDATA[event]]></MsgType>
<Event><![CDATA[subscribe]]></Event>
</xml>`

	unsubscribeEventMessage = `<xml>
<ToUserName><![CDATA[toUser]]></ToUserName>
<FromUserName><![CDATA[FromUser]]></FromUserName>
<CreateTime>123456789</CreateTime>
<MsgType><![CDATA[event]]></MsgType>
<Event><![CDATA[unsubscribe]]></Event>
</xml>`

	scanSubscribeEventMessage = `<xml><ToUserName><![CDATA[toUser]]></ToUserName>
<FromUserName><![CDATA[FromUser]]></FromUserName>
<CreateTime>123456789</CreateTime>
<MsgType><![CDATA[event]]></MsgType>
<Event><![CDATA[subscribe]]></Event>
<EventKey><![CDATA[qrscene_123123]]></EventKey>
<Ticket><![CDATA[TICKET]]></Ticket>
</xml>`

	scanEventMessage = `<xml>
<ToUserName><![CDATA[toUser]]></ToUserName>
<FromUserName><![CDATA[FromUser]]></FromUserName>
<CreateTime>123456789</CreateTime>
<MsgType><![CDATA[event]]></MsgType>
<Event><![CDATA[SCAN]]></Event>
<EventKey><![CDATA[SCENE_VALUE]]></EventKey>
<Ticket><![CDATA[TICKET]]></Ticket>
</xml>`

	locationEventMessage = `<xml>
<ToUserName><![CDATA[toUser]]></ToUserName>
<FromUserName><![CDATA[fromUser]]></FromUserName>
<CreateTime>123456789</CreateTime>
<MsgType><![CDATA[event]]></MsgType>
<Event><![CDATA[LOCATION]]></Event>
<Latitude>23.137466</Latitude>
<Longitude>113.352425</Longitude>
<Precision>119.385040</Precision>
</xml>`

	menuClickEventMessage = `<xml>
<ToUserName><![CDATA[toUser]]></ToUserName>
<FromUserName><![CDATA[FromUser]]></FromUserName>
<CreateTime>123456789</CreateTime>
<MsgType><![CDATA[event]]></MsgType>
<Event><![CDATA[CLICK]]></Event>
<EventKey><![CDATA[EVENTKEY]]></EventKey>
</xml>`

	menuViewEventMessage = `<xml>
<ToUserName><![CDATA[toUser]]></ToUserName>
<FromUserName><![CDATA[FromUser]]></FromUserName>
<CreateTime>123456789</CreateTime>
<MsgType><![CDATA[event]]></MsgType>
<Event><![CDATA[VIEW]]></Event>
<EventKey><![CDATA[www.qq.com]]></EventKey>
<MenuId>12312</MenuId>
</xml>`

	scancodePushEventMessage = `<xml><ToUserName><![CDATA[toUser]]></ToUserName>
<FromUserName><![CDATA[oMgHVjngRipVsoxg6TuX3vz6glDg]]></FromUserName>
<CreateTime>1408090502</CreateTime>
<MsgType><![CDATA[event]]></MsgType>
<Event><![CDATA[scancode_push]]></Event>
<EventKey><![CDATA[6]]></EventKey>
<ScanCodeInfo><ScanType><![CDATA[qrcode]]></ScanType>
<ScanResult><![CDATA[1]]></ScanResult>
</ScanCodeInfo>
</xml>`

	scancodeWaitmsgEventMessage = `<xml><ToUserName><![CDATA[toUser]]></ToUserName>
<FromUserName><![CDATA[oMgHVjngRipVsoxg6TuX3vz6glDg]]></FromUserName>
<CreateTime>1408090606</CreateTime>
<MsgType><![CDATA[event]]></MsgType>
<Event><![CDATA[scancode_waitmsg]]></Event>
<EventKey><![CDATA[6]]></EventKey>
<ScanCodeInfo><ScanType><![CDATA[qrcode]]></ScanType>
<ScanResult><![CDATA[2]]></ScanResult>
</ScanCodeInfo>
</xml>`

	picSysphotoEventMessage = `<xml><ToUserName><![CDATA[toUser]]></ToUserName>
<FromUserName><![CDATA[oMgHVjngRipVsoxg6TuX3vz6glDg]]></FromUserName>
<CreateTime>1408090651</CreateTime>
<MsgType><![CDATA[event]]></MsgType>
<Event><![CDATA[pic_sysphoto]]></Event>
<EventKey><![CDATA[6]]></EventKey>
<SendPicsInfo><Count>1</Count>
<PicList><item><PicMd5Sum><![CDATA[1b5f7c23b5bf75682a53e7b6d163e185]]></PicMd5Sum>
</item>
</PicList>
</SendPicsInfo>
</xml>`

	picPhotoOrAlbumEventMessage = `<xml><ToUserName><![CDATA[toUser]]></ToUserName>
<FromUserName><![CDATA[oMgHVjngRipVsoxg6TuX3vz6glDg]]></FromUserName>
<CreateTime>1408090816</CreateTime>
<MsgType><![CDATA[event]]></MsgType>
<Event><![CDATA[pic_photo_or_album]]></Event>
<EventKey><![CDATA[6]]></EventKey>
<SendPicsInfo><Count>1</Count>
<PicList><item><PicMd5Sum><![CDATA[5a75aaca956d97be686719218f275c6b]]></PicMd5Sum>
</item>
</PicList>
</SendPicsInfo>
</xml>`

	picWeixinEventMessage = `<xml><ToUserName><![CDATA[toUser]]></ToUserName>
<FromUserName><![CDATA[oMgHVjngRipVsoxg6TuX3vz6glDg]]></FromUserName>
<CreateTime>1408090816</CreateTime>
<MsgType><![CDATA[event]]></MsgType>
<Event><![CDATA[pic_weixin]]></Event>
<EventKey><![CDATA[6]]></EventKey>
<SendPicsInfo><Count>1</Count>
<PicList><item><PicMd5Sum><![CDATA[5a75aaca956d97be686719218f275c6b]]></PicMd5Sum>
</item>
</PicList>
</SendPicsInfo>
</xml>`

	locationSelectEventMessage = `<xml><ToUserName><![CDATA[toUser]]></ToUserName>
<FromUserName><![CDATA[oMgHVjngRipVsoxg6TuX3vz6glDg]]></FromUserName>
<CreateTime>1408091189</CreateTime>
<MsgType><![CDATA[event]]></MsgType>
<Event><![CDATA[location_select]]></Event>
<EventKey><![CDATA[6]]></EventKey>
<SendLocationInfo><Location_X><![CDATA[23]]></Location_X>
<Location_Y><![CDATA[113]]></Location_Y>
<Scale><![CDATA[15]]></Scale>
<Label><![CDATA[xxff]]></Label>
<Poiname><![CDATA[]]></Poiname>
</SendLocationInfo>
</xml>`

	templateSendJobFinishEventSuccessMessage = `<xml>
   <ToUserName><![CDATA[toUser]]></ToUserName>
   <FromUserName><![CDATA[oia2TjuEGTNoeX76QEjQNrcURxG8]]></FromUserName>
   <CreateTime>1395658920</CreateTime>
   <MsgType><![CDATA[event]]></MsgType>
   <Event><![CDATA[TEMPLATESENDJOBFINISH]]></Event>
   <MsgID>200163836</MsgID>
   <Status><![CDATA[success]]></Status>
   </xml>`

	templateSendJobFinishEventUserBlockMessage = `<xml>
   <ToUserName><![CDATA[toUser]]></ToUserName>
   <FromUserName><![CDATA[oia2TjuEGTNoeX76QEjQNrcURxG8]]></FromUserName>
   <CreateTime>1395658984</CreateTime>
   <MsgType><![CDATA[event]]></MsgType>
   <Event><![CDATA[TEMPLATESENDJOBFINISH]]></Event>
   <MsgID>200163840</MsgID>
   <Status><![CDATA[failed:user block]]></Status>
   </xml>`

	templateSendJobFinishEventSystemFailedMessage = `<xml>
   <ToUserName><![CDATA[toUser]]></ToUserName>
   <FromUserName><![CDATA[oia2TjuEGTNoeX76QEjQNrcURxG8]]></FromUserName>
   <CreateTime>1395658984</CreateTime>
   <MsgType><![CDATA[event]]></MsgType>
   <Event><![CDATA[TEMPLATESENDJOBFINISH]]></Event>
   <MsgID>200163841</MsgID>
   <Status><![CDATA[failed: system failed]]></Status>
   </xml>`
)

func TestTextMessage(t *testing.T) {

	msgs := []string{
		textMessage,
		imageMessage,
		voiceMessage,
		voiceMessageWithRecongnition,
		videoMessage,
		shortVedioMessage,
		locationMessage,
		linkMessage,
		subscribeEventMessage,
		unsubscribeEventMessage,
		scanSubscribeEventMessage,
		scanEventMessage,
		locationEventMessage,
		menuClickEventMessage,
		menuViewEventMessage,
		scancodePushEventMessage,
		scancodeWaitmsgEventMessage,
		picPhotoOrAlbumEventMessage,
		picSysphotoEventMessage,
		picWeixinEventMessage,
		locationSelectEventMessage,
		templateSendJobFinishEventSuccessMessage,
		templateSendJobFinishEventSystemFailedMessage,
		templateSendJobFinishEventUserBlockMessage,
	}

	for _, msg := range msgs {
		dft := &defaultRequestMessage{}
		err := dft.Unmarshal([]byte(msg))
		if err != nil {
			t.Fatal(err)
		}

		if dft.ToUserName() != "toUser" {
			t.Fatal("username is error")
		}

		if dft.MsgType() == ImageType {
			if dft.PicUrl() != "this is a url" {
				t.Fatal("image's picurl is not correct")
			}
		}

		if dft.MsgType() == EventType && dft.EventType() == ScancodeWaitmsgEventType {
			if dft.EventKey() != "6" ||
				dft.ScanCodeInfo().ScanResult != "2" ||
				dft.ScanCodeInfo().ScanType != "qrcode" {
				t.Fatal("ScancodeWaitmsgEventType value is not correct")
			}
		}

		if dft.MsgType() == EventType && dft.EventType() == ScancodePushEventType {
			if dft.EventKey() != "6" ||
				dft.ScanCodeInfo().ScanResult != "1" ||
				dft.ScanCodeInfo().ScanType != "qrcode" {
				t.Fatal("ScancodeWaitmsgEventType value is not correct")
			}
		}

		if dft.MsgType() == EventType && dft.EventType() == PicSysphotoEventType {
			if dft.EventKey() != "6" ||
				dft.SendPicsInfo().Count != 1 ||
				len(dft.SendPicsInfo().PicList) != 1 {
				t.Fatal("PicSysphotoEventType value is not correct")
			}

			if dft.SendPicsInfo().PicList[0].Item.PicMd5Sum != "1b5f7c23b5bf75682a53e7b6d163e185" {
				t.Fatal("PicSysphotoEventType piclist item is not correct")
			}
		}

		if dft.MsgType() == EventType && dft.EventType() == PicPhotoOrAlbumEventType || dft.MsgType() == EventType && dft.EventType() == PicWeixinEventType {
			if dft.EventKey() != "6" ||
				dft.SendPicsInfo().Count != 1 ||
				len(dft.SendPicsInfo().PicList) != 1 {
				t.Fatal("PicPhotoOrAlbumEventType value is not correct")
			}

			if dft.SendPicsInfo().PicList[0].Item.PicMd5Sum != "5a75aaca956d97be686719218f275c6b" {
				t.Fatal("PicPhotoOrAlbumEventType piclist item is not correct")
			}
		}

		if dft.MsgType() == EventType && dft.EventType() == LocationSelectEvenType {
			if dft.EventKey() != "6" {
				t.Fatal("LocationSelectEvenType value is not correct")
			}

			if dft.SendLocationInfo().Location_X != 23 ||
				dft.SendLocationInfo().Location_Y != 113 ||
				dft.SendLocationInfo().Label != "xxff" ||
				dft.SendLocationInfo().Scale != 15 ||
				dft.SendLocationInfo().Poiname != "" {
				t.Fatal("LocationSelectEvenType SendLocationInfo is not correct")
			}
		}

		if dft.MsgType() == EventType && dft.EventType() == TemplateSendJobFinishEventType {
			if dft.MsgId() == 200163836 && !strings.Contains(dft.Status(), "success") {
				t.Fatal("TemplateSendJobFinishEventType success is not correct")
			}

			if dft.MsgId() == 200163840 && !strings.Contains(dft.Status(), "user") {
				t.Fatal("TemplateSendJobFinishEventType user is not correct")
			}

			if dft.MsgId() == 200163841 && !strings.Contains(dft.Status(), "system") {
				t.Fatal("TemplateSendJobFinishEventType system is not correct")
			}
		}
	}

}
