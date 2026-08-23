package queue

import (
	"github.com/infrago/base"
	"github.com/infrago/infra"
)

func (m *Module) Ready() bool {
	m.mutex.RLock()
	defer m.mutex.RUnlock()
	return m.started && len(m.instances) > 0
}

func (m *Module) Health() infra.ModuleHealth {
	m.mutex.RLock()
	started := m.started
	connections := len(m.instances)
	queues := len(m.queues)
	m.mutex.RUnlock()
	return infra.NewModuleHealth("queue", started && connections > 0, nil, base.Map{
		"connections": connections,
		"queues":      queues,
	})
}

func (m *Module) Stats() infra.ModuleStats {
	m.mutex.RLock()
	started := m.started
	connections := len(m.instances)
	queues := len(m.queues)
	declares := len(m.declares)
	m.mutex.RUnlock()
	return infra.NewModuleStats("queue", started && connections > 0, base.Map{
		"connections": connections,
		"queues":      queues,
		"declares":    declares,
	})
}
