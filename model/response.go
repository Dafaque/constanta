package model

import "encoding/json"

type Response struct {
	Data map[string]string
}

//? Add response data by URL as key
func (r *Response) Add(key string, value *[]byte) {
	//? Response might be not a JSON, the task is not clear here
	r.Data[key] = string(*value)
}

//? Check if URL already processed
func (r *Response) KeyExists(key string) bool {
	_, exists := r.Data[key]
	return exists
}

//? Serialize response
func (r *Response) Bytes() []byte {
	b, err := json.Marshal(r)
	if err != nil {
		//? Most likeley dev mistake, no need for client handle
		panic(err)
	}
	return b
}

//? Constructor
func NewResponse() (resp *Response) {
	resp = &Response{}
	resp.Data = make(map[string]string)
	return resp
}
