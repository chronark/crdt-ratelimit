package server

import (
	pb "github.com/chronark/crdt/gen/proto/gossip/v1"

	"google.golang.org/grpc"
)

func NewClient(address string) (pb.GossipServiceClient, func() error, error) {

	conn, err := grpc.Dial(address, grpc.WithInsecure())

	if err != nil {
		return nil, nil, err
	}

	client := pb.NewGossipServiceClient(conn)

	return client, conn.Close, nil

}
