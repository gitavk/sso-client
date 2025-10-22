package main

import (
	"context"
	"fmt"
	"log"
	"time"

	pb "github.com/gitavk/sso-proto/auth"
	"github.com/golang-jwt/jwt/v5"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func main() {
	// Connect to server
	conn, err := grpc.Dial("localhost:50051", grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("Failed to connect: %v", err)
	}
	defer conn.Close()

	client := pb.NewAuthServiceClient(conn)

	// Login
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()

	resp, err := client.Login(ctx, &pb.LoginRequest{
		Username: "admin",
		Password: "password123",
	})
	if err != nil {
		log.Fatalf("Login failed: %v", err)
	}

	if resp.Error != "" {
		log.Fatalf("Login error: %s", resp.Error)
	}

	fmt.Println("Login successful!")
	fmt.Printf("JWT Token: %s\n\n", resp.Token)

	// Parse and display JWT claims
	token, _, err := new(jwt.Parser).ParseUnverified(resp.Token, jwt.MapClaims{})
	if err != nil {
		log.Fatalf("Failed to parse token: %v", err)
	}

	if claims, ok := token.Claims.(jwt.MapClaims); ok {
		fmt.Println("Token Claims:")
		fmt.Printf("  Username: %v\n", claims["username"])
		fmt.Printf("  Issued At: %v\n", time.Unix(int64(claims["iat"].(float64)), 0))
		fmt.Printf("  Expires At: %v\n", time.Unix(int64(claims["exp"].(float64)), 0))
	}
}
