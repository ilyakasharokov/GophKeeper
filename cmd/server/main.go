package main

import (
	"context"
	"fmt"
	"gophkeeper/cmd/server/configuration"
	"gophkeeper/internal/app/server/authorization"
	"gophkeeper/internal/app/server/database"
	"os"
	"os/signal"
	"syscall"
	"time"

	grpc_middleware "github.com/grpc-ecosystem/go-grpc-middleware"
	grpc_auth "github.com/grpc-ecosystem/go-grpc-middleware/auth"
	"github.com/jmoiron/sqlx"

	grpcserver "gophkeeper/internal/app/server/grpc"
	"gophkeeper/internal/app/server/service"
	proto "gophkeeper/pkg/grpc/proto"
	"net"

	"google.golang.org/grpc"
)

func main() {
	ctx, cancel := context.WithCancel(context.Background())
	cfg := configuration.New()
	s := grpc.NewServer(
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
		s.Stop()
		cancel()
	case <-ctx.Done():
	}
	_, cancelt := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancelt()

}

/*func handleSignals(cancel context.CancelFunc, serverCancel func(ctx context.Context) error) {
	sigint := make(chan os.Signal, 1)
	signal.Notify(sigint, os.Interrupt)
	ctxt, cancelt := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancelt()
	select {
	case <-sigint:
		cancel()
		err := serverCancel(ctxt)
		if err != nil {
			log.Info().Err(err)
		}
		return
	}
}
*/