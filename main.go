package main

import (
	"7-solutions/config"
	grpcserver "7-solutions/grpc"
	userpb "7-solutions/proto"

	"7-solutions/handler"
	"7-solutions/middleware"

	"7-solutions/repository"
	"7-solutions/usecase"
	"context"
	"log"
	"net"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"google.golang.org/grpc"
)

func main() {
	_ = godotenv.Load()

	client := config.ConnectDB()
	db := client.Database("7-solutions-db")
	userRepo := repository.NewUserRepository(db)
	userUC := usecase.NewUserUsecase(userRepo, os.Getenv("JWT_SECRET"))

	go func() {
		for {
			count, err := userUC.CountUsers(nil)
			if err == nil {
				log.Printf("Total users: %d\n", count)
			}
			time.Sleep(10 * time.Second)
		}
	}()

	// ! === Setup Gin HTTP Server ===
	ginRouter := gin.Default()
	handler.NewUserHandler(ginRouter, userUC, middleware.JWTAuth(os.Getenv("JWT_SECRET"), userRepo))
	httpSrv := &http.Server{
		Addr:    ":8080",
		Handler: ginRouter,
	}

	// ? === Setup gRPC Server ===
	authInterceptor := grpcserver.NewAuthInterceptor(os.Getenv("JWT_SECRET"))
	grpcSrv := grpc.NewServer(
		grpc.UnaryInterceptor(authInterceptor.Unary()),
	)
	userpb.RegisterUserServiceServer(grpcSrv, grpcserver.NewUserGRPCServer(userUC))
	grpcListener, err := net.Listen("tcp", ":50051")
	if err != nil {
		log.Fatalf("failed to listen for gRPC: %v", err)
	}

	go func() {
		log.Println("HTTP server is running at http://localhost:8080")
		if err := httpSrv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("HTTP server error: %v", err)
		}
	}()

	go func() {
		log.Println("gRPC server is running at :50051")
		if err := grpcSrv.Serve(grpcListener); err != nil {
			log.Fatalf("gRPC server error: %v", err)
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	log.Println("Shutting down servers...")

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	if err := httpSrv.Shutdown(ctx); err != nil {
		log.Printf("HTTP forced to shutdown: %v", err)
	}

	grpcSrv.GracefulStop()

	log.Println("Servers exited gracefully")
}
