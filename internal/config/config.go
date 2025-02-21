package config

import (
	"encoding/json"
	"errors"
	"io"

	"github.com/BurntSushi/toml"
)

type Config interface {
	FromEnviron() error
}

const (
	FromEnviron = iota
	FromJSON
	FromTOML
)

var ErrInvalidWhence = errors.New("invalid whence")

// Load populates config based on whence,
// if whence is a serialisation format, it is decoded from r.
func Load(config Config, whence int, r io.Reader) error {
	switch whence {
	case FromEnviron:
		return config.FromEnviron()

	case FromJSON:
		return json.NewDecoder(r).Decode(config)

	case FromTOML:
		_, err := toml.NewDecoder(r).Decode(config)
		return err

	default:
		return ErrInvalidWhence
	}
}
