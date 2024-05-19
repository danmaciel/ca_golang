package main

import (
	"database/sql"
	"fmt"

	"github.com/danmaciel/clean_arch_golang/configs"
	"github.com/danmaciel/clean_arch_golang/internal/infra/web/webserver"
	"github.com/danmaciel/clean_arch_golang/pkg/events"

	// mysql
	_ "github.com/go-sql-driver/mysql"
)

func main() {
	configs, err := configs.LoadConfig(".")
	if err != nil {
		panic(err)
	}

	db, err := sql.Open(configs.DBDriver, fmt.Sprintf("%s:%s@tcp(%s:%s)/%s", configs.DBUser, configs.DBPassword, configs.DBHost, configs.DBPort, configs.DBName))
	if err != nil {
		panic(err)
	}
	defer db.Close()
	eventDispatcher := events.NewEventDispatcher()

	//createOrderUseCase := NewCreateOrderUseCase(db, eventDispatcher)

	webserver := webserver.NewWebServer(configs.WebServerPort)
	webOrderHandler := NewWebOrderHandler(db, eventDispatcher)
	webserver.AddHandler("/order", webOrderHandler.Create)
	webserver.AddHandler("/order-list", webOrderHandler.GetAll)
	fmt.Println("Starting web server on port", configs.WebServerPort)
	//go webserver.Start()
	webserver.Start()

	print("Rodou")
}
