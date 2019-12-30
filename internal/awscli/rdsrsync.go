package awscli

import (
	"fmt"
	"log"
	"regexp"

	"github.com/spf13/viper"

	"github.com/weiqiang333/devops/internal/database"
	"github.com/weiqiang333/devops/internal/model"
)

type SearchTableOrderLogs struct {
	model.TableRdsRsyncOrder
	model.TableRdsRsyncOrderLogs
}

type SearchTableOrder struct {
	model.TableRdsRsyncOrder
	Status interface{} `json:"status" form:"status"`
}


//CreateWorkorder 创建同步工单
func CreateWorkorder(databaseName, username string) error {
	match := viper.GetString("aws.rds.match")
	if match == "" {
		match = "ci-"
	}
	matched, _ := regexp.MatchString(fmt.Sprintf("^(%s).*", match), databaseName)
	if ! matched {
		return fmt.Errorf("CreateWorkorder fail: 你所输入的数据库目前并不支持")
	}
	err := insertWorkorder(databaseName, username)
	if err != nil {
		return err
	}
	return nil
}


func insertWorkorder(databaseName, username string) error {
	sql := fmt.Sprintf(`
		INSERT INTO rds_rsync_workorder (database, username)
		VALUES ('%s', '%s');`, databaseName, username)
	db := database.Db()
	defer db.Close()
	row, err := db.Query(sql)
	defer row.Close()
	if err != nil {
		log.Printf("insert rds_rsync_workorder error: %s %s - %v", databaseName, username, err)
		return fmt.Errorf("insertWorkorder fail: %v", err)
	}
	return nil
}


//SearchWorkorder
func SearchWorkorder(id int) ([]model.TableRdsRsyncWorkorder, error) {
	selectWorkorder := []model.TableRdsRsyncWorkorder{}
	sql := fmt.Sprintf("SELECT * FROM rds_rsync_workorder ORDER BY id DESC;")
	if id != 0 {
		sql = fmt.Sprintf("SELECT * FROM rds_rsync_workorder WHERE id=%v ORDER BY id DESC;", id)
	}
	db := database.Db()
	defer db.Close()
	row, err := db.Query(sql)
	defer row.Close()
	if err != nil {
		log.Printf("select rds_rsync_workorder error: %v", err)
		return selectWorkorder, fmt.Errorf("select rds_rsync_workorder fail: %v", err)
	}
	for row.Next() {
		var wo model.TableRdsRsyncWorkorder
		if err := row.Scan(&wo.Id, &wo.Database, &wo.Username, &wo.CreatedAt, &wo.PassAt, &wo.OrderStatus); err != nil {
			log.Printf("db rows scan fail for rds_rsync_workorder: %v", err)
		}
		selectWorkorder = append(selectWorkorder, wo)
	}
	return selectWorkorder, nil
}


//SearchOrder
func SearchOrder(workorderId int, databaseName string) ([]SearchTableOrder, error) {
	selectOrder := []SearchTableOrder{}
	//sql := fmt.Sprintf("SELECT * FROM rds_rsync_order WHERE database='default' ORDER BY priority ASC;")
	sql := fmt.Sprintf(`SELECT o.id, o.priority, o.authorized_user, 
		(SELECT status FROM rds_rsync_order_logs WHERE workorderid=%v AND orderid=o.id) AS status
		FROM rds_rsync_order o
		WHERE database='default'
		ORDER BY o.priority ASC;`, workorderId)
	db := database.Db()
	defer db.Close()
	row, err := db.Query(sql)
	defer row.Close()
	if err != nil {
		log.Printf("select rds_rsync_order error: %v", err)
		return selectOrder, fmt.Errorf("select rds_rsync_workorder fail: %v", err)
	}
	for row.Next() {
		var ro SearchTableOrder
		if err := row.Scan(&ro.Id, &ro.Priority, &ro.AuthorizedUser, &ro.Status); err != nil {
			log.Printf("db rows scan fail for rds_rsync_workorder: %v", err)
		}
		selectOrder = append(selectOrder, ro)
	}
	return selectOrder, nil
}


//SearchOrderLogs
func SearchOrderLogs(workorderId int) ([]SearchTableOrderLogs, error) {
	selectOrderLogs := []SearchTableOrderLogs{}
	//sql := fmt.Sprintf("SELECT * FROM rds_rsync_order_logs WHERE workorderid=%v ORDER BY created_at DESC;", workorderId)
	sql := fmt.Sprintf(`SELECT o.priority, o.authorized_user, ol.created_at, ol.status
		FROM rds_rsync_order o, rds_rsync_order_logs ol
		WHERE ol.workorderid=%v and o.id=ol.orderid
		ORDER BY ol.created_at DESC;`, workorderId)
	db := database.Db()
	defer db.Close()
	row, err := db.Query(sql)
	defer row.Close()
	if err != nil {
		log.Printf("select rds_rsync_order error: %v", err)
		return selectOrderLogs, fmt.Errorf("select rds_rsync_workorder fail: %v", err)
	}
	for row.Next() {
		var rol SearchTableOrderLogs
		if err := row.Scan(&rol.Priority, &rol.AuthorizedUser, &rol.CreatedAt, &rol.Status); err != nil {
			log.Printf("db rows scan fail for rds_rsync_order_logs: %v", err)
		}
		selectOrderLogs = append(selectOrderLogs, rol)
	}
	return selectOrderLogs, nil
}


//InsertOrderLog
func InsertOrderLog(workorderId int, orderId int, status bool) error {
	sql := fmt.Sprintf(`
		INSERT INTO rds_rsync_order_logs (workorderid, orderid, status)
		VALUES (%v, %v, %v);`, workorderId, orderId, status)
	db := database.Db()
	defer db.Close()
	row, err := db.Query(sql)
	defer row.Close()
	if err != nil {
		log.Printf("insert rds_rsync_order_logs error: %v %v - %v", workorderId, orderId, err)
		return fmt.Errorf("insertWorkorder fail: %v", err)
	}
	return nil
}


//IfRsyncStatus
func IfRsyncStatus(workorderId int) {
	sql := fmt.Sprintf(`
		SELECT all_count
		FROM (
			SELECT COUNT(o.id)
			FROM rds_rsync_order_logs ol
			JOIN rds_rsync_order o
			ON o.id = ol.orderid
			WHERE ol.workorderid=2 AND ol.status=TRUE
			AND o.id IN (
				SELECT id
				FROM rds_rsync_order
				WHERE database='default' AND priority IN (
					SELECT MAX(priority)
					FROM rds_rsync_order
					WHERE database='default')
			)) ok_count, (
				SELECT COUNT(id)
				FROM rds_rsync_order
				WHERE database='default' AND priority IN (
					SELECT MAX(priority)
					FROM rds_rsync_order
					WHERE database='default')
			) all_count
		WHERE ok_count=all_count;`)
	updateSql := fmt.Sprintf(`
		UPDATE rds_rsync_workorder
		SET order_status = 'rsync', pass_at = now() at time zone 'utc'
		WHERE id=%v;`, workorderId)
	db := database.Db()
	defer db.Close()
	row, err := db.Query(sql)
	defer row.Close()
	if err != nil {
		log.Printf("select rds_rsync_order sqlAll error: %v - %v", workorderId, err)
		return
	}
	if row.Next() {
		_, err = db.Query(updateSql)
		if err != nil {
			log.Printf("update rds_rsync_workorder order_status error: %v - %v", workorderId, err)
			return
		}
	}
	return
}
