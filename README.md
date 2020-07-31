# tarent/middleware-jwt [![Build Status](https://travis-ci.org/tarent/middleware-jwt.svg)](https://travis-ci.org/tarent/middleware-jwt)

`middleware-jwt` provides JWT authentication for Go. It is compatible with Go's own `net/http` and anything that speaks the `http.Handler` interface.

## Examples

See the examples folder for working examples.

## Installation

```sh
$ go get github.com/tarent/middleware-jwt
```


## Usage

Since it's all `http.Handler`, `middleware-jwt` works with [gorilla/mux](https://github.com/gorilla/mux) and most other routers.

### Basic usage

```go
[include](File:examples/gorilla/simple/main.go)

```

### Advanced usage

For complex authentication scenarios, you can access the "claims" in the http context which is passed to the subsequent http handlers. This allows you to:

```go
[include](File:examples/gorilla/advanced/main.go)

```


## License

MIT Licensed. See the LICENSE file for details.