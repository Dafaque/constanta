package model

import "encoding/json"

type Response struct {
	Data map[string]*[]byte
}

func (r *Response) Add(key string, value *[]byte) {
	r.Data[key] = value
}

func (r *Response) Bytes() []byte {
	b, err := json.Marshal(r)
	if err != nil {
		//? Could be only dev mistake, no need for client handle
		panic(err)
	}
	return b
}

func NewResponse() (resp *Response) {
	resp = &Response{}
	resp.Data = make(map[string]*[]byte)
	return resp
}
