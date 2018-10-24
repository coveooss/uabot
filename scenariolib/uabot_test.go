package scenariolib_test

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"path/filepath"
	"reflect"
	"runtime"
	"testing"
)

// assert fails the test if the condition is false.
func assert(tb testing.TB, condition bool, msg string, v ...interface{}) {
	if !condition {
		_, file, line, _ := runtime.Caller(1)
		fmt.Printf("%s:%d: "+msg+"\n\n", append([]interface{}{filepath.Base(file), line}, v...)...)
		tb.FailNow()
	}
}

// ok fails the test if an err is not nil.
func ok(tb testing.TB, err error) {
	if err != nil {
		_, file, line, _ := runtime.Caller(1)
		fmt.Printf("%s:%d: unexpected error: %s\n\n", filepath.Base(file), line, err.Error())
		tb.FailNow()
	}
}

// notok fails the test if an err is nil.
func notok(tb testing.TB, err error) {
	if err == nil {
		_, file, line, _ := runtime.Caller(1)
		fmt.Printf("%s:%d: was expecting error but was nil\n\n", filepath.Base(file), line)
		tb.FailNow()
	}
}

// equals fails the test if exp is not equal to act.
func equals(tb testing.TB, exp, act interface{}) {
	if !reflect.DeepEqual(exp, act) {
		_, file, line, _ := runtime.Caller(1)
		fmt.Printf("%s:%d:\n\n\texp: %#v\n\n\tgot: %#v\n\n", filepath.Base(file), line, exp, act)
		tb.FailNow()
	}
}

type RestRequest struct {
	Method  string
	Headers map[string]string
	Body    []byte
}

// JSONBytesEqual compares the JSON in two byte slices.
func JSONBytesEqual(a, b []byte) (bool, error) {
	var j, j2 interface{}
	if err := json.Unmarshal(a, &j); err != nil {
		return false, err
	}
	if err := json.Unmarshal(b, &j2); err != nil {
		return false, err
	}
	return reflect.DeepEqual(j2, j), nil
}

// This a proxy to intercept requests sent to platform endpoints and add them into a map (receivedRequests).
// receivedRequests will then be used for validation in unit tests.
func createTestServer(t testing.TB, receivedRequests map[string]RestRequest) *httptest.Server {
	// Start a local HTTP server to intercept requests
	// Using url as key for the map of received requests.
	server := httptest.NewServer(http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {
		url := req.URL.String()

		headers := make(map[string]string)

		for k := range req.Header {
			headers[k] = req.Header.Get(k)
		}
		body, err := ioutil.ReadAll(req.Body)
		if err != nil {
			panic(err)
		}
		receivedRequests[url] = RestRequest{
			Method:  req.Method,
			Headers: headers,
			Body:    body,
		}

		// Send back a static response
		rw.Write([]byte(`{"status":"OK"}`))
	}))

	return server
}
