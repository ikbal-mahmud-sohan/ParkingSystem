package main

import (
	"context"
	"embed"
	"fmt"
	"log"
	"net"
	"net/http"
	"os/signal"
	"syscall"

	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"

	plCore "parking/api/core/parkinglot"
	psCore "parking/api/core/parkingslot"
	prsCore "parking/api/core/parkingstatement"
	vCore "parking/api/core/vehicle"
	plSvc "parking/api/services/parkinglot"
	psSvc "parking/api/services/parkingslot"
	prsSvc "parking/api/services/parkingstatement"
	pvSvc "parking/api/services/vehicle"
	"parking/api/storage/postgres"
	"parking/utility"
	env "parking/utility/envvariable"
	"parking/utility/logging"
)

var apiMethodsContent embed.FS

var (
	module         = "api"
	devEnvironment = "development"
)

type RegisterGatewayFunc func(ctx context.Context, mux *runtime.ServeMux, gwAddr string, opts []grpc.DialOption) error

// RegisterGateway ...
func (f RegisterGatewayFunc) RegisterGateway(ctx context.Context, mux *runtime.ServeMux, gwAddr string, opts []grpc.DialOption) error {
	return f(ctx, mux, gwAddr, opts)
}

type Service interface {
	RegisterSvc(*grpc.Server) error
	RegisterGateway(context.Context, *runtime.ServeMux, string, []grpc.DialOption) error
}

type Svcs struct {
	external []Service
	internal []Service
}

func main() {
	cfg, err := utility.NewConfig("env/config")
	if err != nil {
		log.Fatal(err)
	}

	logger := logging.NewLogger(cfg).WithFields(logrus.Fields{
		"module":      module,
		"environment": cfg.GetString("runtime.environment"),
	})

	ctx, cancel := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer cancel()

	dbString := utility.NewDBString(cfg)
	db, err := postgres.NewStorage(dbString, logger)
	if err != nil {
		logger.WithError(err).Errorln("unable to connect DB")
		return
	}
	defer db.Close()

	if err := db.RunMigration(cfg.GetString(env.DatabaseMigrationDir)); err != nil {
		logger.WithError(err).Errorln("unable to run DB migrations")
		return
	}

	svc := setupGRPCServices(cfg, logger, db)
	go setupGRPCInternalServer(ctx, cfg, logger, db, svc)
	go setupGRPCExternalServer(ctx, cfg, logger, svc)

	<-ctx.Done()

	logger.Info("All services and servers have been shut down gracefully.")
}

func setupGRPCInternalServer(ctx context.Context, config *viper.Viper, logger *logrus.Entry, store *postgres.Storage, s *Svcs) {
	grpcServer := grpc.NewServer()
	for _, svc := range s.internal {
		svc.RegisterSvc(grpcServer)
	}
	addr := fmt.Sprintf(":%d", config.GetInt(env.ServerPort))
	lis, err := net.Listen("tcp", addr)
	if err != nil {
		logger.WithError(err).Errorf("unable to listen port on: %d\n", config.GetInt(env.ServerPort))
		return
	}

	go func() {
		logger.WithField("url", lis.Addr()).Info("starting gRPC server")
		if err := grpcServer.Serve(lis); err != nil {
			logger.Fatalf("Failed to serve gRPC server: %v\n", err)
		}
	}()

	<-ctx.Done()
	logger.Infoln("Stopped gRPC server.")
}

func setupGRPCExternalServer(ctx context.Context, config *viper.Viper, logger *logrus.Entry, s *Svcs) {
	opts := []grpc.DialOption{grpc.WithBlock(), grpc.WithTransportCredentials(insecure.NewCredentials())}
	addr := fmt.Sprintf(":%d", config.GetInt(env.ServerGatewayPort))
	grpcAddr := fmt.Sprintf(":%d", config.GetInt(env.ServerPort))

	mux := runtime.NewServeMux()
	for _, r := range s.external {
		if err := r.RegisterGateway(ctx, mux, grpcAddr, opts); err != nil {
			logger.Fatalf("unable to register gateway: %v\n", err)
			return
		}
	}

	gwServer := &http.Server{
		Addr:    addr,
		Handler: mux,
	}
	go func() {
		logger.WithField("url", addr).Println("starting gRPC gateway server")
		if err := gwServer.ListenAndServe(); err != nil {
			logger.Fatalf("Failed to serve gRPC gateway server: %v\n", err)
		}
	}()

	<-ctx.Done()
	logger.Infoln("Stopped gRPC gateway server.")
}

func setupGRPCServices(config *viper.Viper, logger *logrus.Entry, store *postgres.Storage) *Svcs {
	plSvc := plSvc.New(plCore.New(store, logger))
	psSvc := psSvc.New(psCore.New(store, logger))
	pvSvc := pvSvc.New(vCore.New(store, logger))
	prsSvc := prsSvc.New(prsCore.New(store, logger))
	return &Svcs{
		external: []Service{plSvc, psSvc, pvSvc, prsSvc},
		internal: []Service{plSvc, psSvc, pvSvc, prsSvc},
	}
}
