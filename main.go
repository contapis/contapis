package main

func main() {

	catalog := NewCatalog()
	executionChan := make(chan *Job, 10)
	catalog.RegisterNewJobNotification(executionChan)
	go executor(executionChan)
	apiInit(newJobs, catalog)
	return

	// var err error
	// var reader io.ReadCloser
	// var path = "./test.yml"
	// if path == "-" {
	// 	reader = os.Stdin
	// } else {
	// 	reader, err = os.Open(path)
	// 	if err != nil {
	// 		fmt.Println(err)
	// 	}
	// }
	// defer reader.Close()

	// pipelineConfig, err := frontend.Parse(reader)

	// var engine backend.Engine
	// engine, err = docker.NewEnv()
	// if err != nil {
	// 	fmt.Println("error", err)
	// }
	// ctx, cancel := context.WithTimeout(context.Background(), time.Second*30)
	// defer cancel()
	// ctx = interrupt.WithContext(ctx)

	// // encodedCmds := base64.StdEncoding.EncodeToString([]byte("touch /data/first-file"))

	// // step0 := &cncdbackend.Step{
	// // 	Name:       "step0_test",
	// // 	Image:      "ubuntu:14.04",
	// // 	Entrypoint: []string{"/bin/sh", "-c"},
	// // 	Command:    []string{"echo " + encodedCmds + " | base64 -d | /bin/sh -e"},
	// // 	// Command: []string{"echo $CI_SCRIPT | base64 -d | /bin/sh -e"},
	// // 	// Environment: map[string]string{
	// // 	// 	"MY_SCRIPT": encodedCmds,
	// // 	// },
	// // 	Volumes:   []string{"data:/data"},
	// // 	OnSuccess: true,
	// // }

	// // encodedCmds = base64.StdEncoding.EncodeToString([]byte("ls -l /data/"))

	// // step1 := &cncdbackend.Step{
	// // 	Name:       "step1_test",
	// // 	Image:      "ubuntu:14.04",
	// // 	Entrypoint: []string{"/bin/sh", "-c"},
	// // 	Command:    []string{"echo " + encodedCmds + " | base64 -d | /bin/sh -e"},
	// // 	// Command: []string{"echo $CI_SCRIPT | base64 -d | /bin/sh -e"},
	// // 	// Environment: map[string]string{
	// // 	// 	"MY_SCRIPT": encodedCmds,
	// // 	// },
	// // 	Volumes:   []string{"data:/data"},
	// // 	OnSuccess: true,
	// // }

	// // stage0 := &backend.Stage{Name: "stage_name", Steps: []*backend.Step{step0}}

	// // config := &backend.Config{
	// // 	// Stages:  []*backend.Stage{stage1},
	// // 	Stages:  []*backend.Stage{stage0, &backend.Stage{Name: "stage_name", Steps: []*backend.Step{step1}}},
	// // 	Volumes: []*backend.Volume{{Name: "data"}},
	// // }

	// config, err := back.Convert(pipelineConfig, )
	// logger := back.NewPipelineLogger(config.Stages)
	// err = pipeline.New(config,
	// 	pipeline.WithContext(ctx),
	// 	pipeline.WithLogger(logger),
	// 	// pipeline.WithLogger(defaultLogger),
	// 	pipeline.WithTracer(defaultTracer),
	// 	pipeline.WithEngine(engine),
	// ).Run()
	// if err != nil {
	// 	fmt.Println("-> error during pipe execution:\n", err)
	// }
	// for i, logs := range logger.Logs {
	// 	fmt.Println("logs from step", i, " are:\n", logs.String())
	// }
}

// var defaultTracer = pipeline.TraceFunc(func(state *pipeline.State) error {
// 	if state.Process.Exited {
// 		fmt.Printf("-> proc %q exited with status %d\n", state.Pipeline.Step.Name, state.Process.ExitCode)
// 		if state.Pipeline.Error != nil {
// 			fmt.Printf("error for proc %q: %q\n", state.Pipeline.Step.Name, state.Pipeline.Error)
// 		}
// 	} else {
// 		fmt.Printf("-> proc %q started\n", state.Pipeline.Step.Name)
// 		// state.Pipeline.Step.Environment["CI_BUILD_STATUS"] = "success"
// 		// state.Pipeline.Step.Environment["CI_BUILD_FINISHED"] = strconv.FormatInt(time.Now().Unix(), 10)
// 		// if state.Pipeline.Error != nil {
// 		// 	state.Pipeline.Step.Environment["CI_BUILD_STATUS"] = "failure"
// 		// }
// 	}
// 	return nil
// })

func executor(jobsChan chan *Job) {
	for {
		job := <-jobsChan
		executeJob(job)
	}
}
