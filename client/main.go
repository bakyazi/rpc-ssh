package main

import (
	"bufio"
	"context"
	"errors"
	"flag"
	"fmt"
	"github.com/bakyazi/rpc-ssh/api"
	"google.golang.org/grpc"
	"os"
)

var (
	server   = flag.String("server", "localhost:8080", "RPC-SSH Server Address in <ip:port> format")
	hostname = flag.String("host", "", "IP addr or Hostname of remote machine you want to execute command in")
	username = flag.String("user", "", "Username for SSH Login")
	command  = flag.String("command", "", "Command you want to execute")
	scanner  = bufio.NewScanner(os.Stdin)
)

func main() {
	flag.Parse()

	addr := *server
	conn, err := grpc.Dial(addr, grpc.WithInsecure())
	if err != nil {
		panic(err)
	}

	client := api.NewSshServiceClient(conn)
	ctx := context.Background()

	err = nil
	for err == nil && *command != "exit" {
		sendCommand(client, ctx)
		err = getCommandFromCommandLine()
	}

	fmt.Println("There is no provided command!")
	fmt.Println("RPC SSH CLI is exiting!")

	conn.Close()

}

func getCommandFromCommandLine() error {
	fmt.Println("Please enter new command (just type 'exit' to exit):")
	fmt.Printf(">>> ")
	if scanner.Scan() {
		*command = scanner.Text()
		return nil
	}
	return errors.New("cannot scan new command")
}

func sendCommand(client api.SshServiceClient, ctx context.Context) {
	resp, err := client.RunCommand(ctx, &api.SshCommand{Username: *username, Host: *hostname, Command: *command})
	fmt.Printf("Result of \"%s\" command:\n", *command)
	if err != nil {
		fmt.Printf("\033[1;33m%s\033[0m", err)
		fmt.Printf("\n")
		return
	}
	fmt.Printf("\033[1;33m%s\033[0m", resp.Response)
}
