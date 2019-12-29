package main

import (
	"context"
	"fmt"
	demo "github.com/brharrelldev/dexauthpoc/api"
	"github.com/coreos/go-oidc"
	grpc_auth "github.com/grpc-ecosystem/go-grpc-middleware/auth"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"log"
	"net"
	"os"
)

type Server struct {
}

func (s *Server) DemoService(ctx context.Context, request *demo.DemoRequest) (*demo.DemoResponse, error) {

	return &demo.DemoResponse{
		Message: request.Message,
	}, nil

}

func StartServer() error {

	lis, err := net.Listen("tcp", ":50051")
	if err != nil {
		return fmt.Errorf("error estabilishing listener %v", err)
	}

	g := grpc.NewServer(grpc.UnaryInterceptor(grpc_auth.UnaryServerInterceptor(authFunction)))

	s := Server{}
	demo.RegisterDemoServiceServer(g, &s)

	log.Println("starting server")
	return g.Serve(lis)
}

func authFunction(ctx context.Context) (context.Context, error) {
	token, err := grpc_auth.AuthFromMD(ctx, "bearer")
	if err != nil {
		return nil, status.Error(codes.Unauthenticated, fmt.Sprintf("could not get new token due to %v", err))
	}

	fmt.Println(token)

	provider, err := oidc.NewProvider(ctx, "http://127.0.0.1:5556/dex")
	if err != nil {
		return nil, status.Error(codes.Unauthenticated, fmt.Sprintf("error reaching provider %v", err))
	}

	fmt.Println(provider)







	//idToken, err := verifier.Verify(ctx, token)
	//if err != nil {
	//	return nil, status.Error(codes.Unauthenticated, fmt.Sprintf("error verifying idToken %v", err))
	//}
	//
	//fmt.Println(idToken)

	return nil, nil

}

func main() {

	sigChan := make(chan os.Signal, 1)
	errChan := make(chan error)

	go func() {
		if err := StartServer(); err != nil {
			errChan <- fmt.Errorf("error starting grpc server %v", err)
			log.Fatalf("could not start server %v", err)
		}
		<-sigChan

	}()

	select {
	case <-sigChan:
		log.Println("server stopped")
		os.Exit(0)
	case err := <-errChan:
		log.Fatalf("error occured due to %v", err)

	}
}
