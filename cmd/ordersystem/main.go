package main

import (
	"database/sql"
	"fmt"
	"net"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	graphql_handler "github.com/99designs/gqlgen/graphql/handler"
	"github.com/99designs/gqlgen/graphql/playground"
	"github.com/lucas4ndrade/FullcycleCleanArch/configs"
	"github.com/lucas4ndrade/FullcycleCleanArch/internal/event/handler"
	"github.com/lucas4ndrade/FullcycleCleanArch/internal/infra/graph"
	"github.com/lucas4ndrade/FullcycleCleanArch/internal/infra/grpc/pb"
	"github.com/lucas4ndrade/FullcycleCleanArch/internal/infra/grpc/service"
	"github.com/lucas4ndrade/FullcycleCleanArch/internal/infra/web"
	"github.com/lucas4ndrade/FullcycleCleanArch/internal/infra/web/webserver"
	"github.com/lucas4ndrade/FullcycleCleanArch/internal/usecase"
	"github.com/lucas4ndrade/FullcycleCleanArch/pkg/events"
	"github.com/streadway/amqp"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"

	// mysql
	_ "github.com/go-sql-driver/mysql"
)

func main() {
	conf, err := configs.LoadConfig(".")
	if err != nil {
		panic(err)
	}

	db := connectDB(conf)
	defer db.Close()

	eventDispatcher := getEventDispatcher(conf)

	createOrderUseCase := NewCreateOrderUseCase(db, eventDispatcher)
	listOrderUseCase := NewListOrderUseCase(db)

	go startWebServer(conf, createOrderUseCase, listOrderUseCase)
	go startGRPCServer(conf, createOrderUseCase, listOrderUseCase)
	go startGrapQLServer(conf, createOrderUseCase, listOrderUseCase)

	stopChan := make(chan os.Signal, 1)
	signal.Notify(stopChan, syscall.SIGTERM, syscall.SIGINT)
	<-stopChan
}

func connectDB(conf *configs.Config) *sql.DB {
	db, err := sql.Open(conf.DBDriver, fmt.Sprintf("%s:%s@tcp(%s:%s)/%s", conf.DBUser, conf.DBPassword, conf.DBHost, conf.DBPort, conf.DBName))
	if err != nil {
		panic(err)
	}

	return db
}

func getEventDispatcher(conf *configs.Config) *events.EventDispatcher {
	rabbitMQChannel := getRabbitMQChannel(conf.AmqpURL)

	eventDispatcher := events.NewEventDispatcher()
	err := eventDispatcher.Register("OrderCreated", &handler.OrderCreatedHandler{
		RabbitMQChannel: rabbitMQChannel,
	})
	if err != nil {
		panic(err)
	}

	return eventDispatcher
}

func getRabbitMQChannel(url string) *amqp.Channel {
	conn, err := amqp.Dial(url)
	if err != nil {
		panic(err)
	}
	ch, err := conn.Channel()
	if err != nil {
		panic(err)
	}
	return ch
}

func startWebServer(
	conf *configs.Config,
	createOrderUC *usecase.CreateOrderUseCase,
	listOrderUC *usecase.ListOrderUseCase,
) {
	webserver := webserver.NewWebServer(":" + conf.WebServerPort)
	webserver.AddHandler(http.MethodPost, "/order", web.CreateOrderHandler(createOrderUC))
	webserver.AddHandler(http.MethodGet, "/order", web.ListOrderHandler(listOrderUC))

	fmt.Println("Starting web server on port", conf.WebServerPort)
	webserver.Start()
}

func startGRPCServer(
	conf *configs.Config,
	createOrderUC *usecase.CreateOrderUseCase,
	listOrderUC *usecase.ListOrderUseCase,
) {
	grpcServer := grpc.NewServer()
	orderService := service.NewOrderService(*createOrderUC, *listOrderUC)
	pb.RegisterOrderServiceServer(grpcServer, orderService)
	reflection.Register(grpcServer)

	lis, err := net.Listen("tcp", fmt.Sprintf(":%s", conf.GRPCServerPort))
	if err != nil {
		panic(err)
	}

	fmt.Println("Starting gRPC server started on port", conf.GRPCServerPort)
	err = grpcServer.Serve(lis)
	if err != nil {
		panic(err)
	}
}

func startGrapQLServer(
	conf *configs.Config,
	createOrderUC *usecase.CreateOrderUseCase,
	listOrderUC *usecase.ListOrderUseCase,
) {
	srv := graphql_handler.NewDefaultServer(graph.NewExecutableSchema(graph.Config{Resolvers: &graph.Resolver{
		CreateOrderUseCase: *createOrderUC,
		ListOrderUseCase:   *listOrderUC,
	}}))
	http.Handle("/", playground.Handler("GraphQL playground", "/query"))
	http.Handle("/query", srv)

	fmt.Println("Starting GraphQL server on port", conf.GraphQLServerPort)
	err := http.ListenAndServe(":"+conf.GraphQLServerPort, nil)
	if err != nil {
		panic(err)
	}
}
