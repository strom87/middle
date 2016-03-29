# Middle - go middleware
A middleware handler for golang that helps you keep your HTTP requests simple and elegant

### Get package
```sh
$ go get github.com/strom87/middle
```

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

http.Handle("/one", m.Then(yourHttpRequest))
http.HandleFunc("/two", m.ThenFunc(yourHttpRequest))
```

### Using chaining
Chain your middlewares in the http handler
```go
m := middle.New()
m.UseBefore(middleware1)

http.HandleFunc("/one", m.Before(middleware2).ThenFunc(request))
http.HandleFunc("/two", m.Before(middleware3).After(middleware4).ThenFunc(request))
```

### Wrap the request
Wrap encapsulates the middlewares and runs all requests inside the same context
```go
func wrapper(w http.ResponseWriter, r *http.Request, next middle.Request) {
    w.Write([]byte("Ex: open database connection")
    next(w, r)
    w.Write([]byte("Close the connection")
}

m := middle.New()
m.UseWrap(wrapper)

http.Handle("/", m.Wrap(wrapper).Then(request))
```
### Before
**UseBefore(...)** and **Before(...)** exectues the included middlewares before the main request.  
**UseBefore** adds the middleware to be executed in all of the requests.  
**Before** chaining just executes the inserted middlewares for that single request.

### After
**UseAfter(...)** and **After(...)** exectues the included middlewares after the main request.  
**UseAfter** adds the middleware to be executed in all of the requests.  
**After** chaining just executes the inserted middlewares for that single request.

### Wrap
**UseWrap()** and **Wrap()** executes first and last, all the middlewares and the main request is excecuted in between wrap functions **next()** statement.  
**UseWrap()** will be executed for all the request and can be overidden for a single request with the chaining function **Wrap()**.