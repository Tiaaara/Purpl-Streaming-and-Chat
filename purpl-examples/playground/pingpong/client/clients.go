package main

import (
	"bufio"
	"context"
	"fmt"
	"io"
	"log"
	"os"
	"strings"
	"time"

	"example.com/m/v2/pb"
	jwt "github.com/Siar-Akbayin/jwt-go-auth"
	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
)

func main() {
	// === Input user ===
	reader := bufio.NewReader(os.Stdin)

	fmt.Println("Pilih mode (unary / stream / chat):")
	modeRaw, _ := reader.ReadString('\n')
	mode := strings.TrimSpace(modeRaw)

	fmt.Println("Pilih policy (all_allowed / all_denied / mixed / maximized):")
	policyRaw, _ := reader.ReadString('\n')
	policyChoice := strings.TrimSpace(policyRaw)

	// mapping policy file
	policyFiles := map[string]string{
		"all_allowed": "client/policy_all_allowed.json",
		"all_denied":  "client/policy_all_denied.json",
		"mixed":       "client/policy_mixed.json",
		"maximized":   "client/policy_maximized.json",
	}
	policyFile, ok := policyFiles[policyChoice]
	if !ok {
		log.Fatalf("Policy tidak dikenal: %s", policyChoice)
	}

	// === Connect gRPC ===
	conn, err := grpc.Dial("localhost:50051", grpc.WithInsecure(), grpc.WithBlock())
	if err != nil {
		log.Fatalf("Failed to connect: %v", err)
	}
	defer conn.Close()
	c := pb.NewPingPongClient(conn)

	// === Generate token sesuai policy ===
	token, err := jwt.GenerateToken(
		policyFile,
		"trackingService-maximal",
		"purpose1",
		"client/private.pem",
		1,
	)
	if err != nil {
		log.Fatalf("Error generating token: %v", err)
	}
	ctx := metadata.AppendToOutgoingContext(context.Background(), "authorization", token)

	// === Jalankan RPC sesuai mode ===
	switch mode {
	case "unary":
		resp, err := c.SayHello(ctx, &pb.HelloRequest{Name: "unary-client"})
		if err != nil {
			log.Fatalf("Unary error: %v", err)
		}
		log.Printf("Unary: %s", resp)

	case "stream":
		stream, err := c.StreamHello(ctx, &pb.HelloRequest{Name: "stream-client"})
		if err != nil {
			log.Fatalf("Stream error: %v", err)
		}
		for {
			msg, err := stream.Recv()
			if err == io.EOF {
				break
			}
			if err != nil {
				log.Fatalf("Stream recv error: %v", err)
			}
			log.Printf("StreamHello: %s", msg)
		}

	case "chat":
		chat, err := c.Chat(ctx)
		if err != nil {
			log.Fatalf("Chat error: %v", err)
		}
		names := []string{"Alice", "Bob", "Charlie"}

		// kirim
		go func() {
			for _, n := range names {
				if err := chat.Send(&pb.HelloRequest{Name: n}); err != nil {
					log.Printf("Chat send error: %v", err)
					return
				}
				time.Sleep(500 * time.Millisecond)
			}
			chat.CloseSend()
		}()

		// terima
		for {
			reply, err := chat.Recv()
			if err == io.EOF {
				break
			}
			if err != nil {
				log.Fatalf("Chat recv error: %v", err)
			}
			log.Printf("Chat response: %s", reply)
		}

	default:
		log.Fatalf("Mode tidak dikenali: %s", mode)
	}
}
