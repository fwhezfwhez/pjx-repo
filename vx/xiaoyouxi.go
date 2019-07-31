package vx

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"runtime/debug"
	"time"

	"github.com/fwhezfwhez/errorx"
)

var c = &http.Client{
	Timeout: 30 * time.Second,
}

var Mode string
func init() {
	log.SetFlags(log.Llongfile | log.LstdFlags)
}

// 检查sessionKey是否失效
func CheckSessionKey(cli *http.Client, openId, accessToken, sessionKey string) (bool, error) {
	if cli == nil {
		cli = c
	}
	signature := hmacHs256("", sessionKey)

	var result = struct {
		Errcode int    `json:"errcode"`
		Errmsg  string `json:"errmsg"`
	}{}
	req, e := http.NewRequest("GET", fmt.Sprintf("https://api.weixin.qq.com/wxa/checksession?access_token=%s&signature=%s&openid=%s&sig_method=hmac_sha256", accessToken, signature, openId), nil)
	if e != nil {
		return false, errorx.Wrap(e)
	}
	rsp, e := cli.Do(req)
	if e != nil {
		return false, errorx.Wrap(e)
	}
	if rsp != nil && rsp.Body != nil {
		defer rsp.Body.Close()
	}
	buf, e := ioutil.ReadAll(rsp.Body)
	if e != nil {
		return false, errorx.Wrap(e)
	}
	e = json.Unmarshal(buf, &result)
	if e != nil {
		return false, errorx.Wrap(e)
	}
	// 正常
	if result.Errcode == 0 && result.Errmsg == "ok" {
		return true, nil
	}
	// 非法
	if result.Errcode == 87009 || result.Errmsg == "invalid signature" {
		fmt.Println(result)
		return false, nil
	}
	return false, errorx.NewFromStringf("got errcode '%d', errmsg '%s'", result.Errcode, result.Errmsg)
}

type MidasPayRequest struct {
	Openid      string `json:"openid"`
	Appid       string `json:"appid"`
	OfferId     string `json:"offer_id"`
	Ts          int64  `json:"ts"`
	ZoneId      string `json:"zone_id"`
	Pf          string `json:"pf"`
	Amt         int    `json:"amt"`     // 不能为0
	BillNo      string `json:"bill_no"` // 唯一
	Sig         string `json:"sig"`
	AccessToken string `json:"access_token,omitempty"` //请求时消除
	MpSig       string `json:"mp_sig"`
}

type MidasPayResponse struct {
	Errcode    int    `json:"errcode"`
	Errmsg     string `json:"errmsg"`
	BillNo     string `json:"bill_no"`
	Balance    int    `json:"balance"`      // 预扣后的余额
	UsedGenAmt int    `json:"used_gen_amt"` // 本次扣的赠送币的金额
}

// 减钻石
func MidasPay(param MidasPayRequest, url string, orgLoc string, offerSecret string, sessionKey string) error {
	// 构造sig
	var e error
	param.Sig, param.MpSig, e = GenerateSigAndMpSig(param, offerSecret, sessionKey, orgLoc, param.AccessToken)
	// 调用远程接口
	param.AccessToken = ""
	buf, e := json.Marshal(param)
	if e != nil {
		return errorx.NewWithParam(e, param)
	}
	req, e := http.NewRequest("POST", url, bytes.NewReader(buf))

	logbuf, e := json.MarshalIndent(map[string]interface{}{"url": url, "param": param}, "  ", "  ")
	log.SetFlags(log.Llongfile | log.LstdFlags)
	log.Println(string(logbuf))

	req.Header.Set("Content-Type", "application/json")
	if e != nil {
		return errorx.Wrap(e)
	}
	rsp, e := c.Do(req)
	if e != nil {
		return errorx.Wrap(e)
	}
	if rsp != nil && rsp.Body != nil {
		defer rsp.Body.Close()
	}
	rsBuf, e := ioutil.ReadAll(rsp.Body)
	if e != nil {
		return errorx.Wrap(e)
	}
	var rs MidasPayResponse
	e = json.Unmarshal(rsBuf, &rs)
	if e != nil {
		return errorx.Wrap(e)
	}
	if rs.Errcode != 0 {
		return errorx.NewFromStringf("got errcode '%d', errmsg '%s'", rs.Errcode, rs.Errmsg)
	}
	return nil
}

