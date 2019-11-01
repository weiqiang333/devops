package crontab

import (
	"fmt"
	"log"
	"strings"

	"github.com/spf13/viper"
	"github.com/weiqiang333/devops/internal/database"
	"github.com/weiqiang333/devops/internal/model"
	"github.com/weiqiang333/devops/internal/sshclient"
)


func getService(server string) (string, error) {
	privateKey := viper.GetString("sshcline.private_key")
	username := viper.GetString("sshcline.username")
	port := viper.GetString("sshcline.port")
	cmd := `systemctl list-unit-files --state=enabled --type=service | grep service | awk -F '.service' '{print $1}' | 
		grep -vE '^(acpid|auditd|autovt@|atd|update-motd|sysstat|rpcbind|rngd|postfix|microcode|lvm2-monitor|libstoragemgmt|irqbalance|hibinit-agent|getty@|dmraid-activation|chronyd|crond|ntpd|rsyslog|mdmonitor|cloud-|amazon-|systemd-)'`
	response, err := sshclient.SSHCline(privateKey, username, server, port, cmd)
	if err != nil {
		return response, err
	}
	return response, nil
}


func getServers() []string {
	var server string
	servers := []string{}
	sql := fmt.Sprintf(`SELECT server FROM server_list;`)
	db := database.Db()
	row, err := db.Query(sql)
	defer row.Close()
	defer db.Close()
	if err != nil {
		log.Printf("getServers search server error: %v", err)
		return servers
	}
	for row.Next() {
		if err := row.Scan(&server); err != nil {
			log.Printf("db rows scan fail: %v", err)
		}
		servers = append(servers, server)
	}
	return servers
}


func insertService(server string, service string)  {
	sql := fmt.Sprintf(`
		INSERT INTO service (server, service)
		VALUES ('%s', '%s');`, server, service)
	db := database.Db()
	row, err := db.Query(sql)
	defer row.Close()
	defer db.Close()
	if err != nil {
		log.Printf("insert server_list error: %s %s - %v", server, service, err)
	}
}


func deleteService(server string, service string)  {
	sql := fmt.Sprintf(`
		DELETE FROM service WHERE server = '%s' AND service = '%s';`, server, service)
	db := database.Db()
	row, err := db.Query(sql)
	defer row.Close()
	defer db.Close()
	if err != nil {
		log.Printf("delete server_list error: %s %s - %v", server, service, err)
	}
}


func selectService(server string) []string {
	services := []string{}
	sql := fmt.Sprintf(`
		SELECT service FROM service WHERE server = '%s'
	`, server)
	db := database.Db()
	row, err := db.Query(sql)
	defer row.Close()
	defer db.Close()
	if err != nil {
		log.Printf("selectService error: %s - %v", server, err)
	}
	for row.Next() {
		var s model.TableService
		if err := row.Scan(&s.Service); err != nil {
			log.Printf("selectService db rows scan fail: %v", err)
		}
		services = append(services, s.Service)
	}
	return services
}


//UpdateService update database tables service
func UpdateService()  {
	log.Println("cron update service start")
	servers := getServers()
	for _, server := range servers {
		response, err := getService(server)
		if err != nil {
			continue
		}
		servicesCurrent := strings.Split(response, "\n")
		servicesTables := selectService(server)
		for _, service := range servicesCurrent {
			if service != "" && ! checkValue(service, servicesTables) {
				insertService(server, service)
			}
		}
		for _, service := range servicesTables {
			if ! checkValue(service, servicesCurrent) {
				deleteService(server, service)
			}
		}
	}
	log.Println("cron update service done")
}


func checkValue(service string, services []string) bool {
	if len(services) == 0 {
		return false
	}
	if service == services[0] {
		return true
	}
	return checkValue(service, services[1:])
}