package handlers

import (
	"net/http"

	connectioncontext "github.com/Dafaque/constanta/connection_context"
	"github.com/Dafaque/constanta/model"
)

//? Max allowed URLs in payload
const MAX_URLS int = 20

func Root(w http.ResponseWriter, r *http.Request) {
	//? 429 check
	var connectionsCount = r.Context().Value(connectioncontext.ConnectionsKeyCounter).(int)
	if connectionsCount > connectioncontext.MAX_CONNECTIONS {
		w.WriteHeader(http.StatusTooManyRequests)
		return
	}
	//? POST Only
	if r.Method != http.MethodPost {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}

	//? Read req body into struct
	req, errParseRequest := model.NewRequest(r)
	if errParseRequest != nil {
		w.WriteHeader(http.StatusUnprocessableEntity)
		return
	}
	//? To many URLs
	if len(req.Urls) > MAX_URLS {
		w.WriteHeader(http.StatusRequestEntityTooLarge)
		return
	}

	//? Main stuff here
	resp, err := req.Process(r.Context().Done())

	//? Something went wrong
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(err.Error()))
		return
	}

	w.Write(resp.Bytes())
}
