package middle

import (
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

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

func TestSimpleRequest(t *testing.T) {
	const actual = "final"

	result, err := generateRequest(New().Then(final))
	if err != nil {
		t.Error(err)
	}

	if actual != result {
		t.Errorf("Expected: %s Got: %s", actual, result)
	}
}

func TestUseBefore(t *testing.T) {
	const actual = "one two final"

	m := New()
	m.UseBefore(middleware1, middleware2)

	result, err := generateRequest(m.Then(final))
	if err != nil {
		t.Error(err)
	}

	if actual != result {
		t.Errorf("Expected: %s Got: %s", actual, result)
	}
}

func TestUseBeforeWithChaining(t *testing.T) {
	const actual = "one two three four final"

	m := New()
	m.UseBefore(middleware1, middleware2)

	result, err := generateRequest(m.Before(middleware3, middleware4).Then(final))
	if err != nil {
		t.Error(err)
	}

	if actual != result {
		t.Errorf("Expected: %s Got: %s", actual, result)
	}
}

func TestUseAfter(t *testing.T) {
	const actual = "final one two"

	m := New()
	m.UseAfter(middleware1, middleware2)

	result, err := generateRequest(m.Then(final))
	if err != nil {
		t.Error(err)
	}

	if actual != result {
		t.Errorf("Expected: %s Got: %s", actual, result)
	}
}

func TestUseAfterWithChaining(t *testing.T) {
	const actual = "final one two three four"

	m := New()
	m.UseAfter(middleware1, middleware2)

	result, err := generateRequest(m.After(middleware3, middleware4).Then(final))
	if err != nil {
		t.Error(err)
	}

	if actual != result {
		t.Errorf("Expected: %s Got: %s", actual, result)
	}
}

func TestBeforeAndAfter(t *testing.T) {
	const actual = "one two final three four"

	result, err := generateRequest(New().Before(middleware1, middleware2).After(middleware3, middleware4).Then(final))
	if err != nil {
		t.Error(err)
	}

	if actual != result {
		t.Errorf("Expected: %s Got: %s", actual, result)
	}
}

func TestUseWrap(t *testing.T) {
	const actual = "wrapper1_start final wrapper1_end"
	m := New()
	m.UseWrap(wrapper1)

	result, err := generateRequest(m.Then(final))
	if err != nil {
		t.Error(err)
	}

	if actual != result {
		t.Errorf("Expected: %s Got: %s", actual, result)
	}
}

func TestOverrideUseWrap(t *testing.T) {
	const actual = "wrapper2_start final wrapper2_end"
	m := New()
	m.UseWrap(wrapper1)

	result, err := generateRequest(m.Wrap(wrapper2).Then(final))
	if err != nil {
		t.Error(err)
	}

	if actual != result {
		t.Errorf("Expected: %s Got: %s", actual, result)
	}
}

func TestStopMiddleware(t *testing.T) {
	const actual = "one stop"

	result, err := generateRequest(New().Before(middleware1, middlewareStop, middleware2).Then(final))
	if err != nil {
		t.Error(err)
	}

	if actual != result {
		t.Errorf("Expected: %s Got: %s", actual, result)
	}
}
