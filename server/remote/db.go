package remote

import (
	"fmt"
	"log"
	"os"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

const (
	CODE_SUCCESS          = 0
	DB_ERROR              = 1
	DB_EMPTY_RESULT_ERROR = 2
	DB_NOT_UPDATED        = 3
)

var db *gorm.DB

func InitDB() error {
	//username := os.Getenv("MYSQL_USERNAME")
	//password := os.Getenv("MYSQL_PASSWORD")
	host := os.Getenv("MYSQL_HOST")
	port := os.Getenv("MYSQL_PORT")
	dbName := os.Getenv("MYSQL_DB")
	//dsn := fmt.Sprintf("%v:%v@tcp(%v:%v)/%v?charset=utf8mb4&parseTime=True&loc=Local", username, password, host, port, dbName)
	dsn := fmt.Sprintf("%v:%v@tcp(%v:%v)/%v?charset=utf8mb4&parseTime=True&loc=Local", "root", "123456789", host, port, dbName)
	fmt.Println("dsn: ", dsn)
	var err error
	db, err = gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Printf("[ERROR] init db error: %v", err.Error())
	}
	return err
}

func GetDB() *gorm.DB {
	return db
}