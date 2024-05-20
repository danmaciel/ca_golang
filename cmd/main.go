package main

import (
	"database/sql"
	"fmt"
	"net"
	"net/http"

	"github.com/99designs/gqlgen/graphql/playground"
	"github.com/danmaciel/ca_golang/configs"
	"github.com/danmaciel/ca_golang/internal/infra/graph"
	"github.com/danmaciel/ca_golang/internal/infra/grpc/pb"
	"github.com/danmaciel/ca_golang/internal/infra/grpc/service"
	"github.com/danmaciel/ca_golang/internal/infra/web/webserver"
	"github.com/danmaciel/ca_golang/pkg/events"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"

	graphql_handler "github.com/99designs/gqlgen/graphql/handler"
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

	webserver := webserver.NewWebServer(configs.WebServerPort)
	webOrderHandler := NewWebOrderHandler(db, eventDispatcher)
	webserver.AddHandler("/order", webOrderHandler.Create)
	webserver.AddHandler("/list-orders", webOrderHandler.GetAll)
	fmt.Println("Starting web server on port", configs.WebServerPort)
	go webserver.Start()

	createOrderUseCase := NewCreateOrderUseCase(db, eventDispatcher)
	listOrderUseCase := NewListOrderUseCase(db, eventDispatcher)

	// graphql
	srv := graphql_handler.NewDefaultServer(graph.NewExecutableSchema(graph.Config{Resolvers: &graph.Resolver{
		CreateOrderUseCase: *createOrderUseCase,
	}}))

	http.Handle("/", playground.Handler("GraphQL playground", "/query"))
	http.Handle("/query", srv)

	fmt.Printf("connect to http://localhost%s/ for GraphQL playground\n", configs.GraphQLServerPort)
	go http.ListenAndServe(":"+configs.GraphQLServerPort, nil)

	grpcServer := grpc.NewServer()
	createOrderService := service.NewOrderService(*createOrderUseCase, *listOrderUseCase)

	pb.RegisterOrderServiceServer(grpcServer, createOrderService)

	reflection.Register(grpcServer)

	fmt.Println("Starting gRPC server on port", configs.GRPCServerPort)
	lis, err := net.Listen("tcp", fmt.Sprintf(":%s", configs.GRPCServerPort))
	if err != nil {
		panic(err)
	}

	grpcServer.Serve(lis)

	print("Rodou")
}
