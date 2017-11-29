package backend

import (
	"bytes"
	"fmt"
	"io"
	"io/ioutil"

	"github.com/cncd/pipeline/pipeline/backend"
	"github.com/cncd/pipeline/pipeline/multipart"
)

type PipelineLogger struct {
	Stages []*backend.Stage
	Logs   []bytes.Buffer
}

func NewPipelineLogger(stages []*backend.Stage) *PipelineLogger {
	p := new(PipelineLogger)
	p.Stages = stages
	p.Logs = make([]bytes.Buffer, len(p.Stages))
	return p
}

func (p PipelineLogger) Log(step *backend.Step, reader multipart.Reader) error {
	var stepNumber int
	for i, s := range p.Stages {
		if s.Steps[0] == step {
			stepNumber = i
			break
		}
	}
	for {
		part, err := reader.NextPart()
		if err == io.EOF {
			return nil
		}
		if err != nil {
			return err
		}
		_, err = copier(&p.Logs[stepNumber], part)
		if err != nil {
			return err
		}
	}
}

func copier(b *bytes.Buffer, r io.Reader) (n int64, err error) {
	allBytes, err := ioutil.ReadAll(r)
	// fmt.Println("read :", len(allBytes), err)
	c, e := b.Write(allBytes)
	// fmt.Println("write :", c, e)
	return int64(c), e
}

func copier2(b *bytes.Buffer, r io.Reader) (n int64, err error) {
	var total int64
	var buf = make([]byte, 0, 512)
	for {
		fmt.Println("buffer at (len, capacity):", len(buf), cap(buf))
		c, e := r.Read(buf)
		fmt.Println("read OP: (count, error):", c, e)
		if c == 0 {
			fmt.Println("no data to copy (read is zero)")
			return 0, nil
		}
		if c > 0 {
			total = total + int64(c)
			fmt.Println("read total bytes:", total)
			c2, e2 := b.Write(buf)
			fmt.Println("result of buffer write (count, error): ", c2, e2)
			if e2 != nil {
				fmt.Println("WHHHHHAT", e2)
				fmt.Println(c2)
				return total, e
			}
		}
		if e == io.EOF {
			fmt.Println("done copying... due to EOF")
			return total, nil
		}
		if e != nil {
			fmt.Println("ERR WHILE reading", e)
			fmt.Println(c)
			return total, e
		}
	}
}
