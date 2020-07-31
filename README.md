# tarent/middleware-jwt [![Build Status](https://travis-ci.org/tarent/middleware-jwt.svg)](https://travis-ci.org/tarent/middleware-jwt)

`middleware-jwt` provides JWT authentication for Go. It is compatible with Go's own `net/http` and anything that speaks the `http.Handler` interface.

## Examples

See the `examples` folder for working examples.

## Installation

```sh
$ go get github.com/tarent/middleware-jwt
```

## Usage

Since it's all `http.Handler`, `middleware-jwt` works with [gorilla/mux](https://github.com/gorilla/mux) and most other routers.


### Basic usage

To run the simple example execute:
```sh
go run examples/gorilla/simple/main.go
```

To test it, execute the following curl in another terminal:
```sh
curl -v -H "Authorization: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJzdWIiOiIxMjM0NTY3ODkwIiwibmFtZSI6IkpvaG4gRG9lIiwiZ3JvdXBzIjpbImdyb3VwQSIsImdyb3VwQiJdLCJpYXQiOjE1MTYyMzkwMjJ9.pPJGnFh4FUJnIcnReZlrrraG0Ep_bqEadYo6iH4KdHY" localhost:8080
```

### Advanced usage

For complex authentication scenarios, you can access the "claims" in the http context which is passed to the subsequent http handlers. This allows you to:

To run the advanced example execute:
```sh
go run examples/gorilla/advanced/main.go
```

To test it, execute the following curl in another terminal:
```sh
curl -v -H "Authorization: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJzdWIiOiIxMjM0NTY3ODkwIiwibmFtZSI6IkpvaG4gRG9lIiwiZ3JvdXBzIjpbImdyb3VwQSIsImdyb3VwQiJdLCJpYXQiOjE1MTYyMzkwMjJ9.pPJGnFh4FUJnIcnReZlrrraG0Ep_bqEadYo6iH4KdHY" localhost:8080
```

## License

MIT Licensed. See the LICENSE file for details.