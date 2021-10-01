package model

import (
	"encoding/json"
	"errors"
	"io/ioutil"
	"net/http"
	"sync"
	"time"
)

//? Expected request
type Request struct {
	Urls []string `json:"urls"`
}

func (r *Request) Process() (*Response, error) {
	var resp *Response = NewResponse()
	var client *http.Client = &http.Client{
		Timeout: 1 * time.Second,
	}

	var err chan error = make(chan error)
	var done chan bool = make(chan bool)

	var wg sync.WaitGroup

	wg.Add(len(r.Urls))
	//? Waitgroup routine
	go func() {
		//? Request loop
		for _, u := range r.Urls {
			var urlString = u
			go func() {
				httpReq, errMakeReq := http.NewRequest("GET", urlString, nil)
				if errMakeReq != nil {
					err <- errMakeReq
					return
				}
				response, errSend := client.Do(httpReq)
				if errSend != nil {
					err <- errSend
					return
				}
				switch response.StatusCode {
				case http.StatusOK, http.StatusAccepted, http.StatusCreated:
					data, errReadBody := ioutil.ReadAll(response.Body)
					if errReadBody != nil {
						err <- errReadBody
						return
					}
					response.Body.Close()
					resp.Add(urlString, &data)
					wg.Done()
				default:
					err <- errors.New(response.Status)
				}
			}()
		}
		wg.Wait()
		done <- true
	}()

	select {
	case e := <-err:
		return nil, e
	case <-done:
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
