# regexrouter

regexrouter is a golang router

## Usage

Usage looks like this:

```go
subrouter := regexrouter.NewRouter()
subrouter.Handle("^/builds/?$", map[string]http.Handler{
    http.MethodGet:  http.HandlerFunc(getRecentBuildsHandler),
    http.MethodPost: http.HandlerFunc(createBuildHandler),
})
subrouter.Handle("^/builds/([0-9]+)/?$", map[string]http.Handler{
    http.MethodGet: http.HandlerFunc(getBuildByIDHandler),
})

router := regexrouter.NewRouter()
router.HandleAny("^/builds/", subrouter)
router.HandleAny("^/", http.HandlerFunc(rootHandler))
```

You can also use named capture groups and access them through `http.Request.Context()`.
