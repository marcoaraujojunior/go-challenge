package database

import(
	"os"
	"log"
	"github.com/go-sql-driver/mysql"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql"
)

var (
	Db *gorm.DB
)

func Connect() {
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

}

func dbConn(dsn string) {
	log.Printf("Connecting: %s", dsn)
	conn, err := gorm.Open("mysql", dsn)

	if err != nil {
		log.Fatal("[DB err ]: %s", err)
	}

	Db = conn
}

