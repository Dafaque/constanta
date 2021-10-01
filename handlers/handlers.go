package handlers

import (
	"fmt"
	"net/http"

	"github.com/Dafaque/constanta/model"
)

const MAX_URLS int = 20

func Root(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		//? POST Only
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}

	//? Read req body into struct
	req, errParseRequest := model.NewRequest(r)
	if errParseRequest != nil || len(req.Urls) > MAX_URLS {
		w.WriteHeader(http.StatusBadRequest)
		fmt.Println(errParseRequest)
		return
	}
	resp, err := req.Process()

	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(err.Error()))
	}

	w.Write(resp.Bytes())
}
