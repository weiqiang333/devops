/*
	将二次验证码发送至邮箱，进行验证
	二次验证码复用 google-authenticator, 有效期 5m
*/

package auth

import (
	"net/http"
	"log"

	"github.com/gin-gonic/gin"

	"github.com/weiqiang333/devops/internal/authentication"
	"github.com/weiqiang333/devops/internal/mailcli"
)


//PostMailVerificationCode 发送邮箱验证码
func PostMailVerificationCode(c *gin.Context)  {
	name := c.PostForm("username")
	// 需要优化：确认用户真实存在
	search, err := SearchQRcodeSecret(name)
	if err != nil {
		search, err = createSecret(name)
		if err != nil {
			log.Printf("GetMailVerificationCode createSecret error: %v", err)
			c.JSON(http.StatusNotImplemented, gin.H{"response": err.Error()})
			return
		}
	}
	code, err := authentication.NewGoogleAuth().GetMailCode(search)
	if err != nil {
		log.Printf("GetMailVerificationCode GetMailCode error: %v", err)
		c.JSON(http.StatusNotImplemented, gin.H{"response": err.Error()})
		return
	}
	//发送
	if err := mailcli.PostMailCode(name, code); err != nil {
		c.JSON(http.StatusNotImplemented, gin.H{"response": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"response": "Send successfully"})
	return
}


//GetForgetPwd 忘记密码页面
func GetForgetPwd(c *gin.Context)  {
	c.HTML(http.StatusOK, "ldapadmin/forgetPwd.tmpl", gin.H{"ldapAdmin": "active",})
}


//ModifyUserPwd 修改用户密码 handler
func ModifyUserPwd(c *gin.Context) {
	user, userOk := c.GetPostForm("username")
	password, pwdOk := c.GetPostForm("password")
	qrcode, qrOK := c.GetPostForm("qrcode")
	if ! userOk || ! pwdOk || ! qrOK {
		log.Printf("ModifyUserPwd fail, Please check the parameters: %s", user)
		c.JSON(http.StatusBadRequest, gin.H{
			"response": "ModifyUserPwd fail, Please check the parameters",
		})
		return
	}

	secret, err := SearchQRcodeSecret(user)
	if err != nil {
		log.Printf("ModifyUserPwd fail, Please confirm to enable secondary verification: %s; %v", user, err)
		c.JSON(http.StatusUnauthorized, gin.H{
			"response": "ModifyUserPwd fail, Please confirm to enable secondary verification",
		})
		return
	}

	ok, _ := authentication.NewGoogleAuth().VerifyCode(secret, qrcode)
	mOk, _ := authentication.NewGoogleAuth().MailVerifyCode(secret, qrcode)
	if ok || mOk {
		//修改密码
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
	} else {
		log.Printf("ModifyUserPwd fail, Secondary verification failed, please verify again: %s; %v", user, err)
		c.JSON(http.StatusUnauthorized, gin.H{
			"response": "ModifyUserPwd fail, Secondary verification failed, please verify again",
		})
		return
	}
}
