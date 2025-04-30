package main

import (
	"github.com/joho/godotenv"
	"github.com/spf13/viper"
	_ "go-booking-system/cmd/docs"
	"go-booking-system/internal/database"
	"go-booking-system/internal/handlers"
	"go-booking-system/internal/server"
	"go-booking-system/internal/services"
	"log"
)

// @title go-booking-system
// @version 1.0
// @description Booking system written in go. Implemented Authorization, Handlers, Services, Repositories, Models for booking.

// @contact.name Rasul Khiriev
// @contact.url https://github.com/MKhiriev/go-booking-system
// @contact.email khiriev.rasul@inbox.ru
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
