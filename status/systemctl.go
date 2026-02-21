package status

import (
	"bytes"
	"fmt"
	"os/exec"
	"strings"
)

type SystemctlStatusManager struct {
	ServiceName string
}

func (manager *SystemctlStatusManager) Status() StatusResponse {
	cmd := exec.Command(
		"systemctl",
		"show",
		"-p",
		"ActiveState",
		"--value",
		manager.ServiceName,
	)
	stderr := new(bytes.Buffer)
	stdout := new(bytes.Buffer)
	cmd.Stderr = stderr
	cmd.Stdout = stdout
	err := cmd.Run()
	if err != nil {
		return StatusResponse{
			Status:      StatusError,
			Description: stderr.String(),
		}
	}
	if strings.TrimSpace(stdout.String()) == "active" {
		return StatusResponse{
			Status:      StatusRunning,
			Description: "",
		}
	} else {
		return StatusResponse{
			Status:      StatusStopped,
			Description: stdout.String(),
		}
	}
}

func (manager *SystemctlStatusManager) Start() ActionResult {
	cmd := exec.Command(
		"systemctl",
		"restart",
		fmt.Sprintf("%s.service", manager.ServiceName),
	)
	stderr := new(bytes.Buffer)
	stdout := new(bytes.Buffer)
	cmd.Stderr = stderr
	cmd.Stdout = stdout
	err := cmd.Run()
	if err != nil {
		return ActionResult{
			Success: false,
			Message: stderr.String(),
		}
	}

	return ActionResult{
		Success: true,
		Message: "Started service. " + stdout.String(),
	}
}

func (manager *SystemctlStatusManager) Stop() ActionResult {
	cmd := exec.Command(
		"systemctl",
		"stop",
		fmt.Sprintf("%s.service", manager.ServiceName),
	)
	stderr := new(bytes.Buffer)
	stdout := new(bytes.Buffer)
	cmd.Stderr = stderr
	cmd.Stdout = stdout
	err := cmd.Run()
	if err != nil {
		return ActionResult{
			Success: false,
			Message: stderr.String(),
		}
	}

	return ActionResult{
		Success: true,
		Message: "Stopped service. " + stdout.String(),
	}
}

func CreateSystemctlManager(data map[string]interface{}) (StatusManager, error) {
	service_name, ok := data["service_name"]
	if !ok {
		return nil, fmt.Errorf("Incorrect configuration. Could not find `service_name` value.")
	}
	service_name_str, ok := service_name.(string)
	if !ok {
		return nil, fmt.Errorf("`service_name` must be a string.")
	}
	manager := &SystemctlStatusManager{
		ServiceName: service_name_str,
	}
	return manager, nil
}
