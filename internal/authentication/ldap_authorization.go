/*
	Ldap 用户安全组权限确认
*/
package authentication

import (
	"log"
)


// authorization LDAP
// 授权确认，用户存在于安全组中
func SecurityGroupAuthorization(username string) bool {
	groupDN, err := LdapGetDN("group","sre")
	if err != nil {
		log.Println(err.Error())
		return false
	}

	users, err := LdapGroupUser(groupDN)
	if err != nil {
		log.Println(err.Error())
		return false
	}
	return recursionUsers(username, users)
}


// recursionUsers 确认用户授权在用户组中
func recursionUsers(username string, users []string) bool {
	if len(users) == 0 {
		return false
	}
	if username == users[0] {
		return true
	} else  {
		return recursionUsers(username, users[1:])
	}
}