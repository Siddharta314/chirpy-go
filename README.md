go build -o out && ./out
http://localhost:8080



An http.Handler is just an interface:

```go
type Handler interface {
	ServeHTTP(ResponseWriter, *Request)
}
```
Any type with a ServeHTTP method that matches the http.HandlerFunc signature above is an http.Handler. Take a second to think about it: it makes a lot of sense! To handle an incoming HTTP request, all a function needs is a way to write a response and the request itself.