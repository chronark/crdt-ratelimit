package server

import (
	"context"
	"fmt"
	"github.com/labstack/echo/v4"
	"net"
	"os"
	"sync"
	"time"

	pb "github.com/chronark/crdt/gen/proto/gossip/v1"
	"github.com/chronark/crdt/pkg/crdt"
	"github.com/google/uuid"
	"github.com/rs/zerolog"

	"google.golang.org/grpc"
)

type Peer struct {
	id       string
	host     string
	port     int64
	client   pb.GossipServiceClient
	close    func() error
	lock     *sync.RWMutex
	lastPush int64
}

type Service struct {
	lock    *sync.RWMutex
	ID      string
	Host    string
	Port    int64
	Counter map[string]*crdt.GCounter
	peers   map[string]*Peer
	logger  zerolog.Logger
	pb.UnimplementedGossipServiceServer
}

func New(Host string, Port int64, logger zerolog.Logger) *Service {
	ID := uuid.NewString()
	Counter := map[string]*crdt.GCounter{}
	peers := map[string]*Peer{}
	logger = logger.With().Str("ID", ID).Logger()
	lock := &sync.RWMutex{}
	return &Service{
		lock,
		ID,
		Host,
		Port,
		Counter,
		peers,
		logger,
		pb.UnimplementedGossipServiceServer{},
	}
}

func (svc *Service) AddPeer(id, host string, port int64) error {
	svc.logger.Info().Str("rpc", "Join").Str("id", id).Str("host", host).Int64("port", port).Send()

	for knownID := range svc.peers {
		if knownID == id {
			svc.logger.Warn().Str("id", id).Msg("Peer has already joined")
			return nil
		}
	}

	client, close, err := NewClient(fmt.Sprintf("%s:%d", host, port))
	if err != nil {
		return err
	}
	lastPush := time.Now().UnixNano()
	lock := &sync.RWMutex{}
	svc.peers[id] = &Peer{id, host, port, client, close, lock, lastPush}
	return nil

}

func (svc *Service) Merge(ctx context.Context, req *pb.MergeRequest) (*pb.MergeResponse, error) {
	other := crdt.GCounter{
		ID: req.GetCrdtId(),
		C:  map[int64]int64{},
	}
	for _, m := range req.GetMap() {
		other.C[m.GetKey()] = m.GetVal()
	}

	_, exists := svc.Counter[req.GetCrdtId()]
	if !exists {
		svc.Counter[req.GetCrdtId()] = crdt.NewGCounter(req.GetCrdtId())
	}
	svc.Counter[req.GetCrdtId()].Merge(other)
	return &pb.MergeResponse{}, nil
}

// }
// func (svc *Service) Increment(ctx context.Context, req *pb.IncrementRequest) (*pb.IncrementResponse, error) {
// 	svc.logger.Info().Str("rpc", "Increment").Send()

// 	err := svc.Counter.Increment(req.GetIncrement())
// 	if err != nil {
// 		return &pb.IncrementResponse{}, err
// 	}
// 	svc.lock.Lock()
// 	defer svc.lock.Unlock()

// 	if req.GetBroadcast() {

// 		errors := []error{}
// 		for _, peer := range svc.peers {
// 			peer.lock.Lock()
// 			last := peer.lastPush

// 			mergeReq := &pb.MergeRequest{
// 				Broadcast: true,
// 				CrdtId:    svc.Counter.ID,
// 				Map:       []*pb.Int64KV{},
// 			}
// 			for ts, incr := range svc.Counter.C {

// 				if ts > peer.lastPush {
// 					if ts > last {
// 						last = ts
// 					}
// 					mergeReq.Map = append(mergeReq.Map, &pb.Int64KV{Key: ts, Val: incr})
// 				}
// 			}
// 			svc.logger.Info().Int("changes", len(mergeReq.Map)).Int64("since", peer.lastPush).Int64("last", last).Msg("Syncing changes")

// 			_, err := peer.client.Merge(ctx, mergeReq)
// 			if err != nil {
// 				svc.logger.Error().Err(err).Send()
// 				errors = append(errors, err)
// 			}

// 			peer.lastPush = last

// 			peer.lock.Unlock()

// 		}
// 		if len(errors) > 0 {
// 			return &pb.IncrementResponse{}, errors[0]
// 		}
// 	}

// 	return &pb.IncrementResponse{}, nil

// }

func (svc *Service) ratelimit(id string) (bool, error) {

	limit := int64(100)

	_, exists := svc.Counter[id]
	if !exists {
		svc.Counter[id] = crdt.NewGCounter(id)
	}

	usage := svc.Counter[id].Value()
	if usage < limit {

		c := svc.Counter[id]
		svc.lock.Lock()
		defer svc.lock.Unlock()

		err := c.Increment(1)
		if err != nil {
			panic(err)
		}

		errors := []error{}
		for _, peer := range svc.peers {
			peer.lock.Lock()
			last := peer.lastPush

			mergeReq := &pb.MergeRequest{
				Broadcast: true,
				CrdtId:    c.ID,
				Map:       []*pb.Int64KV{},
			}
			for ts, incr := range c.C {

				if ts > peer.lastPush {
					if ts > last {
						last = ts
					}
					mergeReq.Map = append(mergeReq.Map, &pb.Int64KV{Key: ts, Val: incr})
				}
			}
			svc.logger.Info().Int("changes", len(mergeReq.Map)).Int64("since", peer.lastPush).Int64("last", last).Msg("Syncing changes")

			_, err := peer.client.Merge(context.TODO(), mergeReq)
			if err != nil {
				svc.logger.Error().Err(err).Send()
				errors = append(errors, err)
			}

			peer.lastPush = last

			peer.lock.Unlock()

		}
		if len(errors) > 0 {
			svc.logger.Fatal().Errs("errors", errors).Send()
		}
	}

	return usage < limit, nil

}

func (svc *Service) Run() {

	e := echo.New()
	e.POST("/ratelimit", func(c echo.Context) error {
		id := "id"
		allow, err := svc.ratelimit(id)
		if err != nil {
			return err
		}
		if allow {
			return c.String(200, "Allow")

		}
		return c.String(429, "Blocked")

	})
	svc.logger.Error().Str("id", svc.ID).Msg("Run")

	go func() {

		e.Start(fmt.Sprintf("%s:%d", svc.Host, svc.Port+1000))

	}()

	go func() {
		grpcServer := grpc.NewServer()

		pb.RegisterGossipServiceServer(grpcServer, svc)
		listener, err := net.Listen("tcp", fmt.Sprintf("%s:%d", svc.Host, svc.Port))
		if err != nil {
			svc.logger.Error().Err(err).Int64("port", svc.Port).Msg("Could not listen")
			os.Exit(1)
		}

		svc.logger.Info().Str("addr", listener.Addr().String()).Msg("Listening")
		err = grpcServer.Serve(listener)
		if err != nil {
			panic(err)
		}
	}()

}