// 获取sig和mp_sig的公共参数
type SigParam struct {
	Openid      string `json:"openid"`
	Appid       string `json:"appid"`
	OfferId     string `json:"offer_id"`
	Ts          int64  `json:"ts"`
	ZoneId      string `json:"zone_id"`
	Pf          string `json:"pf"`
	Sig         string `json:"sig"`
	AccessToken string `json:"access_token,omitempty"`
	MpSig       string `json:"mp_sig"`
}

type GetDiamondNumberParam struct {
	Openid      string `json:"openid"`
	Appid       string `json:"appid"`
	OfferId     string `json:"offer_id"`
	Ts          int64  `json:"ts"`
	ZoneId      string `json:"zone_id"`
	Pf          string `json:"pf"`
	Sig         string `json:"sig"`
	AccessToken string `json:"access_token,omitempty"`
	MpSig       string `json:"mp_sig"`
}
type GetDiamondNumberParamResult struct {
	Errcode    int    `json:"errcode"`
	Errmsg     string `json:"errmsg"`
	Balance    int64  `json:"balance"`     // 游戏币个数
	GenBalance int64  `json:"gen_balance"` // 赠送的游戏币数量
	FirstSave  int    `json:"first_save"`  // 是否历史首冲
	SaveAmt    int64  `json:"save_atm"`    // 历史充值金额的游戏币数量
	SaveSum    int64  `json:"save_sum"`    // 历史总游戏币金额
	CostSum    int64  `json:"cost_sum"`    // 历史总消费游戏币金额
	PresentSum int64  `json:"present_sum"` // 历史累计的赠送金额
}

// 获取钻石余额
// 包含了检查sessionKey是否过期
func GetBalance(openid string, result *GetDiamondNumberParamResult, accessToken string, offerId string, offerSecret, appId string, sessionKey string, orgLoc string) (string, error) {
	if sessionKey == "" {
		return "", errorx.NewFromString("sessionKey is empty")
	}
	// 检查sessionKey是否过期
	ok, e := CheckSessionKey(&http.Client{Timeout: 15 * time.Second}, openid, accessToken, sessionKey)
	if e != nil {
		return "", errorx.Wrap(e)
	}
	if !ok {
		return "sessionKey expires", errorx.NewFromStringf("open_id '%s'sessionKey exipires", openid)
	}
	log.Println("accessToken:", accessToken)
	log.Println("offerId:", offerId)
	log.Println("appId:", appId)
	var ts = time.Now().Unix()

	// 请求参数
	param := GetDiamondNumberParam{
		Openid:      openid,
		Appid:       appId,
		OfferId:     offerId,
		Ts:          ts,
		ZoneId:      "1",
		Pf:          "android",
		Sig:         "",
		MpSig:       "",
		AccessToken: accessToken,
	}

	param.Sig, param.MpSig, e = GenerateSigAndMpSig(param, offerSecret, sessionKey, orgLoc, accessToken)
	if e != nil {
		return "", errorx.Wrap(e)
	}

	// 清理param里的access_token
	param.AccessToken = ""
	// 请求余额
	buf, e := json.Marshal(param)
	if e != nil {
		return "", errorx.Wrap(e)
	}

	if e != nil {
		return "", errorx.Wrap(e)
	}
    url := fmt.Sprintf("https://api.weixin.qq.com%s?access_token=%s", orgLoc, accessToken)
	req, e := http.NewRequest("POST", url, bytes.NewReader(buf))
	if e != nil {
		return "", errorx.Wrap(e)
	}
	req.Header.Set("Content-Type", "application/json")
	resp, e := c.Do(req)
	if e != nil {
		return "", errorx.Wrap(e)
	}
	if resp != nil && resp.Body != nil {
		defer resp.Body.Close()
	}

	reBuf, e := ioutil.ReadAll(resp.Body)
	if e != nil {
		return "", errorx.Wrap(e)
	}

	e = json.Unmarshal(reBuf, result)
	if e != nil {
		saveError(errorx.Wrap(e), map[string]interface{}{
			"orgLoc": orgLoc,
			"url": url,
			"response": string(reBuf),
			"param": param,
		})
		return "", nil
	}

	if result.Errcode != 0 || result.Errmsg != "ok" {
		return "", errorx.NewFromStringf("获取余额结果，收到 errcode '%d', errmsg '%s'", result.Errcode, result.Errmsg)
	}
	return "", nil
}

