package vx

import (
	"crypto/hmac"
	"crypto/sha256"
	"fmt"
	"github.com/shopspring/decimal"
	"io"
	"math"
	"reflect"
	"strconv"
	"strings"
	"time"
)

// 获取hmac-hs256 签名
func hmacHs256(message string, secret string) string {
	h := hmac.New(sha256.New, []byte(secret))
	io.WriteString(h, message)
	return fmt.Sprintf("%x", h.Sum(nil))
}

func MapToParam(m map[string]interface{}, escapeValue ...string) string{
	var keys = make([]string,0)
	for k,_ :=range m {
		keys = append(keys, k)
	}
	var tmp string
	for i:=0;i<len(keys)-1;i++ {
		for j:=i+1; j< len(keys);j++{
			if keys[i] > keys[j] {
				tmp = keys[i]
				keys[i] = keys[j]
				keys[j] = tmp
			}
		}
	}
	var records = make([]string,0)
L :
	for i,v:=range keys {
		if IfZero(m[v]) {
			continue
		}
		for _,v2 := range escapeValue {
			if v == v2 {
				continue L
			}
		}
		records = append(records,fmt.Sprintf("%s=%v", keys[i], ToString(m[v])))
	}
	return strings.Join(records, "&")
}

// To judge a value whether zero or not.
// By the way, '%' '%%' is regarded as zero.
func IfZero(arg interface{}) bool {
	if arg == nil {
		return true
	}
	switch v := arg.(type) {
	case int, int32, int16, int64:
		if v == 0 {
			return true
		}
	case float32:
		r := float64(v)
		return math.Abs(r-0) < 0.0000001
	case float64:
		return math.Abs(v-0) < 0.0000001
	case string:
		if v == "" || v == "%%" || v == "%" {
			return true
		}
	case *string, *int, *int64, *int32, *int16, *int8, *float32, *float64, *time.Time:
		if v == nil {
			return true
		}
	case time.Time:
		return v.IsZero()
	case decimal.Decimal:
		tmp, _ := v.Float64()
		return math.Abs(tmp-0) < 0.0000001
	default:
		return false
	}
	return false
}

// Change arg to string
// if arg is a ptr kind, then change what it points to  to string
func ToString(arg interface{}) string {
	tmp := reflect.Indirect(reflect.ValueOf(arg)).Interface()
	switch v := tmp.(type) {
	case int:
		return strconv.Itoa(v)
	case int8:
		return strconv.FormatInt(int64(v), 10)
	case int16:
		return strconv.FormatInt(int64(v), 10)
	case int32:
		return strconv.FormatInt(int64(v), 10)
	case int64:
		return strconv.FormatInt(v, 10)
	case string:
		return v
	case float32:
		return strconv.FormatFloat(float64(v), 'f', -1, 32)
	case float64:
		return strconv.FormatFloat(v, 'f', -1, 64)
	case time.Time:
		return v.Format("2006-01-02 15:04:05")
	case fmt.Stringer:
		return v.String()
	default:
		return ""
	}
}
