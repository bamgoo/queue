package queue

func (Queue) RegistryComponent() string {
	return "queue.queue"
}

func (Declare) RegistryComponent() string {
	return "queue.declare"
}

func (Filter) RegistryComponent() string {
	return "queue.filter"
}

func (Handler) RegistryComponent() string {
	return "queue.handler"
}