// 获取sig与mp_sig
// 使用实例:
///*
//	sig, mpSig, e:= GenerateSigAndMpSig(SigParam{
//		Openid: "odkx20ENSNa2w5y3g_qOkOvBNM1g",
//		Appid:"wx1234567",
//		OfferId: "12345678",
//		Ts:1507530737,
//		ZoneId:"1",
//		Pf:"android",
//		AccessToken: "ACCESSTOKEN",
//	},"zNLgAGgqsEWJOg1nFVaO5r7fAlIQxr1u", "V7Q38/i2KXaqrQyl2Yx9Hg==","/cgi-bin/midas/getbalance")
//*/
func GenerateSigAndMpSig(param interface{}, offerSecret string, sessionKey string, orgLoc string, accessToken string) (string, string, error) {
	var sig string
	var mpSig string
	buf, _ := json.Marshal(param)
	var m map[string]interface{}
	json.Unmarshal(buf, &m)
	delete(m, "access_token")

	stringA := MapToParam(m)
	stringSignTemp := stringA + fmt.Sprintf("&org_loc=%s&method=POST&secret=%s", orgLoc, offerSecret)

	sig = hmacHs256(stringSignTemp, offerSecret)

	// 获取mp_sig
	m["sig"] = sig
	m["access_token"] = accessToken
	stringB := MapToParam(m)
	stringSignTempB := stringB + fmt.Sprintf("&org_loc=%s&method=POST&session_key=%s", orgLoc, sessionKey)
	log.Println("GenerateSigAndMpSig.orgLoc", orgLoc)

	log.Println("stringA", stringA)
	log.Println("stringB", stringB)
	log.Println("stringSignTempA", stringSignTemp)
	log.Println("stringSignTempB", stringSignTempB)
	mpSig = hmacHs256(stringSignTempB, sessionKey)
	return sig, mpSig, nil
}

type MidasPresentRequest struct {
	Openid        string `json:"openid"`
	Appid         string `json:"appid"`
	OfferId       string `json:"offer_id"`
	Ts            int64  `json:"ts"`
	ZoneId        string `json:"zone_id"`
	Pf            string `json:"pf"`
	BillNo        string `json:"bill_no"`
	PresentCounts int64  `json:"present_counts"`
	Sig           string `json:"sig"`
	AccessToken   string `json:"access_token,omitempty"`
	MpSig         string `json:"mp_sig"`
}
type MidasPresentResponse struct {
	Errcode int    `json:"errcode"`
	Errmsg  string `json:"errmsg"`
	Balance int64  `json:"balance"`
	BillNo  string `json:"bill_no"`
}

func MidasPresent(param MidasPresentRequest, url string, orgLoc string, offerSecret string, sessionKey string) (int64, error) {
	log.Println("MidasPresnt.orgLoc", orgLoc)
	var e error
	param.Sig, param.MpSig, e = GenerateSigAndMpSig(param, offerSecret, sessionKey, orgLoc, param.AccessToken)
	param.AccessToken = ""
	buf, e := json.Marshal(param)
	if e != nil {
		return 0, errorx.NewWithParam(e, param)
	}
	req, e := http.NewRequest("POST", url, bytes.NewReader(buf))

	logbuf, e := json.MarshalIndent(map[string]interface{}{"url": url, "param": param}, "  ", "  ")
	log.SetFlags(log.Llongfile | log.LstdFlags)
	log.Println(string(logbuf))

	req.Header.Set("Content-Type", "application/json")
	if e != nil {
		return 0, errorx.Wrap(e)
	}
	rsp, e := c.Do(req)
	if e != nil {
		return 0, errorx.Wrap(e)
	}
	if rsp != nil && rsp.Body != nil {
		defer rsp.Body.Close()
	}
	rsBuf, e := ioutil.ReadAll(rsp.Body)
	if e != nil {
		return 0, errorx.Wrap(e)
	}
	var rs MidasPresentResponse
	e = json.Unmarshal(rsBuf, &rs)
	if e != nil {
		return 0, errorx.Wrap(e)
	}
	if rs.Errcode != 0 {
		return 0, errorx.NewFromStringf("got errcode '%d', errmsg '%s'", rs.Errcode, rs.Errmsg)
	}
	log.Println("midasPresent result:", string(rsBuf))
	log.Println(rs.Balance)
	return rs.Balance, nil
}

