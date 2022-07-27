package main

import (
	"context"
	"fmt"
	services "loggercise/cmd/server/services"
	settings "loggercise/configs"
	loggerciseProtoPackage "loggercise/gen/go/service"
	loggerciseHandlerPackage "loggercise/internal/loggerciseHandler"
	"loggercise/internal/store"
	"net"
	"net/http"

	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"github.com/sethvargo/go-envconfig"
	"github.com/sirupsen/logrus"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

const (
	host     = "0.0.0.0"
	port     = "80"
	grpcPort = "50051"
)

func main() {

	ctx := context.Background()
	logger := logrus.New()
	logger.Info("Starting loggercise server")
	// env config
	var config settings.Settings
	if err := envconfig.Process(ctx, &config); err != nil {
		logger.Fatal(err)
	}
	addr := fmt.Sprintf("%s:%s", host, grpcPort)
	httpAddr := fmt.Sprintf("%s:%s", host, port)

	// setup store
	mongo, err := store.NewMongoClient(config.MongoConnString)
	if err != nil {
		logger.Fatal(err)
	}
	store := store.NewStore(mongo, logger)

	// start tcp listener
	lis, err := net.Listen("tcp", addr)
	if err != nil {
		logger.Fatalf("failed to listen: %v", err)
	}

	// setup grpc server
	server := grpc.NewServer(
		grpc.UnaryInterceptor(AuthInterceptor()),
	)
	loggerciseHandler := loggerciseHandlerPackage.NewLoggerciseHandler(logger, store)
	loggerciseService := services.NewLoggerciseService(logger, loggerciseHandler)
	loggerciseProtoPackage.RegisterLoggerciseServer(server, loggerciseService)

	// setup grpc gateway
	go func() {
		mux := runtime.NewServeMux(
			runtime.WithIncomingHeaderMatcher(runtime.DefaultHeaderMatcher),
		)
		logger.Infof("http server: %s", httpAddr)
		loggerciseProtoPackage.RegisterLoggerciseHandlerServer(context.Background(), mux, loggerciseService)

		logger.Fatal(http.ListenAndServe(httpAddr, GrpcGatewayInterceptor(mux, logger)))
	}()

	reflection.Register(server)
	logger.Infof("grpc server: %s", addr)
	server.Serve(lis)
}

func AuthInterceptor() grpc.UnaryServerInterceptor {
	return func(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
		return handler(ctx, req)
	}
}

func GrpcGatewayInterceptor(mux http.Handler, log *logrus.Logger) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		log.Infof("%s %s", r.Method, r.URL.Path)
		mux.ServeHTTP(w, r)
	})
}
