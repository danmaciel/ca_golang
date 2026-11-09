package main

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	graphqlhandler "github.com/99designs/gqlgen/graphql/handler"
	"github.com/99designs/gqlgen/graphql/playground"
	"github.com/danmaciel/ca_golang/configs"
	"github.com/danmaciel/ca_golang/internal/infra/graph"
	"github.com/danmaciel/ca_golang/internal/infra/grpc/pb"
	"github.com/danmaciel/ca_golang/internal/infra/grpc/service"
	"github.com/danmaciel/ca_golang/internal/infra/web/webserver"
	"github.com/danmaciel/ca_golang/pkg/events"
	_ "github.com/go-sql-driver/mysql"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

func main() {
	if err := run(); err != nil {
		log.Fatal(err)
	}
}

func run() error {
	cfg, err := configs.LoadConfig(".")
	if err != nil {
		return fmt.Errorf("load configuration: %w", err)
	}

	db, err := sql.Open(cfg.DBDriver, fmt.Sprintf("%s:%s@tcp(%s:%s)/%s",
		cfg.DBUser, cfg.DBPassword, cfg.DBHost, cfg.DBPort, cfg.DBName))
	if err != nil {
		// Driver errors may echo the DSN (which contains the password).
		return fmt.Errorf("open database: %T", err)
	}
	defer db.Close()

	dispatcher := events.NewEventDispatcher()
	rest := webserver.NewWebServer(cfg.WebServerPort)
	webOrders := NewWebOrderHandler(db, dispatcher)
	rest.AddHandler("/order", webOrders.Create)
	rest.AddHandler("/list-orders", webOrders.GetAll)

	createOrder := NewCreateOrderUseCase(db, dispatcher)
	listOrders := NewListOrderUseCase(db, dispatcher)
	graphql := graphqlhandler.NewDefaultServer(graph.NewExecutableSchema(graph.Config{Resolvers: &graph.Resolver{
		CreateOrderUseCase: *createOrder,
	}}))
	graphqlMux := http.NewServeMux()
	graphqlMux.Handle("/", playground.Handler("GraphQL playground", "/query"))
	graphqlMux.Handle("/query", graphql)

	grpcServer := grpc.NewServer()
	pb.RegisterOrderServiceServer(grpcServer, service.NewOrderService(*createOrder, *listOrders))
	reflection.Register(grpcServer)

	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()
	log.Printf("Starting REST %s, GraphQL :%s and gRPC :%s", cfg.WebServerPort, cfg.GraphQLServerPort, cfg.GRPCServerPort)
	return listenAndServe(ctx, cfg.WebServerPort, ":"+cfg.GraphQLServerPort,
		":"+cfg.GRPCServerPort, rest.Handler(), graphqlMux, grpcServer)
}
