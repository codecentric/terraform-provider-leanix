package leanix

import (
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"reflect"
	"testing"
)

func assertEqual(t *testing.T, actual interface{}, expected interface{}) {
	if !reflect.DeepEqual(expected, actual) {
		t.Fatalf("Expected %v to be equal to %v", expected, actual)
	}
}

type TestRoute map[TestResourceAndMethod]*TestRouteDefinition

type TestRouteDefinition struct {
	ExpectedHeader map[string]string
	ExpectedBody   []byte
	ResponseStatus func(http.Header, []byte) int
	ResponseBody   func(http.Header, []byte) []byte
}

type TestResourceAndMethod struct {
	Resource string
	Method   string
}

func NewTestServer(t *testing.T, route TestRoute) *httptest.Server {
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		requestRoute := TestResourceAndMethod{r.URL.EscapedPath(), r.Method}
		matchingRoute := route[requestRoute]
		if matchingRoute == nil {
			t.Fatalf("No route matches '%s'", requestRoute)
		}

		for k, v := range matchingRoute.ExpectedHeader {
			if r.Header.Get(k) != v {
				t.Fatalf("Expected header '%s:%s'", k, v)
			}
		}

		defer r.Body.Close()
		bodyBytes, err := ioutil.ReadAll(r.Body)
		if err != nil {
			t.Fatal(err)
		}
		assertEqual(t, bodyBytes, matchingRoute.ExpectedBody)

		w.WriteHeader(matchingRoute.ResponseStatus(r.Header, bodyBytes))
		w.Write(matchingRoute.ResponseBody(r.Header, bodyBytes))
	})
	return httptest.NewServer(handler)
}
