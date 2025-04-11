package main

import (
	"github.com/joho/godotenv"
	"github.com/spf13/viper"
	"humoBooking/internal/database"
	"humoBooking/internal/handlers"
	"humoBooking/internal/server"
	"humoBooking/internal/services"
	"log"
)

// TODO добавить JWT токен!!! jwt.CustomClaims
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
	myServer := new(server.Server)
	err := myServer.ServerRun(handler.Init(), viper.GetString("server.port"))
	check(err)
}

func InitConfig() error {
	viper.AddConfigPath("internal/configs")
	viper.SetConfigName("config")
	return viper.ReadInConfig()
}

func check(err error) {
	if err != nil {
		log.Println(err)
	}
}
