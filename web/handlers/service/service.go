package service

import (
	"fmt"
	"log"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/spf13/viper"

	"github.com/weiqiang333/devops/internal/database"
	"github.com/weiqiang333/devops/internal/sshclient"
	"github.com/weiqiang333/devops/internal/authentication"
	"github.com/weiqiang333/devops/web/handlers/auth"
)

type serverList struct {
	Id int
	Server string `json:"server" form:"server"`
	Name string `json:"name" form:"name"`
	App string `json:"app" form:"app"`
	Pillar string `json:"pillar" form:"pillar"`
	Status bool `json:"status" form:"status"`
	services
	Action string
}

type services struct {
	Service string `json:"service" form:"service"`
	Services []string
	ServiceStatus string `json:"service_status" form:"service_status"`
	ServiceStatuss []string
}


//ListService 列出服务
func ListService(c *gin.Context) {
	username := auth.Me(c)
	if c.Request.Method == http.MethodPost {
		server := c.Query("server")
		service := c.Query("service")
		action, _ := c.GetPostForm("action")
		servers := searchService(server)

		if ! authorization(fmt.Sprint(username)) {
			log.Printf("Info ListService action %s Currently not authorized for %s", action, username)
			c.HTML(http.StatusLocked, "service.tmpl", gin.H{
				"service": "active",
				"server": servers,
				"response": "Currently not authorized, please contact SRE",
				"user": username,
			})
			return
		}

		actionResponse := serviceCmd(server, service, action)
		c.HTML(http.StatusOK, "service.tmpl", gin.H{
			"service": "active",
			"server": servers,
			"response": actionResponse,
			"user": username,
		})
		return
	}
	servers := searchService("")
	c.HTML(http.StatusOK, "service.tmpl", gin.H{
		"service": "active",
		"server": servers,
		"user": username,
	})
	return
}


func serviceCmd(server, service, action string) string {
	privateKey := viper.GetString("sshcline.private_key")
	username := viper.GetString("sshcline.username")
	port := viper.GetString("sshcline.port")
	cmd := fmt.Sprintf("sudo systemctl %s %s", action, service)
	response, err := sshclient.SSHCline(privateKey, username, server, port, cmd)
	if err != nil {
		return err.Error()
	}
	return response
}


func searchService(server string) []serverList {
	servers := []serverList{}
	sql := fmt.Sprintf(`SELECT server_list.server, server_list.name, server_list.app, server_list.pillar, server_list.status,
		string_agg(service.service, ',') AS service, string_agg(service.status, ',') AS service_status
		FROM server_list
		JOIN service ON server_list.server = service.server
		GROUP BY server_list.server, server_list.name, server_list.app, server_list.pillar, server_list.status;`)
	db := database.Db()
	row, err := db.Query(sql)
	defer row.Close()
	defer db.Close()
	if err != nil {
		log.Printf("search service error: %v", err)
		return servers
	}

	for row.Next() {
		var s serverList
		if err := row.Scan(&s.Server, &s.Name, &s.App, &s.Pillar, &s.Status, &s.Service, &s.ServiceStatus); err != nil {
			log.Printf("db rows scan fail: %v", err)
		}
		servers = append(servers, s)
	}
	for i, s := range servers {
		servers[i].Id = i
		servers[i].Services = strings.Split(s.Service, ",")
		servers[i].ServiceStatuss = strings.Split(s.ServiceStatus, ",")
		if s.Server == server {
			servers[i].Action = "true"
		}
	}
	return servers
}


// authorization LDAP
func authorization(username string) bool {
	groupDN, err := authentication.LdapGetDN("group","sre")
	if err != nil {
		log.Println(err.Error())
		return false
	}

	users, err := authentication.LdapGroupUser(groupDN)
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