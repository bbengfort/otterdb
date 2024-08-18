/*
Package replica manages replication between peers in an otterdb cluster.
*/
package replica

import "github.com/bbengfort/otterdb/pkg/config"

type Replica struct {
	conf config.ReplicaConfig
}

func New(conf config.ReplicaConfig) (srv *Replica, err error) {
	return &Replica{conf: conf}, nil
}

func (r *Replica) Serve(errc chan<- error) (err error) {
	return nil
}

func (r *Replica) Shutdown() (err error) {
	return nil
}
