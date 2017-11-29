package backend

import (
	"encoding/base64"
	"fmt"
	"strings"

	"github.com/cncd/pipeline/pipeline/backend"
	"github.com/contapis/engine/frontend"
)

// Convert maps a configuration pipeline to an actual implementation of the pipeline
func Convert(cfg *frontend.Config,
	envOverride map[string]string,
	secrets map[string]string) (*backend.Config, error) {
	var err error
	config := new(backend.Config)
	config.Volumes = []*backend.Volume{{Name: "data"}}
	config.Stages, err = buildStages(cfg.Pipeline, envOverride, secrets)
	if err != nil {
		return nil, err
	}

	return config, nil
}

func buildStages(steps frontend.Steps,
	envOverride map[string]string,
	secrets map[string]string) ([]*backend.Stage, error) {
	stages := make([]*backend.Stage, 0)
	for _, step := range steps.Steps {
		stage := backend.Stage{Name: step.Name}
		backendStep, err := buildStep(step, envOverride, secrets)
		if err != nil {
			return nil, err
		}
		stage.Steps = []*backend.Step{backendStep}
		stages = append(stages, &stage)
	}

	return stages, nil
}

func buildStep(step *frontend.Step,
	envOverride map[string]string,
	secrets map[string]string) (*backend.Step, error) {
	ret := backend.Step{}
	ret.Name = step.Name
	ret.Image = step.Image
	encodedCommand := encodeCommand(step.Commands)
	ret.Command = []string{encodedCommand}
	ret.Entrypoint = []string{"/bin/sh", "-c"}
	ret.Volumes = []string{"data:/data"}
	ret.WorkingDir = "/data"
	ret.Pull = true
	ret.OnSuccess = true

	ret.Environment = map[string]string{}
	for k, v := range envOverride {
		ret.Environment[k] = v
	}
	for secretKey, secretValue := range secrets {
		ret.Environment[fmt.Sprintf("SECRET_%v", secretKey)] = secretValue
	}

	return &ret, nil
}

func encodeCommand(commands []string) string {
	encodedCmds := strings.Join(commands, "\n")
	return "echo " + base64.StdEncoding.EncodeToString([]byte(encodedCmds)) + " | base64 -d | /bin/sh -e"
}

// encodedCmds := base64.StdEncoding.EncodeToString([]byte("touch /data/first-file"))

// step0 := &cncdbackend.Step{
// 	Name:       "step0_test",
// 	Image:      "ubuntu:14.04",
// 	Entrypoint: []string{"/bin/sh", "-c"},
// 	Command:    []string{"echo " + encodedCmds + " | base64 -d | /bin/sh -e"},
// 	// Command: []string{"echo $CI_SCRIPT | base64 -d | /bin/sh -e"},
// 	// Environment: map[string]string{
// 	// 	"MY_SCRIPT": encodedCmds,
// 	// },
// 	Volumes:   []string{"data:/data"},
// 	OnSuccess: true,
// }

// encodedCmds = base64.StdEncoding.EncodeToString([]byte("ls -l /data/"))

// step1 := &cncdbackend.Step{
// 	Name:       "step1_test",
// 	Image:      "ubuntu:14.04",
// 	Entrypoint: []string{"/bin/sh", "-c"},
// 	Command:    []string{"echo " + encodedCmds + " | base64 -d | /bin/sh -e"},
// 	// Command: []string{"echo $CI_SCRIPT | base64 -d | /bin/sh -e"},
// 	// Environment: map[string]string{
// 	// 	"MY_SCRIPT": encodedCmds,
// 	// },
// 	Volumes:   []string{"data:/data"},
// 	OnSuccess: true,
// }

// stage0 := &backend.Stage{Name: "stage_name", Steps: []*backend.Step{step0}}

// config := &backend.Config{
// 	// Stages:  []*backend.Stage{stage1},
// 	Stages:  []*backend.Stage{stage0, &backend.Stage{Name: "stage_name", Steps: []*backend.Step{step1}}},
// 	Volumes: []*backend.Volume{{Name: "data"}},
// }
