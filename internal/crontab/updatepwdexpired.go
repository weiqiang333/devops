package crontab

import (
	"fmt"
	"log"
	"strconv"
	"time"

	"github.com/weiqiang333/devops/internal/authentication"
	"github.com/weiqiang333/devops/internal/database"
)


//通过 pwdLastSet 获取180天前修改，在未来一周过期用户
func updatePwdExpired() {
	ts := time.Now().AddDate(0, -6, 7).UnixNano() / 1e6
	tsNT := utcNT(ts)
	res, err := authentication.LdapGetPwdLastSet(tsNT)
	if err != nil {
		log.Printf("updatePwdExpired Failed: %v", err)
		return
	}

	err = truncatePwdExpired()
	if err != nil {
		return
	}
	for name, tsSet := range res {
		tsSet, err := strconv.ParseInt(tsSet, 10, 64)
		tsExp :=  tsSet + ntExpiredTimes(180)
		if err != nil {
			log.Printf("cront updatePwdExpired: %v", err)
			return
		}
		insertPwdExpired(name, formatTimestamp(ntUTC(tsSet) / 1000), formatTimestamp(ntUTC(tsExp) / 1000))
	}
	log.Println("updatePwdExpired Success")
}


func formatTimestamp(ts int64) string {
	return time.Unix(ts, 0).Format("2006-01-02 15:04:05-07")
}


//insertPwdExpired tables for ldap_pwd_expired
func insertPwdExpired(name, pwdLastSet, pwdExpired string)  {
	sql := fmt.Sprintf(`
		INSERT INTO ldap_pwd_expired (name, pwd_last_set, pwd_expired)
		VALUES ('%s', '%s', '%s');`, name, pwdLastSet, pwdExpired)
	db := database.Db()
	row, err := db.Query(sql)
	defer row.Close()
	defer db.Close()
	if err != nil {
		log.Printf("insert ldap_pwd_expired error: %s - %v", sql, err)
		return
	}
}


//truncatePwdExpired tables for ldap_pwd_expired
func truncatePwdExpired() error {
	sql := "TRUNCATE ldap_pwd_expired RESTART IDENTITY RESTRICT;"
	db := database.Db()
	row, err := db.Query(sql)
	defer row.Close()
	defer db.Close()
	if err != nil {
		log.Printf("truncate ldap_pwd_expired error: %s - %v", sql, err)
		return err
	}
	log.Println("truncate ldap_pwd_expired Success")
	return nil
}

//ntExpiredTimes NT time number of 100-nanosecond
func ntExpiredTimes(day int64) int64 {
	return (day * 24 * 60 * 60 * 1000 * 10000)
}


//utcNT return 100-nanosecond
func utcNT(ts int64) int64 {
	return (ts - 57599875) * 10000 + 116445312000000000
}


//ntUTC return Millisecond
func ntUTC(ts int64) int64 {
	return (ts - 116445312000000000) / 10000 + 57599875
}
