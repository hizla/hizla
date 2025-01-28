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
	Address    string `json:"address"`
	DbHost     string `json:"db_host"`
	DbPort     string `json:"db_port"`
	DbUser     string `json:"db_user"`
	DbPassword string `json:"db_password"`
	DbName     string `json:"db_name"`
}

func (s *ServeAPI) FromEnviron() error {
	if addr := os.Getenv("HIZLA_API_LISTEN_ADDRESS"); addr != "" {
		s.Address = addr
	} else {
		return ErrUnsetAddress
	}
	return nil
}
