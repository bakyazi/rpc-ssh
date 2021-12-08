package main

import (
	"context"
	"fmt"
	"github.com/bakyazi/rpc-ssh/api"
	"google.golang.org/grpc"
	"google.golang.org/grpc/peer"
	"net"
)

var Hosts HostDefinitions

func main() {
	lis, err := net.Listen("tcp", "localhost:8080")
	if err != nil {
		panic(err)
	}

	Hosts, err = LoadHosts()
	if err != nil {
		panic(err)
	}

	srv := grpc.NewServer()
	api.RegisterSshServiceServer(srv, &sshCommandService{})
	fmt.Println("starting server...")
	panic(srv.Serve(lis))
}

type sshCommandService struct {
	api.UnimplementedSshServiceServer
}

func (s *sshCommandService) RunCommand(c context.Context, sc *api.SshCommand) (*api.SshResponse, error) {
	p, _ := peer.FromContext(c)

	fmt.Printf("GET ssh command from peer=%s [Parameters=(user:%s|host:%s|cmd:\"%s\")]", p.Addr, sc.Username, sc.Host, sc.Command)
	output, err := RunSsh(sc.Username, sc.Host, sc.Command)
	fmt.Printf(" --> %t\n", err == nil)
	if err != nil {
		return nil, err
	}
	return &api.SshResponse{Response: *output}, nil
}
