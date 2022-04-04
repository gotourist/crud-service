package main

import (
	"fmt"
	"github.com/iman_task/crud-service/domain/service"
	broker "github.com/iman_task/crud-service/pkg/messagebroker"
	"log"
	"net"

	configPkg "github.com/iman_task/crud-service/config"
	"github.com/iman_task/crud-service/events"
	"github.com/iman_task/crud-service/events/handlers"
	pb "github.com/iman_task/crud-service/genproto/post"
	loggerPkg "github.com/iman_task/crud-service/pkg/logger"
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

func main() {
	// =========================================================================
	// Configurations loading...
	config := configPkg.Load()

	// =========================================================================
	// Logger
	logger := loggerPkg.New("debug", "crud-service")
	defer func() {
		err := loggerPkg.Cleanup(logger)
		if err != nil {
			logger.Fatal("failed cleaning up logs: %v", loggerPkg.Error(err))
		}
	}()

	// =========================================================================
	// Postgres
	logger.Info("Postgresql configs",
		loggerPkg.String("host", config.PostgresHost),
		loggerPkg.Int("port", config.PostgresPort),
		loggerPkg.String("database", config.PostgresDatabase),
	)

	psqlString := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=%s",
		config.PostgresHost,
		config.PostgresPort,
		config.PostgresUser,
		config.PostgresPassword,
		config.PostgresDatabase,
		config.PostgresSSL)

	// db initialization
	connDb, err := sqlx.Connect("postgres", psqlString)
	if err != nil {
		logger.Error("postgres connect error", loggerPkg.Error(err))
		return
	}

	// =========================================================================
	// Kafka

	// Publishers
	publishersMap := make(map[string]broker.Producer)

	postChangeTopicPublisher := events.NewKafkaProducer(&config, logger, handlers.PostChangeTopic)
	defer func() {
		err := postChangeTopicPublisher.Stop()
		if err != nil {
			logger.Fatal("Error while publishing: %v", loggerPkg.Error(err))
		}
	}()

	publishersMap[handlers.PostChangeTopic] = postChangeTopicPublisher

	// Listeners
	postAddTopicListener := events.NewKafkaConsumer(connDb, &config, logger, handlers.PostAddTopic)
	go postAddTopicListener.Start()

	// =========================================================================
	// gRPC server
	postService := service.NewPostService(connDb, logger, config, publishersMap)

	listen, err := net.Listen("tcp", config.RPCPort)
	if err != nil {
		logger.Fatal("error while listening: %v", loggerPkg.Error(err))
	}
	s := grpc.NewServer()

	pb.RegisterPostServiceServer(s, postService)
	reflection.Register(s)

	logger.Info("main: server running", loggerPkg.String("port", config.RPCPort))

	if err := s.Serve(listen); err != nil {
		log.Fatalf("Error while listening: %v", loggerPkg.Error(err))
	}
}
