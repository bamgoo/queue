package queue

import (
	"time"

	. "github.com/bamgoo/base"
)

func Publish(name string, values ...Map) error {
	var value Map
	if len(values) > 0 {
		value = values[0]
	}
	return module.publish("", name, value)
}

func PublishTo(conn, name string, values ...Map) error {
	var value Map
	if len(values) > 0 {
		value = values[0]
	}
	return module.publish(conn, name, value)
}

func DeferredPublish(name string, value Map, delay time.Duration) error {
	return module.publish("", name, value, delay)
}

func DeferredPublishTo(conn, name string, value Map, delay time.Duration) error {
	return module.publish(conn, name, value, delay)
}

func RegisterDriver(name string, driver Driver) {
	module.RegisterDriver(name, driver)
}

func RegisterConfig(name string, cfg Config) {
	module.RegisterConfig(name, cfg)
}

func RegisterConfigs(cfgs Configs) {
	module.RegisterConfigs(cfgs)
}
