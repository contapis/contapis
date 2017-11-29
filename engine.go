package main

import (
	"context"
	"fmt"
	"io"
	"os"
	"time"

	"github.com/cncd/pipeline/pipeline"
	"github.com/cncd/pipeline/pipeline/backend"
	"github.com/cncd/pipeline/pipeline/backend/docker"
	"github.com/cncd/pipeline/pipeline/interrupt"
	back "github.com/contapis/engine/backend"
	"github.com/contapis/engine/frontend"
)

// ExecuteJob runs a job
func executeJob(job *Job) {
	job.StartTime = time.Now()
	fmt.Println("starting execution of job", job.ID)
	var err error
	var reader io.ReadCloser
	var path = fmt.Sprintf("./%v.yml", job.Name)
	// if path == "-" {
	// 	reader = os.Stdin
	// } else {
	reader, err = os.Open(path)
	if err != nil {
		fmt.Println(err)
		job.Failed()
		return
	}
	// }
	defer reader.Close()

	pipelineConfig, err := frontend.Parse(reader)

	var engine backend.Engine
	engine, err = docker.NewEnv()
	if err != nil {
		fmt.Println("error", err)
		job.Failed()
		return
	}
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*30)
	defer cancel()
	ctx = interrupt.WithContext(ctx)
	config, err := back.Convert(pipelineConfig, job.Environment, job.Secrets)
	logger := back.NewPipelineLogger(config.Stages)
	err = pipeline.New(config,
		pipeline.WithContext(ctx),
		pipeline.WithLogger(logger),
		pipeline.WithTracer(defaultTracer),
		pipeline.WithEngine(engine),
	).Run()
	addLogs(job, logger)
	if err != nil {
		fmt.Println("job failed", job.ID)
		fmt.Println("-> error during pipe execution:\n", err)
		job.Failed()
	} else {
		fmt.Println("job succeeded", job.ID)
		job.Succeeded()
	}
	// for i, logs := range logger.Logs {
	// 	fmt.Println("logs from step", i, " are:\n", logs.String())
	// }
}

func addLogs(job *Job, collectedLogs *back.PipelineLogger) {
	for index, stage := range collectedLogs.Stages {
		job.Logs[stage.Name] = collectedLogs.Logs[index].String()
	}
}

var defaultTracer = pipeline.TraceFunc(func(state *pipeline.State) error {
	if state.Process.Exited {
		fmt.Printf("-> proc %q exited with status %d\n", state.Pipeline.Step.Name, state.Process.ExitCode)
		if state.Pipeline.Error != nil {
			fmt.Printf("error for proc %q: %q\n", state.Pipeline.Step.Name, state.Pipeline.Error)
		}
	} else {
		fmt.Printf("-> proc %q started\n", state.Pipeline.Step.Name)
		// state.Pipeline.Step.Environment["CI_BUILD_STATUS"] = "success"
		// state.Pipeline.Step.Environment["CI_BUILD_FINISHED"] = strconv.FormatInt(time.Now().Unix(), 10)
		// if state.Pipeline.Error != nil {
		// 	state.Pipeline.Step.Environment["CI_BUILD_STATUS"] = "failure"
		// }
	}
	return nil
})
