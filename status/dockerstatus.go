package status

import (
	"bytes"
	"fmt"
	"os/exec"
)

type DockerStatusManager struct {
	ContainerName string
	Image         string
	Arguments     []string
}

func (manager *DockerStatusManager) Status() StatusResponse {
	cmd := exec.Command(
		"docker",
		"container",
		"inspect",
		"-f",
		"{{.State.Status}}",
		manager.ContainerName,
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
	if stdout.String() == "running" {
		return StatusResponse{
			Status:      StatusRunning,
			Description: "",
		}
	} else {
		return StatusResponse{
			Status:      StatusStopped,
			Description: "",
		}
	}
}

func (manager *DockerStatusManager) Start() ActionResult {
	arguments := []string{"--name", manager.ContainerName}
	arguments = append(arguments, manager.Arguments...)
	arguments = append(arguments, manager.Image)
	cmd := exec.Command("docker", arguments...)

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
		Message: "Started docker container. " + stdout.String(),
	}
}

func (manager *DockerStatusManager) Stop() ActionResult {
	cmd := exec.Command("docker", "stop", manager.ContainerName)

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
		Message: "Stopped docker container." + stdout.String(),
	}
}

func CreateDockerManager(data map[string]interface{}) (StatusManager, error) {
	container_name, ok := data["container_name"]
	if !ok {
		return nil, fmt.Errorf("Incorrect configuration. Could not find `container_name` value.")
	}
	container_name_str, ok := container_name.(string)
	if !ok {
		return nil, fmt.Errorf("`container_name` must be a string.")
	}
	image, ok := data["image"]
	if !ok {
		return nil, fmt.Errorf("Incorrect configuration. Could not find `image` value.")
	}
	image_str, ok := image.(string)
	if !ok {
		return nil, fmt.Errorf("`image` must be a string.")
	}
	arguments, ok := data["arguments"]
	if !ok {
		return nil, fmt.Errorf("Incorrect configuration. Could not find `arguments` value.")
	}
	arguments_list, ok := arguments.([]interface{})
	if !ok {
		return nil, fmt.Errorf("`arguments` must be a list.")
	}
	arguments_list_str := make([]string, len(arguments_list))
	for i := 0; i < len(arguments_list); i++ {
		arguments_list_str[i], ok = arguments_list[i].(string)
		if !ok {
			return nil, fmt.Errorf("Element %d of the arguments list is not a string.", i)
		}
	}

	manager := &DockerStatusManager{
		ContainerName: container_name_str,
		Image:         image_str,
		Arguments:     arguments_list_str,
	}
	return manager, nil
}
