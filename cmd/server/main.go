package main

import (
	"fmt"
	"log"
	"mygrpc/config"
	"mygrpc/internal/handler"
	"mygrpc/internal/infra/persistence"
	"mygrpc/internal/usecase"
	todopb "mygrpc/pkg/grpc"
	"net"
	"os"
	"os/signal"

	"google.golang.org/grpc/reflection"

	// gRPCサーバーに対応する型が入っている
	"google.golang.org/grpc"
)

func main() {
	// 1. 8080番portのLisnterを作成
	port := 8080
	listener, err := net.Listen("tcp", fmt.Sprintf(":%d", port))
	if err != nil {
		panic(err)
	}

	// 2. gRPCサーバーを作成
	s := grpc.NewServer()

	// 依存性の注入
	db, _ := config.GetConn()
	todoPersistence := persistence.NewTodoPersistence(db)
	todoUseCase := usecase.NewTodoUseCase(todoPersistence)

	// 3. gRPCサーバーに登録
	todopb.RegisterTodoServiceServer(s, handler.NewTodoServer(todoUseCase))

	// 4. サーバーリフレクションの設定
	reflection.Register(s)

	// 5. 作成したgRPCサーバーを、8080番ポートで稼働させる
	go func() {
		log.Printf("start gRPC server port: %v", port)
		s.Serve(listener)
	}()

	// 6.Ctrl+Cが入力されたらGraceful shutdownされるようにする
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt)
	<-quit
	log.Println("stopping gRPC server...")
	s.GracefulStop()
}
