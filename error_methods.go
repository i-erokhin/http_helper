package http_helper

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"

	"github.com/getsentry/raven-go"
)

func (h *Helper) Error500(w http.ResponseWriter, r *http.Request, e error) {
	h.errorHelper(w, r, e, http.StatusInternalServerError, e.Error(), true)
}

func (h *Helper) ErrorUnexpected(w http.ResponseWriter, r *http.Request, e error, status int, message interface{}) {
	h.errorHelper(w, r, e, status, message, true)
}

func (h *Helper) ErrorExpected(w http.ResponseWriter, r *http.Request, e error, status int, message interface{}) {
	h.errorHelper(w, r, e, status, message, false)
}

func (h *Helper) errorHelper(w http.ResponseWriter, r *http.Request, e error, status int, msg interface{}, writeToSentry bool) {
	message := ""
	ok := false
	var jsonPayload []byte

	// msg is not string - marshaling
	if message, ok = msg.(string); !ok {
		var err error
		jsonPayload, err = json.Marshal(msg)
		if err != nil {
			panic(err)
		}
		message = string(jsonPayload)
	}

	if e == nil {
		fmt.Fprintf(os.Stderr, "%s %s - %d: %q\n", r.Method, r.RequestURI, status, message)
	} else {
		fmt.Fprintf(os.Stderr, "%s %s - %d: %q, error: %v\n", r.Method, r.RequestURI, status, message, e)
	}

	// write to Sentry
	if writeToSentry && h.sentryDSN != "" {
		client, err := raven.New(h.sentryDSN)
		if err != nil {
			panic(err)
		}
		defer client.Close()

		client.SetTagsContext(map[string]string{"status": fmt.Sprintf("%d: %s", status, http.StatusText(status))})
		client.SetHttpContext(raven.NewHttp(r))
		if e != nil {
			client.CaptureError(e, map[string]string{"message": message})
		}
	}

	w.WriteHeader(status)

	// jsonPayload have priority
	if jsonPayload != nil {
		w.Header().Set("Content-Type", "application/json")
		w.Write(jsonPayload)
	} else {
		if message == "" {
			message = http.StatusText(status)
		}
		fmt.Fprintln(w, message)
	}
}
