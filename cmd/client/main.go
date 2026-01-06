package main

import (
	"context"
	"flag"
	"log"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"

	// 生成されたパッケージをインポート
	pb "go-grpc/pkg/grpc"
)

func main() {
	// コマンドライン引数で名前を指定できるようにする（デフォルトは "World"）
	name := flag.String("name", "World", "Name to greet")
	flag.Parse()

	// 1. サーバーに接続 (localhost:8080)
	// ※ 本番ではSSL/TLSを使いますが、練習なので insecure (暗号化なし) を使います
	conn, err := grpc.NewClient("localhost:8080", grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}
	defer conn.Close()

	// 2. クライアントを作成
	c := pb.NewGreeterClient(conn)

	// 3. タイムアウトを設定（1秒以内に返事がなければ諦める）
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	// 4. サーバーの SayHello メソッドを呼び出す
	r, err := c.SayHello(ctx, &pb.HelloRequest{Name: *name})
	if err != nil {
		log.Fatalf("could not greet: %v", err)
	}

	// 5. 結果を表示
	log.Printf("Greeting: %s", r.GetMessage())
}
