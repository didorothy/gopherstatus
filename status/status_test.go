package status

import (
	"reflect"
	"testing"
)

// Ensure that we can create the varius status types and that their String()
// method returns the expected value.
func TestStatusType(t *testing.T) {
	if StatusRunning.String() != "Running" {
		t.Errorf("StatusRunning did not return the text 'Running'")
	}
	if StatusStopped.String() != "Stopped" {
		t.Errorf("StatusStopped did not return the text 'Stopped'")
	}
	if StatusError.String() != "Error" {
		t.Errorf("StatusError did not return the text 'Error'")
	}
}

type FakeStatusManager struct{}

func (manager *FakeStatusManager) Status() StatusResponse {
	return StatusResponse{
		Status:      StatusRunning,
		Description: "Fake working.",
	}
}

func (manager *FakeStatusManager) Start() ActionResult {
	return ActionResult{
		Success: true,
		Message: "Fake started.",
	}
}

func (manager *FakeStatusManager) Stop() ActionResult {
	return ActionResult{
		Success: true,
		Message: "Fake stopped.",
	}
}

// Ensure that we can make a new factory, register a callable to create our
// FakeStatusManager, and create an instance of the FakeStatusManager.
func TestStatusManagerFactory(t *testing.T) {
	factory := NewStatusManagerFactory()
	err := factory.Register("fake", func(data map[string]interface{}) (StatusManager, error) {
		manager := &FakeStatusManager{}
		return manager, nil
	})
	if err != nil {
		t.Errorf("Failed to register FakeStatusManager as 'fake': %s", err)
	}
	manager, err := factory.Create("fake", make(map[string]interface{}))
	if err != nil {
		t.Errorf("Failed to create an instance of the FakeStatusManager: %s", err)
	}
	manager_type := reflect.TypeOf(manager).String()
	if manager_type != "*status.FakeStatusManager" {
		t.Errorf("Got unexpected type of StatusManager interface returned: %s", manager_type)
	}
}

// Ensure that we cannot register a name twice.
func TestStatusManagerFactoryDuplicateRegistration(t *testing.T) {
	factory := NewStatusManagerFactory()
	err := factory.Register("fake", func(data map[string]interface{}) (StatusManager, error) {
		manager := &FakeStatusManager{}
		return manager, nil
	})
	if err != nil {
		t.Errorf("Failed to register FakeStatusManager as 'fake' the first time: %s", err)
	}
	err = factory.Register("fake", func(data map[string]interface{}) (StatusManager, error) {
		manager := &FakeStatusManager{}
		return manager, nil
	})
	if err == nil {
		t.Errorf("Expected error when registering the fake function a second time but did not recieve one.")
	}
}

// Ensure that when a string is not registered that an error is returned.
func TestStatusManagerFactoryNotRegistered(t *testing.T) {
	factory := NewStatusManagerFactory()
	_, err := factory.Create("fake", make(map[string]interface{}))
	if err == nil {
		t.Errorf("Expected error when trying to create unregistered StatusManager.")
	}
}
