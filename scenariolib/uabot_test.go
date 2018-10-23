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

type ExpectedRequest struct {
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

func createMockServer(t testing.TB, expected map[string]ExpectedRequest) *httptest.Server {
	// Start a local HTTP server to intercept requests
	// Using url to match the expected responses above.
	server := httptest.NewServer(http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {
		url := req.URL.String()

		expReq, exists := expected[url]
		assert(t, exists, "MISSING expected request for %s", url)

		// Test request parameters
		equals(t, req.Method, expReq.Method)

		for k, v := range expReq.Headers {
			equals(t, v, req.Header.Get(k))
		}

		body, err := ioutil.ReadAll(req.Body)
		if err != nil {
			panic(err)
		}

		eq, err := JSONBytesEqual(expReq.Body, body)
		assert(t, eq, "The Request's body is not what we expected\nGot: %s\nExp: %s", body, req.Body)

		// Send back a static response
		rw.Write([]byte(`{"status":"OK"}`))
	}))

	return server
}
