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
	Status bool `json:"status" form:"status"`
	services
}

type services struct {
	Service string `json:"service" form:"service"`
	Services []string
}


//ListService 列出服务
func ListService(c *gin.Context) {
	//cmdList := `systemctl list-unit-files | grep -E '\.service\s+(generated|enabled)' | awk -F'.service' '{print $1}' |
	//	grep -vE 'acpid|atd|auditd|autovt@|chronyd|crond|cloud-config|cloud-final|cloud-init|dmraid-activation|getty@|hibinit-agent|irqbalance|lvm2-monitor|libstoragemgmt|mdmonitor|microcode|postfix|rngd|rpcbind|rsyslog|sysstat|systemd-readahead-collect|systemd-readahead-drop|systemd-readahead-replay|update-motd|amazon-ssm-agent'`
	if c.Request.Method == http.MethodPost {
		server := c.Query("server")
		service := c.Query("service")
		action := c.Query("action")
		c.JSON(http.StatusOK, gin.H{
			"server": server,
			"service": service,
			"action": action,
		})
		return
	}
	servers := searchService()
	c.HTML(http.StatusOK, "service.tmpl", gin.H{
		"server": servers,
	})
}


func cmd(cmd string) []string {
	response := sshclient.SSHCline("/data/wei_loacl/sshkey/apps_rsa","apps","10.0.1.55","222", cmd)
	serviceList := strings.Split(response, "\n")
	return serviceList
}


func searchService() []serverList {
	//sql := fmt.Sprintf("SELECT server, status FROM server_list;")
	sql := fmt.Sprintf(`SELECT server_list.server, server_list.status, string_agg(service.service, ',') AS service
		FROM server_list
		JOIN service ON server_list.server = service.server
		GROUP BY server_list.server, server_list.status;`)
	db := database.Db()
	row, err := db.Query(sql)
	if err != nil {
		log.Fatalf("search service error: %v", err)
	}
	defer row.Close()
	servers := []serverList{}
	for row.Next() {
		var s serverList
		if err := row.Scan(&s.Server, &s.Status, &s.Service); err != nil {
			log.Fatalf("db rows scan fail: %v", err)
		}
		servers = append(servers, s)
	}
	fmt.Println(servers)
	for i, a := range servers {
		servers[i].Id = i
		servers[i].Services = strings.Split(a.Service, ",")
	}
	return servers
}