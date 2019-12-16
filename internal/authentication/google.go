/*
主要对外调用方法说明：
    GetSecret() ：获取秘钥（32位字符串）
    GetCode() ：获取动态码
    GetQrcode() ：获取动态码二维码内容
    GetQrcodeUrl() ：获取动态码二维码图片地址
    VerifyCode() ：验证动态码

    秘钥生成原理（基于时间）：

    1、时间戳，精确到微秒，除以1000，除以30（动态6位数字每30秒变化一次）

    2、对时间戳余数 hmac_sha1 编码

    3、然后 base32 encode 标准编码

    4、输出大写字符串，即秘钥

    动态6位数字验证：

    Google Authenticator 会基于密钥和时间计算一个 HMAC-SHA1 的 hash 值，这个 hash 是160 bit 的，然后将这个hash值随机取连续的4个字节生成32位整数，最后将整数取31位，再取模得到一个的整数。

    这个就是 Google Authenticator 显示的数字。

    在服务器端验证的时候，同样的方法来计算出数字，然后比较计算出来的结果和用户输入的是否一致。
*/

package authentication

import (
	"bytes"
	"crypto/hmac"
	"crypto/sha1"
	"encoding/base32"
	"encoding/binary"
	"fmt"
	"strings"
	"time"
)

type GoogleAuth struct {
}


func NewGoogleAuth() *GoogleAuth {
	return &GoogleAuth{}
}


func (this *GoogleAuth) un() int64 {
	return time.Now().UnixNano() / 1000 / 30
}


func (this *GoogleAuth) hmacSha1(key, data []byte) []byte {
	h := hmac.New(sha1.New, key)
	if total := len(data); total > 0 {
		h.Write(data)
	}
	return h.Sum(nil)
}


func (this *GoogleAuth) base32encode(src []byte) string {
	return base32.StdEncoding.EncodeToString(src)
}


func (this *GoogleAuth) base32decode(s string) ([]byte, error) {
	return base32.StdEncoding.DecodeString(s)
}


func (this *GoogleAuth) toBytes(value int64) []byte {
	var result []byte
	mask := int64(0xFF)
	shifts := [8]uint16{56, 48, 40, 32, 24, 16, 8, 0}
	for _, shift := range shifts {
		result = append(result, byte((value>>shift)&mask))
	}
	return result
}


func (this *GoogleAuth) toUint32(bts []byte) uint32 {
	return (uint32(bts[0]) << 24) + (uint32(bts[1]) << 16) +
		(uint32(bts[2]) << 8) + uint32(bts[3])
}


func (this *GoogleAuth) oneTimePassword(key []byte, data []byte) uint32 {
	hash := this.hmacSha1(key, data)
	offset := hash[len(hash)-1] & 0x0F
	hashParts := hash[offset : offset+4]
	hashParts[0] = hashParts[0] & 0x7F
	number := this.toUint32(hashParts)
	return number % 1000000
}


func (this *GoogleAuth) GetSecret() string {
	var buf bytes.Buffer
	binary.Write(&buf, binary.BigEndian, this.un())
	return strings.ToUpper(this.base32encode(this.hmacSha1(buf.Bytes(), nil)))
}


func (this *GoogleAuth) GetCode(secret string) (string, error) {
	secretUpper := strings.ToUpper(secret)
	secretKey, err := this.base32decode(secretUpper)
	if err != nil {
		return "", err
	}
	number := this.oneTimePassword(secretKey, this.toBytes(time.Now().Unix()/30))
	return fmt.Sprintf("%06d", number), nil
}


func (this *GoogleAuth) GetQrcode(user, secret string) string {
	return fmt.Sprintf("otpauth://totp/%s?secret=%s", user, secret)
}


func (this *GoogleAuth) GetQrcodeUrl(user, secret string) string {
	qrcode := this.GetQrcode(user, secret)
	return fmt.Sprintf("http://www.google.com/chart?chs=200x200&chld=M%%7C0&cht=qr&chl=%s", qrcode)
}


func (this *GoogleAuth) VerifyCode(secret, code string) (bool, error) {
	_code, err := this.GetCode(secret)
	if err != nil {
		return false, err
	}
	return _code == code, nil
}
