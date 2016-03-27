# Middle - Go middleware
A simple middleware handler for your HTTP requests

### Basic usage
```go
func middleware1(w http.ResponseWriter, r *http.Request) bool {
    if ok := someCheck(r); !ok {
        // stops the request chain
        return false
    }
    // continue to the next middleware
    return true
}

func yourHttpRequest(w http.ResponseWriter, r *http.Request) {
    w.Write([]byte("Hello simple middleware")
}

m := middle.New()
m.UseBefore(middleware1, middleware2)
m.UseAfter(middleware3, middleware4)

http.HandleFunc("/", m.Then(yourHttpRequest))
```

### Using chaing
Chain your middlewares in the http handler
```go
m := middle.New()
m.UseBefore(middleware1)

http.HandleFunc("/one", m.Before(middleware2).Then(request))
http.HandleFunc("/two", m.Before(middleware3).After(middleware4).Then(request))
```

### Wrap the request
```go
func wrapper(w http.ResponseWriter, r *http.Request, next middle.Request) {
    w.Write([]byte("Ex: open database connection")
    next(w, r)
    w.Write([]byte("Close the connection")
}

m := middle.New()
m.UseWrap(wrapper)

http.HandleFunc("/", m.Wrap(wrapper).Then(request))
```