func MsgSecCheck(accessToken, content string) (bool, error) {
	url := fmt.Sprintf("https://api.weixin.qq.com/wxa/msg_sec_check?access_token=%s", accessToken)
	buf, e := json.Marshal(struct{ Content string `json:"content"` }{Content: content})
	if e != nil {
		return false, errorx.Wrap(e)
	}
	req, e := http.NewRequest("POST", url, bytes.NewReader(buf))
	if e != nil {
		return false, errorx.Wrap(e)
	}
	req.Header.Set("Content-Type", "application/json")
	rsp, e := c.Do(req)
	if e != nil {
		return false, errorx.Wrap(e)
	}
	if rsp != nil && rsp.Body != nil {
		defer rsp.Body.Close()
	}
	rsBuf, e := ioutil.ReadAll(rsp.Body)
	if e != nil {
		return false, errorx.Wrap(e)
	}

	type Result struct {
		Errcode int    `json:"errcode"`
		ErrMsg  string `json:"errMsg"`
	}
	var result Result
	e = json.Unmarshal(rsBuf, &result)
	if e != nil {
		return false, errorx.Wrap(e)
	}
	if result.Errcode == 0 && result.ErrMsg == "ok" {
		return true, nil
	}
	if result.Errcode == 87014 && result.ErrMsg != "" {
		return false, nil
	}
	return false, errorx.NewFromStringf("got errcode '%d', errMsg '%s'", result.Errcode, result.ErrMsg)
}

func saveError(e error, context ... map[string]interface{}) string {
L:
	switch v := e.(type) {
	case errorx.Error:
		break L
	case error:
		return saveError(errorx.NewFromString(string(fmt.Sprintf("err '%s' \n %s", v.Error(), debug.Stack()))), context...)
	}

	if len(context) > 1 {
		panic("context max length 1")
	}

	var tmp map[string]interface{}
	if len(context) != 0 {
		tmp = context[0]
	}
	type ResultError struct {
		Id          int       `gorm:"column:id;default:" json:"id" form:"id"`
		Message     string    `gorm:"column:message;default:" json:"message" form:"message"`
		Keyword     string    `gorm:"column:keyword;default:" json:"keyword" form:"keyword"`
		Times       int       `gorm:"column:times;default:" json:"times" form:"times"`
		CreatedAt   time.Time `gorm:"column:created_at;default:" json:"created_at" form:"created_at"`
		CreatedDate time.Time `gorm:"column:created_date;default:" json:"created_date" form:"created_date"`
	}
	type Result1 struct {
		Message string      `json:"message"`
		RE      ResultError `json:"data"`
	}
	var result1 Result1

	var m = map[string]interface{}{
		"message": e.Error(),
		"request": tmp,
	}
	buf, e := json.Marshal(m)
	if e != nil {
		log.Println(errorx.Wrap(e).Error())
		return ""
	}
	log.Println("url", fmt.Sprintf("%s%s", "https://xyx.zonst.com/err", "/error/"))
	log.Println("buf", string(buf))
	req, e := http.NewRequest("POST", fmt.Sprintf("%s%s", "https://xyx.zonst.com/err", "/error/"), bytes.NewReader(buf))
	if e != nil {
		log.Println(errorx.Wrap(e).Error())
		return ""
	}
	req.Header.Set("Content-Type", "application/json")
	rsp, e := c.Do(req)
	if e != nil {
		log.Println(errorx.Wrap(e).Error())
		return ""
	}
	if rsp != nil && rsp.Body != nil {
		defer rsp.Body.Close()
	}
	rsBuf, e := ioutil.ReadAll(rsp.Body)
	if rsp.Status != "200 OK" {
		if e != nil {
			log.Println(errorx.Wrap(e).Error())
			return ""
		}
		log.Println("/error/ api throws " + string(rsBuf))
		return ""
	}
	e = json.Unmarshal(rsBuf, &result1)
	if e != nil {
		log.Println(errorx.Wrap(e).Error())
		return ""
	}
	fmt.Println(result1)
	return ""
}
