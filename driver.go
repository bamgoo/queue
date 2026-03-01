package queue

import (
	"time"

	base "github.com/infrago/base"
)

type (
	Driver interface {
		Connect(*Instance) (Connection, error)
	}

	Connection interface {
		Open() error
		Close() error
		Start() error
		Stop() error

		Register(name string) error
		Publish(name string, data []byte) error
		DeferredPublish(name string, data []byte, delay time.Duration) error
	}

	Instance struct {
		conn    Connection
		Name    string
		Config  Config
		Setting base.Map
	}
)
