package main

import (
	"fmt"
	"io"

	"math"
	"math/rand"
	"net/http"
	"time"

	gossipv1 "github.com/chronark/crdt/gen/proto/gossip/v1"
	"github.com/chronark/crdt/pkg/logging"
	"github.com/chronark/crdt/pkg/server"
)

func main() {

	logger := logging.NewLogger()

	servers := []*server.Service{}
	for port := 9001; port <= 9005; port++ {
		server := server.New("0.0.0.0", int64(port), logger)
		server.Run()
		servers = append(servers, server)

	}

	time.Sleep(time.Second)
	for _, me := range servers {
		for _, peer := range servers {
			if me.ID == peer.ID {
				continue
			}

			err := me.AddPeer(peer.ID, peer.Host, peer.Port)
			if err != nil {
				panic(err)
			}

		}
	}

	clients := []gossipv1.GossipServiceClient{}
	for _, srv := range servers {
		client, close, err := server.NewClient(fmt.Sprintf("%s:%d", srv.Host, srv.Port))
		if err != nil {
			panic(err)
		}
		defer func() {
			err := close()
			if err != nil {
				panic(err)
			}
		}()
		clients = append(clients, client)
	}

	n := 10000
	success := 0
	for i := 0; i < n; i++ {
		logger.Info().Float64("%", math.RoundToEven(100*float64(i)/float64(n))).Send()
		srv := servers[rand.Intn(len(servers))]

		res, err := http.Post(fmt.Sprintf("http://%s:%d/ratelimit", srv.Host, srv.Port+1000), "", nil)
		if err != nil {
			panic(err)
		}
		defer res.Body.Close()
		body, err := io.ReadAll(res.Body)
		if err != nil {
			panic(err)
		}
		if res.StatusCode == 200 {
			success++
		}
		logger.Info().Bytes("body", body).Int("status", res.StatusCode).Send()
	}

	logger.Info().Int("success", success).Send()

	time.Sleep(time.Second * 2)

}
