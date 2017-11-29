package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"

	"github.com/go-chi/chi"
	"github.com/go-chi/render"
)

func apiInit(newJobs chan *Job, catalog *JobCatalog) {
	r := chi.NewRouter()
	r.Route("/jobs", func(r chi.Router) {
		r.Post("/", newJob(catalog)) // POST /jobs
		r.Get("/{ID}/status", jobStatus(catalog))
		r.Get("/{ID}/wait", jobWait(catalog))
	})
	err := http.ListenAndServe("127.0.0.1:3000", r)
	if err != nil {
		fmt.Println("error serving requests:", err)
	}
}

func newJob(catalog *JobCatalog) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		jobDescription := JobDescription{}
		if err := decodeBody(r.Body, &jobDescription); err != nil {
			render.Render(w, r, ErrInvalidRequest(err))
			return
		}
		job := BuildJobFromDescription(jobDescription)
		fmt.Println("scheduling job:", job.ID)
		catalog.Add(job)
		w.Write([]byte(fmt.Sprintf("{\"ID\": \"%v\"}", job.ID)))
	}
}

func jobStatus(catalog *JobCatalog) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		jobID := chi.URLParam(r, "ID")
		if jobID == "" {
			render.Render(w, r, ErrInvalidRequest(errors.New("ID not specified")))
			return
		}
		if job := catalog.Get(jobID); job != nil {
			render.Render(w, r, job)
			return
		}
		render.Render(w, r, ErrNotFound)
	}
}

func jobWait(catalog *JobCatalog) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		jobID := chi.URLParam(r, "ID")
		if jobID == "" {
			render.Render(w, r, ErrInvalidRequest(errors.New("ID not specified")))
			return
		}
		if job := catalog.Get(jobID); job != nil {
			fmt.Println("ptr:", job)
			job.Wait()
			w.WriteHeader(200)
			w.Write([]byte{})
			return
		}
		render.Render(w, r, ErrNotFound)
	}
}

func (job *Job) Render(w http.ResponseWriter, r *http.Request) error {
	// Pre-processing before a response is marshalled and sent across the wire
	return nil
}

type JobDescription struct {
	Name        string
	Environment map[string]string
	Secrets     map[string]string
}

func decodeBody(r io.Reader, v interface{}) error {
	var (
		err      error
		respBody []byte
	)
	if respBody, err = ioutil.ReadAll(r); err != nil {
		return err
	}
	if err := json.Unmarshal(respBody, &v); err != nil {
		return err
	}
	return nil
}

func ErrInvalidRequest(err error) render.Renderer {
	return &ErrResponse{
		Err:            err,
		HTTPStatusCode: 400,
		StatusText:     "Invalid request.",
		ErrorText:      err.Error(),
	}
}

var ErrNotFound = &ErrResponse{HTTPStatusCode: 404, StatusText: "Resource not found."}

func (e *ErrResponse) Render(w http.ResponseWriter, r *http.Request) error {
	render.Status(r, e.HTTPStatusCode)
	return nil
}

// ErrResponse renderer type for handling all sorts of errors.
//
// In the best case scenario, the excellent github.com/pkg/errors package
// helps reveal information on the error, setting it on Err, and in the Render()
// method, using it to set the application-specific error code in AppCode.
type ErrResponse struct {
	Err            error `json:"-"` // low-level runtime error
	HTTPStatusCode int   `json:"-"` // http response status code

	StatusText string `json:"status"`          // user-level status message
	AppCode    int64  `json:"code,omitempty"`  // application-specific error code
	ErrorText  string `json:"error,omitempty"` // application-level error message, for debugging
}
