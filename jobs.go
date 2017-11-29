package main

import (
	"encoding/json"
	"fmt"
	"sync"
	"time"

	"github.com/rs/xid"
)

// Job describes a submitted job
type Job struct {
	ID          string
	Name        string
	Environment map[string]string
	Secrets     map[string]string
	SubmitTime  time.Time
	StartTime   time.Time
	EndTime     time.Time
	Logs        map[string]string
	// Repo        string
	// Branch      string
	// Path        string
	doneChans  []chan bool
	done       bool
	successful bool
	mutex      sync.Mutex
}

func (j *Job) IsDone() bool {
	return j.done
}

func (j *Job) IsSuccessful() bool {
	return j.successful
}

func BuildJobFromDescription(desc JobDescription) *Job {
	job := Job{
		ID:          xid.New().String(),
		Name:        desc.Name,
		Environment: desc.Environment,
		Secrets:     desc.Secrets,
		Logs:        map[string]string{},
		SubmitTime:  time.Now(),
	}
	job.doneChans = make([]chan bool, 0)
	job.mutex = sync.Mutex{}
	return &job
}

func (job *Job) Wait() {
	if !job.done {
		job.mutex.Lock()
		if !job.done {
			ch := make(chan bool)
			job.doneChans = append(job.doneChans, ch)
			job.mutex.Unlock()
			_ = <-ch
		}
	}
}

func (job *Job) Failed() {
	if job.done {
		fmt.Println("called Failed() after already marked as done")
		return
	}
	job.mutex.Lock()
	job.EndTime = time.Now()
	job.successful = false
	job.done = true
	for _, ch := range job.doneChans {
		ch <- true
	}
	job.mutex.Unlock()
}

func (job *Job) Succeeded() {
	if job.done {
		fmt.Println("called Failed() after already marked as done")
		return
	}
	job.mutex.Lock()
	job.EndTime = time.Now()
	job.successful = true
	job.done = true
	for _, ch := range job.doneChans {
		ch <- true
	}
	job.mutex.Unlock()
}

type jobExternalRep struct {
	ID           string            `json:"id"`
	Name         string            `json:"name"`
	Environment  map[string]string `json:"environment"`
	SubmitTime   time.Time         `json:"submit_time"`
	StartTime    time.Time         `json:"start_time"`
	EndTime      time.Time         `json:"end_time"`
	Logs         map[string]string `json:"logs"`
	IsDone       bool              `json:"is_done"`
	IsSuccessful bool              `json:"is_successful"`
}

func NewExternalRep(j *Job) jobExternalRep {
	return jobExternalRep{
		ID:           j.ID,
		Name:         j.Name,
		Environment:  j.Environment,
		SubmitTime:   j.SubmitTime,
		StartTime:    j.StartTime,
		EndTime:      j.EndTime,
		Logs:         j.Logs,
		IsDone:       j.IsDone(),
		IsSuccessful: j.IsSuccessful(),
	}
}

func (job *Job) MarshalJSON() ([]byte, error) {
	j := NewExternalRep(job)
	return json.Marshal(j)
}

var newJobs = make(chan *Job, 10)

type JobCatalog struct {
	jobs    map[string]*Job
	newJobs []chan *Job
}

func NewCatalog() *JobCatalog {
	catalog := &JobCatalog{
		jobs:    make(map[string]*Job),
		newJobs: make([]chan *Job, 0),
	}
	return catalog
}

func (c *JobCatalog) Add(job *Job) {
	fmt.Println("adding job", job.ID)
	c.jobs[job.ID] = job
	for _, ch := range c.newJobs {
		ch <- job
	}
}

func (c *JobCatalog) Get(ID string) *Job {
	if job, ok := c.jobs[ID]; ok {
		return job
	}
	return nil
}

func (c *JobCatalog) RegisterNewJobNotification(jobNotification chan *Job) {
	c.newJobs = append(c.newJobs, jobNotification)
}
