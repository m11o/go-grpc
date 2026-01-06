package main

import (
	"context"
	"io"
	"log"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"

	pb "go-grpc/pkg/grpc"
)

func main() {
	// 1. サーバーに接続
	conn, err := grpc.NewClient("localhost:8080", grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}
	defer conn.Close()
	c := pb.NewGreeterClient(conn)

	// 2. 双方向ストリーミングを開始
	// (コンテキストとストリームの作成)
	stream, err := c.Chat(context.Background())
	if err != nil {
		log.Fatalf("error creating stream: %v", err)
	}

	// 終了待ち受け用のチャネル
	waitc := make(chan struct{})

	// 3. 【送信担当】別のゴルーチンでメッセージを送り続ける
	go func() {
		names := []string{"Alice", "Bob", "Charlie", "Dave", "Eve"}
		for _, name := range names {
			log.Printf("Sending message: %s", name)

			// メッセージ送信
			if err := stream.Send(&pb.HelloRequest{Name: name}); err != nil {
				log.Fatalf("failed to send: %v", err)
			}

			// 1秒待ってから次を送る（チャットっぽくするため）
			time.Sleep(1 * time.Second)
		}
		// すべて送り終わったら「送信終了」をサーバーに伝える
		stream.CloseSend()
	}()

	// 4. 【受信担当】メインスレッドでサーバーからの返事を受け取り続ける
	go func() {
		for {
			in, err := stream.Recv()
			if err == io.EOF {
				// サーバーからの通信が終わったらループを抜ける
				close(waitc)
				return
			}
			if err != nil {
				log.Fatalf("failed to receive: %v", err)
			}
			log.Printf("Received: %s", in.GetMessage())
		}
	}()

	// 5. 受信が終わるまでここで待機
	<-waitc
}
