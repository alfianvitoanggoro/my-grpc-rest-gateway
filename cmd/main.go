package main

import (
	"context"
	"flag"
	"log"
	"net/http"

	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"google.golang.org/grpc"
	"google.golang.org/grpc/grpclog"

	bankgw "github.com/AlfianVitoAnggoro/my-grpc-proto/protogen/gateway/go/proto/bank"       // Update
	hellogw "github.com/AlfianVitoAnggoro/my-grpc-proto/protogen/gateway/go/proto/hello"     // Update
	reslgw "github.com/AlfianVitoAnggoro/my-grpc-proto/protogen/gateway/go/proto/resiliency" // Update
	"google.golang.org/grpc/credentials"
)

var (
	// command-line options:
	// gRPC server endpoint
	grpcServerEndpoint = flag.String("grpc-server-endpoint", "localhost:9090", "gRPC server endpoint")
)

func run() error {
	ctx := context.Background()
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	// Register gRPC server endpoint
	// Note: Make sure the gRPC server is running properly and accessible
	mux := runtime.NewServeMux()

	var opts []grpc.DialOption

	creds, err := credentials.NewClientTLSFromFile("ssl/ca.crt", "")

	if err != nil {
		log.Fatalln("Can't create client credentials :", err)
	}

	opts = append(opts, grpc.WithTransportCredentials(creds))
	// opts := []grpc.DialOption{grpc.WithTransportCredentials(insecure.NewCredentials())}
	if err := hellogw.RegisterHelloServiceHandlerFromEndpoint(ctx, mux, *grpcServerEndpoint, opts); err != nil {
		return err
	}

	if err := reslgw.RegisterResiliencyServiceHandlerFromEndpoint(ctx, mux, *grpcServerEndpoint, opts); err != nil {
		return err
	}

	if err := bankgw.RegisterBankServiceHandlerFromEndpoint(ctx, mux, *grpcServerEndpoint, opts); err != nil {
		return err
	}

	// Start HTTP server (and proxy calls to gRPC server endpoint)
	return http.ListenAndServe(":8081", mux)
}

func main() {
	flag.Parse()

	if err := run(); err != nil {
		grpclog.Fatal(err)
	}
}
