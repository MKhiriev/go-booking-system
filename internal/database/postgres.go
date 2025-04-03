package database

import (
	"fmt"
	"github.com/spf13/viper"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"log"
	"os"
)

func NewConnectPostgres() *gorm.DB {
	host := viper.GetString("db.host")
	port := viper.GetUint16("db.port")
	username := viper.GetString("db.username")
	password := os.Getenv("DB_PASSWORD")
	dbName := viper.GetString("db.DBName")

	dbParams := fmt.Sprintf("host=%s password=%s user=%s dbname=%s port=%d sslmode=disable TimeZone=Asia/Dushanbe",
		host, password, username, dbName, port)

	postgresDialector := postgres.Open(dbParams)
	connection, err := gorm.Open(postgresDialector, &gorm.Config{
		Logger: logger.Default.LogMode(logger.Silent)})

	if err != nil {
		log.Println("NewConnectPostgres(): error occurred")
		return nil
	}

	log.Println("NewConnectPostgres(): successful connection to db")
	return connection
}
