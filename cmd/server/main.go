package main

import (
	"TapMars/productManager/pkg/config"
	"TapMars/productManager/pkg/productManager"
	pb "TapMars/productManager/pkg/proto"
	"context"
	"google.golang.org/grpc"
	"log"
	"net"
)

//Failing to connect to Firestore is the integration test
func main() {
	log.Printf("starting server...")

	ctx := context.Background()

	port, projectID, err := config.GetEnvironmentVariables()
	if err != nil {
		log.Fatalf("failed to config: %v", err)
	}
	log.Printf("Port: %s", port)
	log.Printf("ProjectID: %s", projectID)

	lis, err := net.Listen("tcp", ":"+port)
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	srv, err := productManager.NewServer(ctx, projectID)
	if err != nil {
		log.Fatalf("Firestore Startup Server Error: %v", err)
	}
	log.Printf("New productManager Server")

	s := grpc.NewServer()

	pb.RegisterProductManagerServer(s, srv)
	if err = s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
