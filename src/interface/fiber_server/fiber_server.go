package fiber_server

import (
	"os"
	"os/signal"
	"rental-vehicle-system/src/use_case"
	"sync"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	log "github.com/sirupsen/logrus"
)

type FiberServer struct {
	useCase *use_case.UseCase
	server  *fiber.App
	config  *ServerConfig
}

type ServerConfig struct {
	AppVersion    string
	RequestLog    bool
	ListenAddress string
}

func New(uc *use_case.UseCase, sc *ServerConfig) *FiberServer {
	server := fiber.New(fiber.Config{
		CaseSensitive:         false,
		StrictRouting:         false,
		DisableStartupMessage: true,
		ReadTimeout:           30 * time.Second,
	})

	server.Use(cors.New(cors.Config{
		AllowOrigins: "*",
		AllowHeaders: "*",
		AllowMethods: "GET, POST, PUT, PATCH, DELETE, OPTIONS",
	}))

	f := &FiberServer{
		uc,
		server,
		sc,
	}

	f.addRouteVehicle(server)
	f.addRouteBooking(server)

	return f
}

func (f FiberServer) Start(wg *sync.WaitGroup) {
	wg.Add(2)

	exit := make(chan os.Signal, 1)
	signal.Notify(exit, os.Interrupt)

	go func() {
		defer wg.Done()
		<-exit
		log.Info("Shutting down server...")

		err := f.server.Shutdown()
		if err != nil {
			log.Info("Server shutdown with error", err)
		} else {
			log.Info("Server gracefully shutdown")
		}
	}()

	go func() {
		defer wg.Done()
		log.Info("Server is starting...")
		err := f.server.Listen(f.config.ListenAddress)
		if err != nil {
			log.Info("Server error", err)
		}
		log.Info("Server has been shutdown")
	}()
}
