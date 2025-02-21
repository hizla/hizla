// Package hst contains shared types useful to third party programs wishing to integrate with hizla.
package hst

import (
	"errors"
	"os"
)

var (
	ErrUnsetAddress = errors.New("listen address not present in environment")
)

type ServeAPI struct {
	Address string `json:"address"`
}

func (s *ServeAPI) FromEnviron() error {
	if addr := os.Getenv("HIZLA_API_LISTEN_ADDRESS"); addr != "" {
		s.Address = addr
	} else {
		return ErrUnsetAddress
	}
	return nil
}
