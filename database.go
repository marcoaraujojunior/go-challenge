package main

import(
	"os"
	"log"
	"github.com/go-sql-driver/mysql"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql"
	"github.com/marcoaraujojunior/go-challenge/model"
)

var (
	db *gorm.DB
)

func connect() {
	dbConf := mysql.Config{
		User:      os.Getenv("DB_USER"),
		Passwd:    os.Getenv("DB_PASS"),
		Net:       "tcp",
		Addr:      os.Getenv("DB_HOST"),
		DBName:    os.Getenv("DB_NAME"),
		ParseTime: true,
	}

	dsn := dbConf.FormatDSN()

	dbConn(dsn)

	db.AutoMigrate(&model.Invoice{})
}

func dbConn(dsn string) {
	log.Printf("Connecting: %s", dsn)
	conn, err := gorm.Open("mysql", dsn)

	if err != nil {
		log.Fatal("[DB err ]: %s", err)
	}

	db = conn
}
