package service

import (
	"fmt"
	"net/http"
	"strings"
	"log"

	"github.com/gin-gonic/gin"
	"github.com/weiqiang333/devops/internal/database"
	"github.com/weiqiang333/devops/internal/sshclient"
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
}


//ListService 列出服务
func ListService(c *gin.Context) {
	//cmdList := `systemctl list-unit-files | grep -E '\.service\s+(generated|enabled)' | awk -F'.service' '{print $1}' |
	//	grep -vE 'acpid|atd|auditd|autovt@|chronyd|crond|cloud-config|cloud-final|cloud-init|dmraid-activation|getty@|hibinit-agent|irqbalance|lvm2-monitor|libstoragemgmt|mdmonitor|microcode|postfix|rngd|rpcbind|rsyslog|sysstat|systemd-readahead-collect|systemd-readahead-drop|systemd-readahead-replay|update-motd|amazon-ssm-agent'`
	if c.Request.Method == http.MethodPost {
		server := c.Query("server")
		service := c.Query("service")
		action, _ := c.GetPostForm("action")
		response := serviceCmd(server, service, action)
		servers := searchService(server)
		c.HTML(http.StatusOK, "service.tmpl", gin.H{
			"server": servers,
			"response": response,
		})
		return
	}
	servers := searchService("")
	c.HTML(http.StatusOK, "service.tmpl", gin.H{
		"server": servers,
	})
	return
}


func serviceCmd(server, service, action string) string {
	cmd := fmt.Sprintf("sudo systemctl %s %s", action, service)
	response, err := sshclient.SSHCline("/data/wei_loacl/sshkey/apps_rsa","apps", server,"222", cmd)
	if err != nil {
		return err.Error()
	}
	return response
}


func searchService(server string) []serverList {
	//sql := fmt.Sprintf("SELECT server, status FROM server_list;")
	servers := []serverList{}
	sql := fmt.Sprintf(`SELECT server_list.server, server_list.name, server_list.app, server_list.pillar, server_list.status,
			string_agg(service.service, ',') AS service
		FROM server_list
		JOIN service ON server_list.server = service.server
		GROUP BY server_list.server, server_list.name, server_list.app, server_list.pillar, server_list.status;`)
	db := database.Db()
	row, err := db.Query(sql)
	defer row.Close()
	if err != nil {
		log.Printf("search service error: %v", err)
		return servers
	}

	for row.Next() {
		var s serverList
		if err := row.Scan(&s.Server, &s.Name, &s.App, &s.Pillar, &s.Status, &s.Service); err != nil {
			log.Printf("db rows scan fail: %v", err)
		}
		servers = append(servers, s)
	}
	for i, s := range servers {
		servers[i].Id = i
		servers[i].Services = strings.Split(s.Service, ",")
		if s.Server == server {
			servers[i].Action = "true"
		}
	}
	return servers
}