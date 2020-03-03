package ldapadmin

import (
	"fmt"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/weiqiang333/devops/web/handlers/auth"
	"github.com/weiqiang333/devops/internal/authentication"
)


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

	ok, err := auth.VerifyCode(user, qrcode)
	if err != nil || ! ok {
		log.Printf("ModifyUserPwd fail: %s; VerifyCode %v; error message: %v", user, ok, err)
		c.JSON(http.StatusUnauthorized, gin.H{
			"response": fmt.Sprintf("VerifyCode %v; error message: %v", ok, err),
		})
		return
	}

	err = authentication.LdapModifyPwd(user, password)
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
	c.HTML(http.StatusOK, "ldapadmin/modifyUserPwd.tmpl", gin.H{
		"ldapAdmin": "active",
		"user": user,
	})
}
