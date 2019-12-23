package ldapadmin

import (
	"net/http"
	"log"

	"github.com/gin-gonic/gin"

	"github.com/weiqiang333/devops/web/handlers/auth"
	"github.com/weiqiang333/devops/internal/database"
	"github.com/weiqiang333/devops/internal/model"
)


//LadpAdmin ldap admin html
func LadpAdmin(c *gin.Context)  {
	username := auth.Me(c)
	pwdExpired := searchPwdExpired()
	c.HTML(http.StatusBadRequest, "ldapadmin.tmpl", gin.H{
		"ldapAdmin": "active",
		"user": username,
		"pwdExpired": pwdExpired,
	})
}

func searchPwdExpired() []model.TableLdapPwdExpired {
	pwdExpireds := []model.TableLdapPwdExpired{}
	sql := "SELECT * FROM ldap_pwd_expired;"
	db := database.Db()
	row, err := db.Query(sql)
	defer row.Close()
	defer db.Close()
	if err != nil {
		log.Printf("search ldap_pwd_expired error: %v", err)
		return nil
	}

	for row.Next() {
		var pe model.TableLdapPwdExpired
		if err := row.Scan(&pe.Id, &pe.Name, &pe.PwdLastSet, &pe.PwdExpired); err != nil {
			log.Printf("db rows scan fail for ldap_pwd_expired: %v", err)
		}
		pwdExpireds = append(pwdExpireds, pe)
	}
	return pwdExpireds
}
