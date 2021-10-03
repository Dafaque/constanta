package model

import (
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"sync"
	"time"
)

//? Expected request
type Request struct {
	Urls []string `json:"urls"`
}

func (r *Request) Process(httpRequestContext <-chan struct{}) (*Response, error) {
	var resp *Response = NewResponse()

	var client *http.Client = &http.Client{
		//? Outgoing requests timeout
		Timeout: 1 * time.Second,
	}

	//? Requests context & cancel func to break processing
	ctx, cancel := context.WithCancel(context.Background())

	//? Listening for http request's cancel to break process
	go func() {
		<-httpRequestContext
		cancel()
	}()

	var err chan error = make(chan error)
	var done chan bool = make(chan bool)

	var groupLen int = len(r.Urls)
	var groupDone int = 0
	//? For sync `groupDone` increment
	var mu sync.Mutex

	go func() {
		//? Request loop
		for _, u := range r.Urls {
			//? Anti strong capture
			var urlString = u
			go func() {
				//? If ctx were cancelled reqs will be instanced with err
				httpReq, errMakeReq := http.NewRequestWithContext(ctx, "GET", urlString, nil)
				if errMakeReq != nil {
					err <- errMakeReq
					return
				}
				//? Sending the request
				response, errSend := client.Do(httpReq)
				if errSend != nil {
					err <- errSend
					return
				}
				switch response.StatusCode {
				//? Some of 2xx statuses
				case http.StatusOK, http.StatusAccepted, http.StatusCreated:
					data, errReadBody := ioutil.ReadAll(response.Body)
					if errReadBody != nil {
						err <- errReadBody
						return
					}
					//? Reading response
					response.Body.Close()
					resp.Add(urlString, &data)
					//? Incrementing counter
					mu.Lock()
					groupDone++
					mu.Unlock()
					if groupDone == groupLen {
						done <- true
					}
					return
				default:
					err <- fmt.Errorf("bad response at %s: %s", urlString, response.Status)
					return
				}
			}()
		}
	}()

	select {
	//? Awaiting for error..
	case e := <-err:
		cancel()
		return nil, e
	//? ..or success
	case <-done:
		cancel()
		return resp, nil
	}
}

//? Constructor
func NewRequest(req *http.Request) (*Request, error) {
	defer req.Body.Close()
	var r *Request = &Request{}
	data, errReadBody := ioutil.ReadAll(req.Body)
	if errReadBody != nil {
		return nil, errReadBody
	}
	errUnmarshalBody := json.Unmarshal(data, r)
	return r, errUnmarshalBody
}
