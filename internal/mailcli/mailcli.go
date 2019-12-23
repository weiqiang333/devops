package mailcli

import (
	"fmt"
	"strconv"

	"github.com/spf13/viper"
	"gopkg.in/mail.v2"
)


//PostMailCode 发送邮箱验证码
func PostMailCode(user, code string) error {
	host := viper.GetString("mail.host")
	port, _ := strconv.Atoi(viper.GetString("mail.port"))
	username := viper.GetString("mail.username")
	password := viper.GetString("mail.password")
	domain := viper.GetString("mail.domain")

	d := mail.NewDialer(host, port, username, password)
	d.StartTLSPolicy = mail.MandatoryStartTLS

	m := mail.NewMessage()
	m.SetHeader("From", username)
	m.SetHeader("To", fmt.Sprintf("%s@%s", user, domain))
	m.SetHeader("Subject", "欢迎使用 DevOps 系统!")
	m.SetBody("text/html", fmt.Sprintf(
		"<h3>很高兴遇见你</h3><hr/><p>您正在使用 LDAP 密码找回功能，您的验证码是：<b>%s</b></p><br/><br/><p>如有问题,请联系 SRE 团队</p>",
		code),
	)


	// Send the email to Bob, Cora and Dan.
	if err := d.DialAndSend(m); err != nil {
		return err
	}
	return nil
}