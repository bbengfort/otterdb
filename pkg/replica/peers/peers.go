package peers

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"os"

	"github.com/bbengfort/otterdb/pkg/replica/raft/v1"
	"github.com/rs/zerolog/log"
	"google.golang.org/grpc"
)

// Peers is a collection of replica information that describes a quorum.
type Peers []*Peer

//===========================================================================
// Accessors
//===========================================================================

// Names returns all the names of every peer in the collection.
func (p Peers) Names() []string {
	names := make([]string, 0, len(p))
	for _, peer := range p {
		names = append(names, peer.Name)
	}
	return names
}

// Get a peer by name; returns an error if no peer with that name exists.
func (p Peers) Get(name string) (_ *Peer, err error) {
	for _, peer := range p {
		if name == peer.Name {
			return peer, nil
		}
	}
	return nil, fmt.Errorf("no peer found named %q", name)
}

// Presiding returns the name of the replica with the lowest pid.
// NOTE: a PID of zero is not allowed and is ignored.
func (p Peers) Presiding() string {
	var pid uint16
	var name string
	for _, peer := range p {
		if pid == 0 || (peer.PID > 0 && peer.PID < pid) {
			pid = peer.PID
			name = peer.Name
		}
	}
	return name
}

//===========================================================================
// Networking Methods
//===========================================================================

// Blocking connect method that connects to all peers in the group and returns any
// errors that might occur while connecting to the remotes.
func (p Peers) Connect(opts ...grpc.DialOption) (err error) {
	for _, peer := range p {
		if cerr := peer.Connect(opts...); cerr != nil {
			err = errors.Join(err, cerr)
		}
	}
	return err
}

// Blocking close method that closes all peer connections in the group and returns any
// errors that might occur while closing the connection to the remotes.
func (p Peers) Close() (err error) {
	for _, peer := range p {
		if cerr := peer.Close(); cerr != nil {
			err = errors.Join(err, cerr)
		}
	}
	return err
}

// Broadcast request vote to all remote peers, sending replies on the recv channel.
// Any errors that occur will be logged after all handling is complete.
func (p Peers) RequestVote(ctx context.Context, in *raft.VoteRequest, recv chan<- *raft.VoteReply) {
	// TODO: use a send channel rather than starting an arbitrary number of go routines.
	for _, peer := range p {
		go func() {
			out, err := peer.RequestVote(ctx, in)
			if err != nil {
				log.Warn().Err(err).Msg("request vote rpc failed")
				return
			}
			recv <- out
		}()
	}
}

func (p Peers) AppendEntries(ctx context.Context, in *raft.AppendRequest, recv chan<- *raft.AppendReply) {
	// TODO: use a send channel rather than starting an arbitrary number of go routines.
	for _, peer := range p {
		go func() {
			out, err := peer.AppendEntries(ctx, in)
			if err != nil {
				log.Warn().Err(err).Msg("append entries rpc failed")
				return
			}
			recv <- out
		}()
	}
}

//===========================================================================
// Serialization and Deserialization of Objects
//===========================================================================

// Load the peers from a path on disk
func Load(path string) (peers Peers, err error) {
	var f *os.File
	if f, err = os.Open(path); err != nil {
		return nil, err
	}
	defer f.Close()

	peers = make(Peers, 0)
	if err = json.NewDecoder(f).Decode(&peers); err != nil {
		return nil, err
	}

	return peers, nil
}

// Dump the peers to disk. Will create the file at the specified path with
// 0644 permissions if it doesn't exist otherwise will truncate and replace.
func (p Peers) Dump(path string) (err error) {
	var f *os.File
	if f, err = os.Create(path); err != nil {
		return err
	}

	if err = json.NewEncoder(f).Encode(p); err != nil {
		return err
	}
	return nil
}
