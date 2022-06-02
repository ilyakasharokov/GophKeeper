package main

import (
	"context"
	"crypto/tls"
	"fmt"
	grpc_middleware "github.com/grpc-ecosystem/go-grpc-middleware"
	grpc_auth "github.com/grpc-ecosystem/go-grpc-middleware/auth"
	"github.com/jmoiron/sqlx"
	"google.golang.org/grpc/credentials"
	"gophkeeper/cmd/server/configuration"
	"gophkeeper/internal/app/server/authorization"
	"gophkeeper/internal/app/server/database"
	"os"
	"os/signal"
	"syscall"

	grpcserver "gophkeeper/internal/app/server/grpc"
	"gophkeeper/internal/app/server/service"
	proto "gophkeeper/pkg/grpc/proto"
	"net"

	_ "github.com/lib/pq"
	"google.golang.org/grpc"
)

var (
	BuildVersion = "N/A"
	BuildDate    = "N/A"
	BuildCommit  = "N/A"
)

func loadTLSCredentials() (credentials.TransportCredentials, error) {
	// Load server's certificate and private key
	serverCert, err := tls.LoadX509KeyPair("cert/server-cert.pem", "cert/server-key.pem")
	if err != nil {
		return nil, err
	}

	// Create the credentials and return it
	config := &tls.Config{
		Certificates: []tls.Certificate{serverCert},
		ClientAuth:   tls.NoClientCert,
	}

	return credentials.NewTLS(config), nil
}

func main() {
	fmt.Printf("Build version: %v\n", BuildVersion)
	fmt.Printf("Build date: %v\n", BuildDate)
	fmt.Printf("Build commit: %v\n", BuildCommit)
	ctx, cancel := context.WithCancel(context.Background())
	cfg := configuration.New()

	tlsCredentials, err := loadTLSCredentials()
	if err != nil {
		fmt.Println("cannot load TLS credentials: ", err)
		return
	}

	s := grpc.NewServer(
		grpc.Creds(tlsCredentials),
		grpc.UnaryInterceptor(
			grpc_middleware.ChainUnaryServer(
				grpc_auth.UnaryServerInterceptor(func(ctx context.Context) (context.Context, error) {
					userID, err := authorization.TokenValid(ctx, cfg.ACCESS_TOKEN_SECRET)
					if err != nil {
						userID = ""
					}
					newCtx := context.WithValue(ctx, "userID", userID)
					return newCtx, nil
				}),
			)))
	listen, err := net.Listen("tcp", ":3200")
	if err != nil {
		fmt.Println(err.Error())
		return
	}
	db, err := sqlx.Connect("postgres", cfg.DBDSN)
	if err != nil {
		fmt.Println(err.Error())
		return
	}
	d := database.NewDB(db)
	us := service.NewUserService(d, cfg.ACCESS_TOKEN_LIFETIME, cfg.REFRESH_TOKEN_LIFETIME, cfg.ACCESS_TOKEN_SECRET, cfg.REFRESH_TOKEN_SECRET)
	ss := service.NewSyncService(d)

	proto.RegisterGophKeeperServer(s, grpcserver.New(us, ss))

	go func() {
		fmt.Println("gRPC server started on :3200")
		if err := s.Serve(listen); err != nil {
			fmt.Println(err.Error())
			return
		}
	}()

	sigint := make(chan os.Signal, 1)
	signal.Notify(sigint,
		os.Interrupt,
		syscall.SIGINT,
		syscall.SIGTERM,
		syscall.SIGQUIT)
	select {
	case <-sigint:
		fmt.Println("stop grpc server")
		s.Stop()
		cancel()
	case <-ctx.Done():
	}
}
