package status

import "fmt"

// StatusType enum.
type StatusType int

const (
	StatusRunning StatusType = iota
	StatusStopped
	StatusError
)

var statusName = map[StatusType]string{
	StatusRunning: "Running",
	StatusStopped: "Stopped",
	StatusError:   "Error",
}

func (st StatusType) String() string {
	return statusName[st]
}

// A struct that is returned to tell the status of a monitored service.
type StatusResponse struct {
	Status      StatusType
	Description string
}

// A struct that is returned to describe the result of performing an action.
type ActionResult struct {
	Success bool
	Message string
}

// The interface that all statuses must implement.
type StatusManager interface {
	Status() StatusResponse
	Start() ActionResult
	Stop() ActionResult
}

// The type of the callable that must be registered in the StatusManagerFactory.
type StatusManagerFactoryCallable func(map[string]interface{}) (StatusManager, error)

// A struct to allow registering various StatusManager types that can be created.
type StatusManagerFactory struct {
	registry map[string]StatusManagerFactoryCallable
}

func NewStatusManagerFactory() *StatusManagerFactory {
	data := &StatusManagerFactory{registry: make(map[string]StatusManagerFactoryCallable)}
	return data
}

// Register the creating function in the StatusManagerFactory instance.
func (factory *StatusManagerFactory) Register(name string, creator StatusManagerFactoryCallable) error {
	_, ok := factory.registry[name]
	if ok {
		return fmt.Errorf("`%s` already exists in the registry and it cannot be registered twice.", name)
	}
	factory.registry[name] = creator
	return nil
}

// Create an instance of a StatusManager.
func (factory *StatusManagerFactory) Create(name string, data map[string]interface{}) (StatusManager, error) {
	creator, ok := factory.registry[name]
	if !ok {
		return nil, fmt.Errorf("`%s` not registered in the StatusManagerFactory.", name)
	}
	return creator(data)
}
