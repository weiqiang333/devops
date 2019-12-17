package authentication

import (
	"crypto/tls"
	"fmt"
	"log"

	"golang.org/x/text/encoding/unicode"
	"gopkg.in/ldap.v2"
	. "github.com/go-ldap/ldap/v3"
	"github.com/spf13/viper"
)


//ldapAdminDial 拥有管理权限，并且 SSL 通信
func ldapAdminDial() (*ldap.Conn, error) {
	address := viper.GetString("authentication.ldap.address")
	port := viper.GetString("authentication.ldap.port")
	bindusername := viper.GetString("authentication.ldap.bindusername")
	bindpassword := viper.GetString("authentication.ldap.bindpassword")

	l, err := ldap.Dial("tcp", fmt.Sprintf("%s:%s", address, port))
	if err != nil {
		return nil, fmt.Errorf("Authentication for LDAP Dial: %v", err)
	}
	// Reconnect with TLS
	err = l.StartTLS(&tls.Config{InsecureSkipVerify: true})
	if err != nil {
		return nil, fmt.Errorf("Authentication for LDAP StartTLS: %v", err)
	}
	//bind
	err = l.Bind(bindusername, bindpassword)
	if err != nil {
		log.Printf("LADP Bind Failed: %v", err)
	}
	return l, err
}


//ldapModifyPwd 修改用户密码
func LdapModifyPwd(userDN, password string) error {
	l, err := ldapAdminDial()
	defer l.Close()
	if err != nil {
		return fmt.Errorf("LdapModifyPwd Failed: %v", err)
	}

	utf16 := unicode.UTF16(unicode.LittleEndian, unicode.IgnoreBOM)
	pwdEncoded, _ := utf16.NewEncoder().String(fmt.Sprintf("\"%s\"", password))
	passReq := &ldap.ModifyRequest{
		DN: userDN, // DN for the user we're resetting
		ReplaceAttributes: []ldap.PartialAttribute{
			{"unicodePwd", []string{pwdEncoded}},
		},
	}
	err = l.Modify(passReq)
	if err != nil {
		return fmt.Errorf("Password could not be changed: %v", err)
	}
	return nil
}


//LdapGetPwdLastSet 获取最后修改密码用户超于 lastSetTime 之前
func LdapGetPwdLastSet(lastSetTime int64) (map[string]string, error) {
	baseDN := viper.GetString("authentication.ldap.basedn")
	l, err := ldapDial()
	defer l.Close()
	if err != nil {
		return nil, fmt.Errorf("LdapGetPwdLastSet Failed: %v", err)
	}

	searchRequest := NewSearchRequest(
		baseDN, // The base dn to search
		ScopeWholeSubtree, NeverDerefAliases, 0, 0, false,
		fmt.Sprintf("(&(objectCategory=person)(objectClass=user)(&(pwdLastSet>=1)(pwdLastSet<=%v)))", lastSetTime),
		[]string{},                    // A list attributes to retrieve
		nil,
	)

	sr, err := l.Search(searchRequest)
	if err != nil {
		return nil, fmt.Errorf("LdapGetPwdLastSet Search Failed: baseDN(%s) %v", baseDN, err)
	}

	responses := map[string]string{}
	for _, entry := range sr.Entries {
		responses[entry.GetAttributeValue("cn")] = entry.GetAttributeValue("pwdLastSet")
	}
	return responses, nil
}