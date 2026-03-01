package queue

import (
	"errors"
	"sync"
	"time"

	"github.com/infrago/infra"
)

func init() {
	infra.Register(infra.DEFAULT, &defaultDriver{})
}

var (
	errQueueRunning    = errors.New("queue is running")
	errQueueNotRunning = errors.New("queue is not running")
)

type (
	defaultDriver struct{}

	defaultConnection struct {
		mutex    sync.RWMutex
		running  bool
		instance *Instance
		queues   map[string]chan *defaultMessage
		done     chan struct{}
		wg       sync.WaitGroup
	}

	defaultMessage struct {
		name    string
		data    []byte
		attempt int
	}
)

func (d *defaultDriver) Connect(inst *Instance) (Connection, error) {
	return &defaultConnection{
		instance: inst,
		queues:   make(map[string]chan *defaultMessage, 0),
		done:     make(chan struct{}),
	}, nil
}

func (c *defaultConnection) Open() error  { return nil }
func (c *defaultConnection) Close() error { return nil }

func (c *defaultConnection) Register(name string) error {
	c.mutex.Lock()
	defer c.mutex.Unlock()
	if _, ok := c.queues[name]; !ok {
		c.queues[name] = make(chan *defaultMessage, 256)
	}
	return nil
}

func (c *defaultConnection) Start() error {
	c.mutex.Lock()
	defer c.mutex.Unlock()

	if c.running {
		return errQueueRunning
	}
	for _, ch := range c.queues {
		queueCh := ch
		c.wg.Add(1)
		go func() {
			defer c.wg.Done()
			for {
				select {
				case msg := <-queueCh:
					req := Request{
						Name:      msg.name,
						Data:      msg.data,
						Attempt:   msg.attempt,
						Timestamp: time.Now(),
					}
					res := c.instance.Serve(req)
					if res.Retry {
						msg.attempt++
						c.publish(msg, res.Delay)
					}
				case <-c.done:
					return
				}
			}
		}()
	}
	c.running = true
	return nil
}

func (c *defaultConnection) Stop() error {
	c.mutex.Lock()
	defer c.mutex.Unlock()

	if !c.running {
		return errQueueNotRunning
	}
	close(c.done)
	c.wg.Wait()
	c.done = make(chan struct{})
	c.running = false
	return nil
}

func (c *defaultConnection) publish(msg *defaultMessage, delay time.Duration) {
	c.mutex.RLock()
	ch := c.queues[msg.name]
	c.mutex.RUnlock()
	if ch == nil {
		return
	}
	if delay > 0 {
		time.AfterFunc(delay, func() {
			ch <- msg
		})
		return
	}
	ch <- msg
}

func (c *defaultConnection) Publish(name string, data []byte) error {
	c.publish(&defaultMessage{name: name, data: data, attempt: 1}, 0)
	return nil
}

func (c *defaultConnection) DeferredPublish(name string, data []byte, delay time.Duration) error {
	c.publish(&defaultMessage{name: name, data: data, attempt: 1}, delay)
	return nil
}
