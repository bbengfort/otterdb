package peers

import "errors"

var (
	ErrNoEndpoint       = errors.New("peer does not have an endpoint to connect on")
	ErrAlreadyConnected = errors.New("already connected to remote peer")
	ErrNotConnected     = errors.New("not connected to remote peer")
)
