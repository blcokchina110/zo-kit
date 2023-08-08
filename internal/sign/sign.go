package sign

import (
	"crypto/hmac"
	"crypto/md5"
	"crypto/sha256"
	"encoding/hex"
	"net/url"
	"reflect"
	"strconv"
	"strings"
)

// 验证签名是否正确
func VerificationSign(data interface{}, apiSecret, signValue string) (bool, string) {
	urlValeus := json2UrlValues(data)
	query, err := url.QueryUnescape(urlValeus.Encode())
	if err != nil || query == "" {
		return false, ""
	}
	query = query + "&key=" + apiSecret
	signResult := sign(query)

	return signValue == signResult, signResult
}

// 签名
func Sign(query string) string {
	return sign(query)
}

// 加密签名
func CryptoSign(param, encryption string) string {
	switch encryption {
	case "HMAC-SHA256":
		return hmacSign(param)
	case "MD5":
		return md5Sign(param)
	}
	return ""
}

// json转成url.Values
func json2UrlValues(data interface{}) url.Values {
	t := reflect.TypeOf(data)
	v := reflect.ValueOf(data)

	if t.Kind() == reflect.Ptr {
		t = t.Elem()
		v = v.Elem()
	}

	values := url.Values{}
	for i := 0; i < v.NumField(); i++ {
		//-和sign字段不参与签名
		if key := t.Field(i).Tag.Get("json"); key != "-" && key != "sign" {
			switch v.Field(i).Kind() {
			case reflect.String:
				if val := v.Field(i).String(); val != "" {
					values.Set(key, val)
				}
			case reflect.Float64:
				if val := v.Field(i).Float(); val > 0 {
					values.Set(key, strconv.FormatFloat(val, 'f', -1, 64))
				}
			case reflect.Int:
				if val := v.Field(i).Int(); val > 0 {
					values.Set(key, strconv.Itoa(int(val)))
				}
			case reflect.Int64:
				if val := v.Field(i).Int(); val > 0 {
					//json为string，但类型为int64的
					if index := strings.Index(key, ",string"); index > -1 {
						key = key[:index]
					}
					values.Set(key, strconv.FormatInt(val, 10))
				}
			}
		}
	}
	return values
}

// 签名
// 将字符串进行md5加密，并转换成32位base64值，并转换成大写
func sign(query string) string {
	return strings.ToUpper(md5Sign(query))
}

// md5加密
func md5Sign(param string) string {
	h := md5.New()
	h.Write([]byte(param))

	return hex.EncodeToString(h.Sum(nil))
}

// sha256加密
func hmacSign(param string) string {
	hash := hmac.New(sha256.New, []byte(param))
	hash.Write([]byte(param))

	return hex.EncodeToString(hash.Sum(nil))
}
