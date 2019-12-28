package auth

import (
	"log"
	"net/http"
	"strings"
	"fmt"

	"github.com/gin-gonic/contrib/sessions"
	"github.com/gin-gonic/gin"

	"github.com/weiqiang333/devops/internal/authentication"
)

const (
	userkey = "user"
)


// AuthRequired is a simple middleware to check the session
func AuthRequired(c *gin.Context) {
	user := Me(c)
	if user == nil {
		log.Printf("Unauthorized")
		// Abort the request with the appropriate error code
		c.AbortWithStatus(http.StatusUnauthorized)
		c.HTML(http.StatusUnauthorized, "user/login.tmpl", gin.H{"Authentication": "Unauthorized"})
		return
	}
	// Continue down the chain to handler etc
	c.Next()
}


// login is a handler that parses a form and checks for specific data
func Login(c *gin.Context) {
	session := sessions.Default(c)
	username := c.PostForm("username")
	password := c.PostForm("password")
	// Validate form input
	if strings.Trim(username, " ") == "" || strings.Trim(password, " ") == "" {
		c.HTML(http.StatusBadRequest, "user/login.tmpl", gin.H{"Authentication": "Parameters can't be empty"})
		return
	}

	// Check for username and password match, LDAP Authentication
	if ! authentication.LdapAuthentication(username, password) {
		c.HTML(http.StatusUnauthorized, "user/login.tmpl", gin.H{"Authentication": "Authentication failed"})
		return
	}

	// Save the username in the session
	session.Set(userkey, username) // In real world usage you'd set this to the users ID
	if err := session.Save(); err != nil {
		c.HTML(http.StatusInternalServerError, "user/login.tmpl", gin.H{"Authentication": "Failed to save session"})
		return
	}
	c.Redirect(http.StatusFound, "/")
}


//Logout
func Logout(c *gin.Context) {
	session := sessions.Default(c)
	user := Me(c)
	if user == nil {
		c.HTML(http.StatusBadRequest, "user/login.tmpl", gin.H{"Authentication": "Invalid session token"})
		return
	}
	session.Delete(userkey)
	if err := session.Save(); err != nil {
		c.HTML(http.StatusInternalServerError, "user/login.tmpl", gin.H{"Authentication": "Failed to save session"})
		return
	}
	log.Printf("Logout user for %s", user)
	c.Redirect(http.StatusFound, "/")
}


//Me authentication username
func Me(c *gin.Context) interface{} {
	session := sessions.Default(c)
	username := session.Get(userkey)
	return username
}


//Users Info
func Users(c *gin.Context) {
	username := Me(c)
	qrCodeUrl, err := SearchQRcodeUrl(fmt.Sprint(username))
	if err != nil {
		log.Printf("Users: %v", err)
		c.HTML(http.StatusInternalServerError, "user/users.tmpl", gin.H{
			"users": "active",
			"user": username,
			"qrCodeUrl": "查询异常",
		})
		return
	}
	c.HTML(http.StatusOK, "user/users.tmpl", gin.H{
		"users": "active",
		"user": username,
		"qrCodeUrl": qrCodeUrl,
	})
}
