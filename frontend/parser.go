package frontend

import (
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"strings"

	"gopkg.in/yaml.v2"
)

// Config describes the entire pipe construction
type Config struct {
	Name        string
	Description string
	Pipeline    Steps `yaml:",flow"`
}

// Steps describes the set of containers needed to execute the pipeline
type Steps struct {
	Steps []*Step
}

// Step describes a single container execution
type Step struct {
	Name     string
	Image    string
	Commands []string `yaml:",flow"`
}

// UnmarshalYAML handles custom unmarshalling of YAML
func (p *Steps) UnmarshalYAML(unmarshal func(interface{}) error) error {
	slice := yaml.MapSlice{}
	if err := unmarshal(&slice); err != nil {
		return err
	}

	p.Steps = make([]*Step, 0)

	for _, s := range slice {
		step := Step{}
		out, _ := yaml.Marshal(s.Value)

		if err := yaml.Unmarshal(out, &step); err != nil {
			return err
		}
		if step.Name == "" {
			step.Name = fmt.Sprintf("%v", s.Key)
		}
		p.Steps = append(p.Steps, &step)
	}
	return nil
}

// Parse parses the pipeline config from an io.Reader.
func Parse(r io.Reader) (*Config, error) {
	pipeline := new(Config)
	b, err := ioutil.ReadAll(r)
	if err != nil {
		return nil, err
	}
	err = yaml.Unmarshal(b, &pipeline)
	if err != nil {
		return nil, err
	}
	return pipeline, nil
}

// ParseFile parses the pipeline config from a file.
func ParseFile(path string) (*Config, error) {
	f, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer f.Close()
	return Parse(f)
}

// ParseString parses the pipeline config from a string.
func ParseString(s string) (*Config, error) {
	return Parse(
		strings.NewReader(s),
	)
}
