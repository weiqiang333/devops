package ldapadmin

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/weiqiang333/devops/web/handlers/auth"
	"log"
	"net/http"

	"github.com/weiqiang333/devops/internal/authentication"
)


//modifyPwd 修改用户密码
func modifyPwd(name, password string) error {
	userDN, err := authentication.LdapGetDN("user", name)
	if err != nil {
		log.Printf("ModifyPwd LdapGetDN fail for %s: %v", name, err)
		return err
	}

	err = authentication.LdapModifyPwd(userDN, password)
	if err != nil {
		log.Printf("ModifyPwd LdapModifyPwd fail for %s: %v", name, err)
		return err
	}
	return nil
}


//ModifyUserPwd 修改用户密码 handler
func ModifyUserPwd(c *gin.Context) {
	username, userOk := c.GetPostForm("username")
	password, pwdOk := c.GetPostForm("password")
	qrcode, qrOK := c.GetPostForm("qrcode")
	user := fmt.Sprint(auth.Me(c))
	if ! userOk || ! pwdOk || ! qrOK || username != user {
		log.Printf("ModifyUserPwd fail, Please check the parameters: %s", username)
		c.JSON(http.StatusBadRequest, gin.H{
			"response": "ModifyUserPwd fail, Please check the parameters",
		})
		return
	}

	secret, err := auth.SearchQRcodeSecret(user)
	if err != nil {
		log.Printf("ModifyUserPwd fail, Please confirm to enable secondary verification: %s; %v", user, err)
		c.JSON(http.StatusUnauthorized, gin.H{
			"response": "ModifyUserPwd fail, Please confirm to enable secondary verification",
		})
		return
	}
	ok, err := authentication.NewGoogleAuth().VerifyCode(secret, qrcode)
	if err != nil || ! ok {
		log.Printf("ModifyUserPwd fail, Secondary verification failed, please verify again: %s; %v", user, err)
		c.JSON(http.StatusUnauthorized, gin.H{
			"response": "ModifyUserPwd fail, Secondary verification failed, please verify again",
		})
		return
	}

	err = modifyPwd(user, password)
	if err != nil {
		log.Printf("ModifyUserPwd fail, Please check if the password policy is met. %s; %v", user, err)
		c.JSON(http.StatusRequestedRangeNotSatisfiable, gin.H{
			"response": "ModifyUserPwd fail, Please check if the password policy is met.",
		})
		return
	}

	log.Printf("ModifyUserPwd Success: %s", user)
	c.JSON(http.StatusOK, gin.H{
		"response": "ModifyUserPwd Success",
	})
	return
}


//GetModifyUserPwd
func GetModifyPwdPage(c *gin.Context)  {
	user := auth.Me(c)
	c.HTML(http.StatusOK, "modifyUserPwd.tmpl", gin.H{
		"ldapAdmin": "active",
		"user": user,
	})
}