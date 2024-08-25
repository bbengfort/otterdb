/*
Package quorum implements configuration details for a set of hosts who must
coordinate together to make decisions. This package implements a set-like
structure that identifies unique participants as well as other details like
if a quorum contains a leader. The package also provides additional
mechanisms for votes and elections that relate to quorum operation.

Note that none of the data structures in this package are thread-safe and it
is expected that data structures that contain these methods wrap them with
thread safety.
*/
package quorum

import "github.com/bbengfort/otterdb/pkg/replica/sequence"

// initialize the package on import
func init() {
	idSequence = sequence.New()
}

var idSequence *sequence.Sequence // monotonically increasing quorum id sequence
var exists = struct{}{}           // helpful to not write everywhere struct{}{}

// New creates a quorum from the specified hosts with a unique id.
func New(hosts ...string) (quorum *Quorum) {
	quorum = &Quorum{
		id:    idSequence.Next(),
		hosts: make(map[string]struct{}, len(hosts)),
	}

	for _, host := range hosts {
		quorum.add(host)
	}

	return quorum
}

// Quorum represents a set of hosts that are configured to work together to
// make decisions. This base structure can be refined on a per-consensus
// basis, e.g. for leader-oriented quorums or other quorum types. Hosts are
// identified by string, e.g. by some hashable key or by a network address.
//
// Quorums cannot be changed once they are initialized, though there are some
// operations that can create a new quorum from old quorums -- they will
// receive a unique id. The quorum id is a monotonically increasing number,
// all quorums in the same process will have a unique id.
type Quorum struct {
	id    uint64              // The quorum id, which must be unique in the system
	hosts map[string]struct{} // The set of hosts that participate in the quorum
}

func (q *Quorum) Election() *Election {
	return NewElection(q)
}

//===========================================================================
// Read-Only Access to Quorum Properties
//===========================================================================

// ID returns the quorums unique ID
func (q *Quorum) ID() uint64 {
	return q.id
}

// Hosts is a read-only accessor for the underlying hosts array.
func (q *Quorum) Hosts() (hosts []string) {
	hosts = make([]string, 0, len(q.hosts))
	for host := range q.hosts {
		hosts = append(hosts, host)
	}
	return hosts
}

// Size returns the number of hosts in the quorum.
func (q *Quorum) Size() int {
	return len(q.hosts)
}

//===========================================================================
// Quorum Set Methods
//===========================================================================

// Contains returns true if the host is a member of the quorum.
func (q *Quorum) Contains(host string) bool {
	_, ok := q.hosts[host]
	return ok
}

// IsSubset tests whether q is a subset of r.
func (q *Quorum) IsSubset(r *Quorum) bool {
	for host := range q.hosts {
		if _, ok := r.hosts[host]; !ok {
			return false
		}
	}
	return true
}

// IsSuperset tests whether r is a superset of q.
func (q *Quorum) IsSuperset(r *Quorum) bool {
	return r.IsSubset(q)
}

// Intersects tests whether r and q have hosts in common.
func (q *Quorum) Intersects(r *Quorum) bool {
	for host := range q.hosts {
		if _, ok := r.hosts[host]; ok {
			return true
		}
	}
	return false
}

//===========================================================================
// Quorum Modification Methods
//===========================================================================

// add a host to the quorum (internal helper method only)
func (q *Quorum) add(host string) {
	if host != "" {
		q.hosts[host] = exists
	}
}
