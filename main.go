package main

import (
	"context"
	"fmt"
	"sync"
	"time"

	"rental-vehicle-system/src/interface/fiber_server"
	"rental-vehicle-system/src/repository/booking_repository"
	"rental-vehicle-system/src/repository/vehicle_repository"
	"rental-vehicle-system/src/use_case"

	"github.com/caarlos0/env/v6"
	"github.com/joho/godotenv"
	log "github.com/sirupsen/logrus"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type config struct {
	AppVersion  string `env:"APP_VERSION" envDefault:"v0.0.0"`
	Port        uint   `env:"PORT" envDefault:"8080"`
	MongoDbUri  string `env:"MONGO_DB_URI" envDefault:"mongodb://root:password@localhost:27017"`
	MongoDBName struct {
		RentalDB string `env:"MONGO_DB_NAME" envDefault:"rental-vehicle"`
	}
}

func main() {
	cfg := initEnvironment()
	vehicleRepo, bookingRepo := initRepo(cfg)

	useCase := use_case.New(vehicleRepo, bookingRepo)

	initInterfaces(cfg, useCase)
}

func initEnvironment() config {
	err := godotenv.Load()
	if err != nil {
		log.Printf("Error loading .env file: %s\n", err)
	}

	var cfg config
	err = env.Parse(&cfg)
	if err != nil {
		log.Fatalf("Error parse env: %s\n", err)
	}

	return cfg
}

func initRepo(cfg config) (use_case.VehicleRepository, use_case.BookingRepository) {
	log.Info("Start to connect DB")
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	client, err := mongo.Connect(ctx, options.Client().ApplyURI(cfg.MongoDbUri))
	if err != nil {
		log.Fatalf("Connection DB error")
	}

	vehicleRepo := vehicle_repository.NewMongoDb(client.Database(cfg.MongoDBName.RentalDB))
	bookingRepo := booking_repository.NewMongoDb(client.Database(cfg.MongoDBName.RentalDB))

	log.Info("Connection DB success")

	return vehicleRepo, bookingRepo
}

func initInterfaces(cfg config, useCase *use_case.UseCase) {
	wg := new(sync.WaitGroup)

	serv := fiber_server.New(useCase, &fiber_server.ServerConfig{
		AppVersion:    cfg.AppVersion,
		ListenAddress: fmt.Sprintf(":%d", cfg.Port),
		RequestLog:    true,
	})
	log.Info("Fiber server initialized")

	serv.Start(wg)
	log.Info("Fiber server started on port")

	wg.Wait()
	log.Info("Application stopped")

}
