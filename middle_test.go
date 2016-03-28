package middle

import (
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

type testCase struct {
	id      string
	actual  string
	handler http.Handler
}

func getTestCases() []testCase {
	return []testCase{
		{
			id:      "SimpleRequest",
			actual:  "final",
			handler: New().Then(final),
		},
		{
			id:     "UseBefore",
			actual: "one two final",
			handler: func() http.Handler {
				m := New()
				m.UseBefore(middleware1, middleware2)
				return m.Then(final)
			}(),
		},
		{
			id:     "UseBeforeWithChaining",
			actual: "one two three four final",
			handler: func() http.Handler {
				m := New()
				m.UseBefore(middleware1, middleware2)
				return m.Before(middleware3, middleware4).Then(final)
			}(),
		},
		{
			id:     "UseAfter",
			actual: "final one two",
			handler: func() http.Handler {
				m := New()
				m.UseAfter(middleware1, middleware2)
				return m.Then(final)
			}(),
		},
		{
			id:     "UseAfterWithChaining",
			actual: "final one two three four",
			handler: func() http.Handler {
				m := New()
				m.UseAfter(middleware1, middleware2)
				return m.After(middleware3, middleware4).Then(final)
			}(),
		},
		{
			id:      "BeforeAndAfter",
			actual:  "one two final three four",
			handler: New().Before(middleware1, middleware2).After(middleware3, middleware4).Then(final),
		},
		{
			id:     "UseWrap",
			actual: "wrapper1_start final wrapper1_end",
			handler: func() http.Handler {
				m := New()
				m.UseWrap(wrapper1)
				return m.Then(final)
			}(),
		},
		{
			id:     "OverrideUseWrap",
			actual: "wrapper2_start final wrapper2_end",
			handler: func() http.Handler {
				m := New()
				m.UseWrap(wrapper1)
				return m.Wrap(wrapper2).Then(final)
			}(),
		},
		{
			id:      "StopMiddleware",
			actual:  "one stop",
			handler: New().Before(middleware1, middlewareStop, middleware2).Then(final),
		},
	}
}

func TestRunTestCases(t *testing.T) {
	for _, testCase := range getTestCases() {
		result, err := generateRequest(testCase.handler)
		if err != nil {
			t.Error(err)
		}

		if testCase.actual != result {
			t.Errorf("ID: %s. Expected: %s Got: %s", testCase.id, testCase.actual, result)
		}
	}
}

func middleware1(w http.ResponseWriter, r *http.Request) bool {
	w.Write([]byte("one "))
	return true
}

func middleware2(w http.ResponseWriter, r *http.Request) bool {
	w.Write([]byte("two "))
	return true
}

func middleware3(w http.ResponseWriter, r *http.Request) bool {
	w.Write([]byte("three "))
	return true
}

func middleware4(w http.ResponseWriter, r *http.Request) bool {
	w.Write([]byte("four "))
	return true
}

func middleware5(w http.ResponseWriter, r *http.Request) bool {
	w.Write([]byte("five "))
	return true
}

func middlewareStop(w http.ResponseWriter, r *http.Request) bool {
	w.Write([]byte("stop "))
	return false
}

func wrapper1(w http.ResponseWriter, r *http.Request, next Request) {
	w.Write([]byte("wrapper1_start "))
	next(w, r)
	w.Write([]byte("wrapper1_end "))
}

func wrapper2(w http.ResponseWriter, r *http.Request, next Request) {
	w.Write([]byte("wrapper2_start "))
	next(w, r)
	w.Write([]byte("wrapper2_end "))
}

func final(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("final "))
}

func generateRequest(handler http.Handler) (string, error) {
	ts := httptest.NewServer(handler)
	defer ts.Close()

	res, err := http.Get(ts.URL)
	if err != nil {
		return "", err
	}

	result, err := ioutil.ReadAll(res.Body)
	res.Body.Close()
	if err != nil {
		return "", err
	}

	return strings.Trim(string(result), " "), nil
}
