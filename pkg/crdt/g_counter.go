package crdt

import (
	"fmt"

	"time"
)

type GCounter struct {
	ID string
	C  map[int64]int64
}

func NewGCounter(id string) *GCounter {
	return &GCounter{
		ID: id,
		C:  map[int64]int64{},
	}
}

func (g *GCounter) Increment(i int64) error {
	if i <= 0 {
		return fmt.Errorf("i must be positive")
	}

	g.C[time.Now().UnixNano()] = i
	return nil
}

func (g *GCounter) Value() int64 {
	var total int64 = 0
	for _, v := range g.C {
		total += v
	}
	return total
}

func (g *GCounter) Merge(other GCounter) {
	for id, otherValue := range other.C {
		myValue, ok := g.C[id]
		// Update my value if it doesn't exist or is smaller than other
		if !ok || myValue < otherValue {
			g.C[id] = otherValue
		}
	}

}

// type server struct {
// 	ID      string
// 	counter *GCounter
// 	Peers   []*server
// }

// func NewServer(id string) *server {
// 	return &server{
// 		ID:      id,
// 		counter: NewGCounter(),
// 		Peers:   []*server{},
// 	}
// }

// func (s *server) AddPeers(peers ...*server) {
// 	s.Peers = append(s.Peers, peers...)
// }

// func (s *server) Receive(o GCounter) error {
// 	s.counter.Merge(o)
// 	fmt.Printf("%s state = %d\n", s.ID, s.counter.Value())

// 	return nil
// }

// func (s *server) Broadcast() error {
// 	for _, peer := range s.Peers {
// 		err := peer.Receive(*s.counter)
// 		if err != nil {
// 			return err
// 		}
// 	}
// 	return nil
// }

// func (s *server) Increment(i int64) error {
// 	err := s.counter.Increment(i)
// 	if err != nil {
// 		return err
// 	}
// 	fmt.Printf("%s state = %d\n", s.ID, s.counter.Value())

// 	return s.Broadcast()
// }

// func print(s ...*server) {
// 	fmt.Println()
// 	for _, srv := range s {
// 		fmt.Printf("%s: [%d] %+v\n", srv.ID, srv.counter.Value(), *srv)

// 	}
// 	fmt.Println()

// }

// func main() {

// 	a := NewServer("a")
// 	b := NewServer("b")
// 	c := NewServer("c")

// 	a.AddPeers(b, c)
// 	b.AddPeers(a, c)
// 	c.AddPeers(a, b)

// 	_ = a.Increment(5)
// 	_ = b.Increment(4)
// 	_ = b.Increment(3)
// 	_ = c.Increment(2)
// 	print(a, b, c)

// }
