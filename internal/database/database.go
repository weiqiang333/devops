package database

import (
	"fmt"
	"log"

	"database/sql"
	_ "github.com/lib/pq"
	"github.com/spf13/viper"
)


// Db open connect
func Db() (db *sql.DB) {
	driverName := viper.GetString("database.driver_name")
	host := viper.GetString("database.host")
	port := viper.GetString("database.port")
	dbname := viper.GetString("database.dbname")
	user := viper.GetString("database.user")
	password := viper.GetString("database.password")
	db, err := sql.Open(driverName, fmt.Sprintf("host=%s port=%s dbname=%s user=%s password=%s", host, port, dbname, user, password))
	if err != nil {
		log.Fatalf("unable to open database: %v", err)
	}
	return db
}