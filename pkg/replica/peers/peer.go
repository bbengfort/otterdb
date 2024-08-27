package peers

import (
	"context"
	"fmt"
	"sync"

	"github.com/bbengfort/otterdb/pkg/replica/raft/v1"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

// Peer represents a replica in a distributed consensus quorum and provides connection
// functionality to maintain a remote connection to that replica for RPCs.
type Peer struct {
	PID    uint16 `json:"pid"`              // The precedence id of the peer
	Name   string `json:"name"`             // The unique name of the replica in the quorum
	Addr   string `json:"addr"`             // The dial address of the peer including port
	Region string `json:"region,omitempty"` // The region that the peer is located in

	sync.RWMutex
	conn   *grpc.ClientConn // grpc dial connection to the remote
	client raft.RaftClient  // grpc raft client
}

//===========================================================================
// Network Connection and RPCs
//===========================================================================

func (p *Peer) Connect(opts ...grpc.DialOption) (err error) {
	p.Lock()
	defer p.Unlock()

	if p.Addr == "" {
		return ErrNoEndpoint
	}

	if p.conn != nil {
		return ErrAlreadyConnected
	}

	// If no options are specified, connect with an insecure client
	if len(opts) == 0 {
		opts = append(opts, grpc.WithTransportCredentials(insecure.NewCredentials()))
	}

	if p.conn, err = grpc.NewClient(p.Addr, opts...); err != nil {
		return fmt.Errorf("could not connect to %s: %w", p.Name, err)
	}

	p.client = raft.NewRaftClient(p.conn)
	return nil
}

func (p *Peer) Close() (err error) {
	p.Lock()
	defer p.Unlock()

	err = p.conn.Close()

	p.conn = nil
	p.client = nil
	return err
}

func (p *Peer) RequestVote(ctx context.Context, in *raft.VoteRequest) (*raft.VoteReply, error) {
	if p.client == nil {
		return nil, ErrNotConnected
	}

	return p.client.RequestVote(ctx, in)
}

func (p *Peer) AppendEntries(ctx context.Context, in *raft.AppendRequest) (*raft.AppendReply, error) {
	if p.client == nil {
		return nil, ErrNotConnected
	}

	return p.client.AppendEntries(ctx, in)
}
