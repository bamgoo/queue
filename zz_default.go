package queue

import "github.com/infrago/infra"

// Register the built-in driver after the queue module has been mounted.
// The earlier registration in default.go is retained for compatibility with
// runtimes that apply registrations lazily.
func init() {
	infra.Register(infra.DEFAULT, &defaultDriver{})
}
