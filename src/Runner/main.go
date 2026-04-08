package main

import (
	"context"
	"log"
	"os"
	"time"

	"github.com/Jakkoble/HexaTask/src/Runner/internal"
	"github.com/Jakkoble/HexaTask/src/Runner/pb"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

//go:generate protoc --go_out=pb --go_opt=paths=source_relative --go-grpc_out=pb --go-grpc_opt=paths=source_relative -I ../../contracts ../../contracts/runner.proto

func main() {
	jobID := os.Getenv("JOB_ID")
	commanderURL := os.Getenv("COMMANDER_URL")

	if jobID == "" || commanderURL == "" {
		log.Fatalf("JOB_ID or COMMANDER_URL is missing.")
	}

	log.Printf("Runner started for Job ID: %s", jobID)

	conn, err := grpc.NewClient(commanderURL, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("Failed to connect to Commander: %v", err)
	}
	defer conn.Close()

	client := pb.NewRunnerServiceClient(conn)

	setupCtx, setupCancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer setupCancel()

	// fetch job from commander
	fetcher := internal.NewGrpcJobFetcher(client)
	commands, err := fetcher.FetchCommands(setupCtx, jobID)
	if err != nil {
		log.Fatalf("Failed to fetch job details: %v", err)
	}

	// setup streaming connection
	execCtx := context.Background()
	stream, err := client.StreamLogs(execCtx)
	if err != nil {
		log.Fatalf("Failed to open log stream: %v", err)
	}

	// execute commands
	executor := internal.NewTaskExecutor(
		&internal.ShellCommandRunner{},
		internal.NewGrpcLogStreamer(stream, jobID),
	)

	executor.Execute(execCtx, commands)
}
