package vx

import (
	"fmt"
	"net/http"
	"testing"
)

// 测试accesstoken 与 sessionkey 是否过期
func TestCheckSessionKey(t *testing.T) {
	var accessToken = "19_z-28Zl3NFTMVXLcYls_CqS5PNAdlNc191tbzlmGEKbacShqVjYpGy7I7-vDeyD92HQ0PUlYwn_KwxFPV_X7IdiYS38wux4-M-lt4PXDmq_S72SpS65r1Ghs3uAY5deqkyjVqfb-Ycf3rpPsnXXRhAJAGGA"
	var openId = "oZKx35Nm5ztxziIKxmf6jMTmeOpY"
	var sessionKey = "YeiOaRfiJFmrr/0KVBaBQg=="
	ok, e := CheckSessionKey(&http.Client{}, openId, accessToken, sessionKey)
	fmt.Println(ok, e)
}

// 测试sig和mp_sig的生成
func TestGenerateSigAndMpSig(t *testing.T) {
	sig, mpSig, e := GenerateSigAndMpSig(SigParam{
		Openid:      "odkx20ENSNa2w5y3g_qOkOvBNM1g",
		Appid:       "wx1234567",
		OfferId:     "12345678",
		Ts:          1507530737,
		ZoneId:      "1",
		Pf:          "android",
		AccessToken: "ACCESSTOKEN",
	}, "zNLgAGgqsEWJOg1nFVaO5r7fAlIQxr1u", "V7Q38/i2KXaqrQyl2Yx9Hg==", "/cgi-bin/midas/getbalance", "ACCESSTOKEN")

	if e != nil {
		t.Fail()
		fmt.Println(e.Error())
	}
	fmt.Println("sig:", sig)
	fmt.Println("mp_sig", mpSig)
}

//func TestMidasPresent(t *testing.T) {
//	var midasPresentParam = MidasPresentRequest{
//		Openid:        "oZKx35Nm5ztxziIKxmf6jMTmeOpY",
//		Appid:         "wxbec7aebf80022eb2",
//		OfferId:       "1450019844",
//		Ts:            time.Now().Unix(),
//		ZoneId:        "1",
//		Pf:            "android",
//		BillNo:        fmt.Sprintf("%s-%d-%d-%d", time.Now().Format("20060102150405"), time.Now().UnixNano(), 33586765, 500),
//		PresentCounts: 500,
//		AccessToken:   ,
//		Sig:           "",
//		MpSig:         "",
//	}
//	e = MidasPresent(
//		midasPresentParam,
//		fmt.Sprintf("https://api.weixin.qq.com%s?access_token=%s", cfg.GetString("vx.midasPresent"), accessToken),
//		c.GetString("vx.midasPresent"),
//		c.GetString("vx.offerSecret"),
//		sessionKey)
//}

func TestMsgSecCheck(t *testing.T) {
	ok, e := MsgSecCheck("20_xGp5ZhQJOnXyWDQzZzJzqV1s9kVrUVoSiXdinHdT3OHMpIDireaFYawW4CDik8q_Xz8ya911hyxxXOYlmMNWtWxn0f6ciLbxG58iRQ6vEkZtIXeLze-7mXydkyQDRShADAMAB", "特3456书yuuo莞6543李zxcz蒜7782法fgnv级")
	fmt.Println(ok, e)
}
