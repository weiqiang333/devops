/*
	用于邮箱验证码的验证，复用 google-authenticationa
*/

package authentication

import (
	"fmt"
	"strings"
	"time"
)


//GetMailCode 获取用于 mail 的 5m 有效验证码
func (this *GoogleAuth) GetMailCode(secret string) (string, error) {
	secretUpper := strings.ToUpper(secret)
	secretKey, err := this.base32decode(secretUpper)
	if err != nil {
		return "", err
	}
	number := this.oneTimePassword(secretKey, this.toBytes(time.Now().Unix()/300))
	return fmt.Sprintf("%06d", number), nil
}


//MailVerifyCode 邮箱码验证
func (this *GoogleAuth) MailVerifyCode(secret, code string) (bool, error) {
	_code, err := this.GetMailCode(secret)
	if err != nil {
		return false, err
	}
	return _code == code, nil
}