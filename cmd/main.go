package main

import (
	"github.com/joho/godotenv"
	"github.com/spf13/viper"
	"humoBooking/internal/database"
	"humoBooking/internal/handlers"
	"humoBooking/internal/services"
	"log"
)

// main TODO add viper and godotenv
func main() {
	if err := InitConfig(); err != nil {
		log.Fatal("InitConfig(): error loading config.yaml", err)
	}
	if err := godotenv.Load(); err != nil {
		log.Fatal("error loading .env file", err)
	}
	conn := database.NewConnectPostgres()
	repository := database.NewDatabase(conn)
	service := services.NewService(repository)
	handler := handlers.NewHandler(service)
	log.Println(handler.Init())
}

func InitConfig() error {
	viper.AddConfigPath("internal/configs")
	viper.SetConfigName("config") // расширение не надо указывать
	return viper.ReadInConfig()
}
