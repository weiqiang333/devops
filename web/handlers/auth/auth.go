package auth

import (
	"fmt"
	"github.com/weiqiang333/devops/internal/authentication"
	"net/http"
	"strings"

	"github.com/gin-gonic/contrib/sessions"
	"github.com/gin-gonic/gin"
)

const (
	userkey = "user"
)


// AuthRequired is a simple middleware to check the session
func AuthRequired(c *gin.Context) {
	user := Me(c)
	if user == nil {
		fmt.Printf("Unauthorized")
		// Abort the request with the appropriate error code
		c.AbortWithStatus(http.StatusUnauthorized)
		c.HTML(http.StatusUnauthorized, "login.tmpl", gin.H{"Authentication": "Unauthorized"})
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
		c.HTML(http.StatusBadRequest, "login.tmpl", gin.H{"Authentication": "Parameters can't be empty"})
		return
	}

	// Check for username and password match, LDAP Authentication
	if ! authentication.LdapAuthentication(username, password) {
		c.HTML(http.StatusUnauthorized, "login.tmpl", gin.H{"Authentication": "Authentication failed"})
		return
	}

	// Save the username in the session
	session.Set(userkey, username) // In real world usage you'd set this to the users ID
	if err := session.Save(); err != nil {
		c.HTML(http.StatusInternalServerError, "login.tmpl", gin.H{"Authentication": "Failed to save session"})
		return
	}
	c.Redirect(http.StatusMovedPermanently, "/")
}


func Logout(c *gin.Context) {
	session := sessions.Default(c)
	user := Me(c)
	if user == nil {
		c.HTML(http.StatusBadRequest, "login.tmpl", gin.H{"Authentication": "Invalid session token"})
		return
	}
	fmt.Println(user)
	session.Delete(userkey)
	fmt.Println("1", session.Get(userkey))
	if err := session.Save(); err != nil {
		c.HTML(http.StatusInternalServerError, "login.tmpl", gin.H{"Authentication": "Failed to save session"})
		return
	}
	fmt.Println("2", session.Get(userkey))
	fmt.Printf("333")
	c.Redirect(http.StatusFound, "/")
}


func Me(c *gin.Context) interface{} {
	session := sessions.Default(c)
	username := session.Get(userkey)
	return username
}


func Status(c *gin.Context) {
	username := Me(c)
	c.JSON(http.StatusOK, gin.H{"status": "You are logged in", "username": username})
}