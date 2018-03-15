package http_helper

import (
	"encoding/json"
	"fmt"
	"net/http"
)

func (h *Helper) Ok(w http.ResponseWriter, response interface{}) {
	if responseString, ok := response.(string); ok {
		fmt.Fprintln(w, responseString)
	} else {
		b, err := json.Marshal(response)
		if err != nil {
			panic(err)
		}
		w.Header().Set("Content-Type", "application/json")
		w.Write(b)
	}
	w.WriteHeader(http.StatusOK)
}
