package main

import (
	"context"
	"fmt"
	"log"
	"net"

	// 生成されたパッケージをインポート (go.modのモジュール名 + パス)
	pb "go-grpc/pkg/grpc"

	"google.golang.org/grpc"
)

// 1. サーバーの構造体を定義
type server struct {
	// これを埋め込むことで、将来APIが増えてもコンパイルエラーになりにくくなる（必須のお作法）
	pb.UnimplementedGreeterServer
}

// 2. SayHello メソッドを実装（これが実際の処理！）
func (s *server) SayHello(ctx context.Context, req *pb.HelloRequest) (*pb.HelloResponse, error) {
	// クライアントから送られてきた名前を取り出す
	name := req.GetName()

	// 返事を作る
	message := fmt.Sprintf("Hello, %s!", name)

	// レスポンスを返す
	return &pb.HelloResponse{
		Message: message,
	}, nil
}

func main() {
	// 3. ポート8080で待ち受けを開始
	lis, err := net.Listen("tcp", ":8080")
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	// 4. gRPCサーバーを作成
	s := grpc.NewServer()

	// 5. 生成されたコードの関数を使って、自作の server を登録
	pb.RegisterGreeterServer(s, &server{})

	log.Printf("server listening at %v", lis.Addr())

	// 6. サーバーを起動
	if err := s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
