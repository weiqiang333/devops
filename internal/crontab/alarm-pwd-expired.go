package crontab

import (
	"fmt"
	"github.com/spf13/viper"
	"github.com/weiqiang333/devops/internal/dingtalk"
	"log"
	"time"

	"github.com/weiqiang333/devops/internal/database"
	"github.com/weiqiang333/devops/internal/model"
)

func alarmPwdExpired() {
	alarmPwdExpiredToken := viper.GetString("cron.alarm_pwd_expired_token")
	pwdExpireds, err := selectPwdExpired()
	if err != nil {
		log.Printf("cron alarmPwdExpired fail: %v", err)
		return
	}

	message := createMessage(pwdExpireds)
	err = dingtalk.Dingtalk(alarmPwdExpiredToken, message, "false")
	if err != nil {
		log.Println(err)
	}
}

func createMessage(pwdExpireds []model.TableLdapPwdExpired) (message string) {
	message = "LDAP 用户密码将过期用户。\n"
	for i, v := range pwdExpireds {
		message += fmt.Sprintf("\t\t%v\t\t%s\t\t%s\n", i, v.Name, v.PwdExpired.Format("2006-01-02 15:04:05"))
	}
	message += fmt.Sprintf("LDAP 用户密码将影响您正常使用：Wifi、VPN、Jenkins、Consul、Zabbix 等内部工具.\n" +
		"重置/忘记密码，可以通过 https://devops-infra.growingio.com/ldapAdmin/ 进行操作")
	return message
}

func selectPwdExpired() ([]model.TableLdapPwdExpired, error) {
	ofTime := time.Now().UTC().AddDate(0, 0, -15).Format("2006-01-02 15:04:05-07")
	sql := fmt.Sprintf(`SELECT name, pwd_expired FROM ldap_pwd_expired WHERE pwd_expired >= '%s' ORDER BY pwd_expired;`, ofTime)
	pwdExpireds := []model.TableLdapPwdExpired{}

	db := database.Db()
	defer db.Close()
	row, err := db.Query(sql)
	defer row.Close()
	if err != nil {
		log.Printf("select ldap_pwd_expired error: %v", err)
		return pwdExpireds, fmt.Errorf("select fail")
	}

	for row.Next() {
		pwdExpired := model.TableLdapPwdExpired{}
		err = row.Scan(&pwdExpired.Name, &pwdExpired.PwdExpired)
		if err != nil {
			log.Printf("selectPwdExpired Scan fail: %v", err)
			continue
		}
		pwdExpireds = append(pwdExpireds, pwdExpired)
	}
	return pwdExpireds, nil
}
