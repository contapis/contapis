package main

import (
	"context"
	"encoding/base64"
	"fmt"
	"io"
	"os"
	"time"

	"github.com/cncd/pipeline/pipeline"
	"github.com/cncd/pipeline/pipeline/backend"
	"github.com/cncd/pipeline/pipeline/backend/docker"
	"github.com/cncd/pipeline/pipeline/interrupt"
	"github.com/cncd/pipeline/pipeline/multipart"
)

func main() {
	fmt.Println("hello world")

	engine, err := docker.NewEnv()
	if err != nil {
		fmt.Println("error", err)
	}
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*30)
	defer cancel()
	ctx = interrupt.WithContext(ctx)

	encodedCmds := base64.StdEncoding.EncodeToString([]byte("touch /data/first-file"))

	step0 := &backend.Step{
		Name:       "step0_test",
		Image:      "ubuntu:14.04",
		Entrypoint: []string{"/bin/sh", "-c"},
		Command:    []string{"echo " + encodedCmds + " | base64 -d | /bin/sh -e"},
		// Command: []string{"echo $CI_SCRIPT | base64 -d | /bin/sh -e"},
		// Environment: map[string]string{
		// 	"MY_SCRIPT": encodedCmds,
		// },
		Volumes:   []string{"data:/data"},
		OnSuccess: true,
	}

	encodedCmds = base64.StdEncoding.EncodeToString([]byte("ls -l /data/"))

	step1 := &backend.Step{
		Name:       "step1_test",
		Image:      "ubuntu:14.04",
		Entrypoint: []string{"/bin/sh", "-c"},
		Command:    []string{"echo " + encodedCmds + " | base64 -d | /bin/sh -e"},
		// Command: []string{"echo $CI_SCRIPT | base64 -d | /bin/sh -e"},
		// Environment: map[string]string{
		// 	"MY_SCRIPT": encodedCmds,
		// },
		Volumes:   []string{"data:/data"},
		OnSuccess: true,
	}

	stage0 := &backend.Stage{Name: "stage_name", Steps: []*backend.Step{step0}}

	config := backend.Config{
		// Stages:  []*backend.Stage{stage1},
		Stages:  []*backend.Stage{stage0, &backend.Stage{Name: "stage_name", Steps: []*backend.Step{step1}}},
		Volumes: []*backend.Volume{{Name: "data"}},
	}

	err = pipeline.New(&config,
		pipeline.WithContext(ctx),
		pipeline.WithLogger(defaultLogger),
		pipeline.WithTracer(defaultTracer),
		pipeline.WithEngine(engine),
	).Run()
	if err != nil {
		fmt.Println("error", err)
	}
}

var defaultLogger = pipeline.LogFunc(func(proc *backend.Step, rc multipart.Reader) error {
	part, err := rc.NextPart()
	if err != nil {
		return err
	}
	io.Copy(os.Stderr, part)
	return nil
})

var defaultTracer = pipeline.TraceFunc(func(state *pipeline.State) error {
	if state.Process.Exited {
		fmt.Printf("proc %q exited with status %d\n", state.Pipeline.Step.Name, state.Process.ExitCode)
	} else {
		fmt.Printf("proc %q started\n", state.Pipeline.Step.Name)
		// state.Pipeline.Step.Environment["CI_BUILD_STATUS"] = "success"
		// state.Pipeline.Step.Environment["CI_BUILD_FINISHED"] = strconv.FormatInt(time.Now().Unix(), 10)
		// if state.Pipeline.Error != nil {
		// 	state.Pipeline.Step.Environment["CI_BUILD_STATUS"] = "failure"
		// }
	}
	return nil
})
