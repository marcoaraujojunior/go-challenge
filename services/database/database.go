package database

import(
	"os"
	"log"
	"github.com/go-sql-driver/mysql"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql"
)

var (
	db *gorm.DB
)

func GetDb() *gorm.DB {
	if (db == nil) {
		Connect()
	}
	return db
}

func Connect() {
	config := GetConfig()
	dsn := config.FormatDSN()
	conn, err := gorm.Open("mysql", dsn)
	if err != nil {
		log.Fatal("[DB err ]: %s", err)
	}
	db = conn
}

func GetConfig() mysql.Config {
	return mysql.Config{
		User:      os.Getenv("DB_USER"),
		Passwd:    os.Getenv("DB_PASS"),
		Net:       "tcp",
		Addr:      os.Getenv("DB_HOST"),
		DBName:    os.Getenv("DB_NAME"),
		ParseTime: true,
	}
}

