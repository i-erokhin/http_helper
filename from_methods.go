package http_helper

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"
)

type Validator interface {
	Validate() []string
}

func (h *Helper) FromGet(w http.ResponseWriter, r *http.Request, target Validator) error {
	return h.from("GET", w, r, target)
}

func (h *Helper) FromPost(w http.ResponseWriter, r *http.Request, target Validator) error {
	return h.from("POST", w, r, target)
}

func (h *Helper) FromPut(w http.ResponseWriter, r *http.Request, target Validator) error {
	return h.from("PUT", w, r, target)
}

func (h *Helper) FromDelete(w http.ResponseWriter, r *http.Request, target Validator) error {
	return h.from("DELETE", w, r, target)
}

func (h *Helper) from(method string, w http.ResponseWriter, r *http.Request, target Validator) error {
	if r.Method != method {
		w.WriteHeader(http.StatusMethodNotAllowed)
		message := http.StatusText(http.StatusMethodNotAllowed)
		err := errors.New(message)
		h.ErrorExpected(w, r, err, http.StatusMethodNotAllowed, message)
		return err
	}

	if target == nil {
		return nil
	}

	// populate target Validator
	if method == "GET" {
		r.URL.Query()

		if err := r.ParseForm(); err != nil {
			panic(err)
		}

		m := map[string]string{}
		for k, v := range r.Form {
			m[k] = v[0]
		}

		b, err := json.Marshal(m)
		if err != nil {
			panic(err)
		}

		if err = json.Unmarshal(b, target); err != nil {
			h.ErrorExpected(w, r, err, http.StatusBadRequest, fmt.Sprintf("Bad GET parameters: %s", err))
			return err
		}
	} else {
		body, err := ioutil.ReadAll(r.Body)
		if err != nil {
			panic(err)
		}

		if err = json.Unmarshal(body, target); err != nil {
			h.ErrorExpected(w, r, err, http.StatusBadRequest, fmt.Sprintf("Bad JSON format: %s", err))
			return err
		}
	}

	if errs := target.Validate(); errs != nil {
		var errStr string
		if len(errs) == 1 {
			errStr = errs[0]
		} else if len(errs) > 1 {
			errStr = "Errors: " + strings.Join(errs, "; ")
		} else {
			panic("Array of errors string can`be empty, nil expected instead.")
		}
		err := errors.New(errStr)
		h.ErrorExpected(w, r, err, http.StatusBadRequest, errStr)
		return err
	}

	return nil
}
