package authentication

import (
	"crypto/tls"
	"fmt"
	"log"

	. "github.com/go-ldap/ldap/v3"
	"github.com/spf13/viper"
)


func ldapDial() (*Conn, error) {
	address := viper.GetString("authentication.ldap.address")
	port := viper.GetString("authentication.ldap.port")
	bindusername := viper.GetString("authentication.ldap.bindusername")
	bindpassword := viper.GetString("authentication.ldap.bindpassword")

	l, err := Dial("tcp", fmt.Sprintf("%s:%s", address, port))
	if err != nil {
		log.Printf("Authentication for LDAP Dial: %v", err)
		return nil, err
	}
	// Reconnect with TLS
	err = l.StartTLS(&tls.Config{InsecureSkipVerify: true})
	if err != nil {
		return nil, fmt.Errorf("Authentication for LDAP StartTLS: %v", err)
	}
	err = l.Bind(bindusername, bindpassword)
	if err != nil {
		log.Printf("LADP Bind Failed: %v", err)
	}
	return l, err
}


// Ldap 验证用户
func LdapAuthentication(username, password string) bool {
	l, _ := ldapDial()
	if l == nil {
		return false
	}
	defer l.Close()

	//Bind as the user to verify their password
	err := l.Bind(username, password)
	if err != nil {
		log.Printf("Authentication user for LDAP Bind: %v", err)
		return false
	}
	log.Printf("Authentication user for LDAP Success: %s", username)
	return true
}


// LdapGetDN 获取 group/user DN 信息
func LdapGetDN(class, name string) (string, error) {
	baseDN := viper.GetString("authentication.ldap.basedn")
	l, err := ldapDial()
	defer l.Close()
	if err != nil {
		return "", fmt.Errorf("ldapDial Failed")
	}

	searchRequest := NewSearchRequest(
		baseDN, // The base dn to search
		ScopeWholeSubtree, NeverDerefAliases, 0, 0, false,
		fmt.Sprintf("(&(objectClass=%s)(CN=%s))", class, name),
		[]string{"dn"},
		nil,
	)

	sr, err := l.Search(searchRequest)
	if err != nil {
		return "", fmt.Errorf("Search Group for LDAP: %s", err.Error())
	}

	if len(sr.Entries) != 1 {
		return "", fmt.Errorf("Groups does not exist or too many entries returned: %s", name)
	}
	return sr.Entries[0].DN, nil
}


// LdapGroupUser 组内用户获取
func LdapGroupUser(groupDN string) ([]string, error) {
	baseDN := viper.GetString("authentication.ldap.basedn")
	groupUser := []string{}
	l, err := ldapDial()
	defer l.Close()
	if err != nil {
		return []string{}, fmt.Errorf("ldapDial Failed")
	}

	searchRequest := NewSearchRequest(
		baseDN, // The base dn to search
		ScopeWholeSubtree, NeverDerefAliases, 0, 0, false,
		fmt.Sprintf("(&(objectClass=user)(memberof:=%s))", groupDN),
		[]string{"cn"},
		nil,
	)

	sr, err := l.Search(searchRequest)
	if err != nil {
		return []string{}, fmt.Errorf("Search Group for LDAP: %s", err.Error())
	}

	for _, entry := range sr.Entries {
		groupUser = append(groupUser, entry.GetAttributeValue("cn"))
	}
	if len(groupUser) == 0 {
		return []string{}, fmt.Errorf("LdapGroupUser Search for nil")
	}
	return groupUser, nil
}
