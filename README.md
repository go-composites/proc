<p align="center"><img src="https://raw.githubusercontent.com/go-composites/brand/main/social/go-composites.png" alt="go-composites/proc" width="720"></p>

# proc

[![ci](https://github.com/go-composites/proc/actions/workflows/ci.yml/badge.svg)](https://github.com/go-composites/proc/actions/workflows/ci.yml)

`proc` is the first-class **callable** composite of [go-composites](https://github.com/go-composites) — the Go analogue of Ruby's `Proc`/lambda. A `Proc` wraps a single function as a value you can pass around, invoke with `Call`, and chain with a railway-oriented `Then`. True to the org's discipline it is interface-first, never returns a bare `nil`, carries errors as a [`Result`](https://github.com/go-composites/result), and ships a [Null-Object](https://github.com/go-composites/null) that is always safe to use.

## Proc vs. compose

[`compose`](https://github.com/go-composites/compose) threads a payload through a *pipeline of steps*. `proc` is one level down: it wraps a *single callable* as a first-class value. `Then` then lets you chain individual `Proc`s railway-style — on success the payload of one becomes the single argument to the next; on error the chain short-circuits.

## Install

```sh
go get github.com/go-composites/proc@main
```

## API

```go
type Interface interface {
	Call(args ...interface{}) Result.Interface  // invoke the wrapped fn (nil fn -> error Result, no panic)
	Then(next Interface) Interface               // railway compose: feed payload to next on success, short-circuit on error
	IsNull() bool
}

func New(fn func(args ...interface{}) Result.Interface) Interface  // wrap a function as a Proc
func Null() Interface                                              // the not-callable Null Proc
```

Semantics:

- **`New(fn)`** wraps `fn` as a `Proc`. A `Proc` built from a `nil` function is still safe — its `Call` returns an error `Result` (`"Proc.Call: nil function"`) rather than panicking.
- **`Call(args...)`** invokes the wrapped function and returns its `Result`.
- **`Then(next)`** returns a new `Proc` that calls the receiver and, on a successful `Result`, feeds that payload as the single argument to `next.Call(...)`. If the receiver yields an error `Result`, `Then` short-circuits and returns it unchanged — `next` is never called.
- **`Null()`** is the null-object: `Call` returns `"Proc.Null: not callable"`, `Then` returns the `Null` Proc again, and `IsNull()` is `true`.

## Usage

```go
import (
	Proc "github.com/go-composites/proc/src"
	Result "github.com/go-composites/result/src"
	Error "github.com/go-composites/error/src"
)

parse := Proc.New(func(args ...interface{}) Result.Interface {
	n, err := strconv.Atoi(args[0].(string))
	if err != nil {
		return Result.New(Result.WithError(Error.New("not an integer")))
	}
	return Result.New(Result.WithPayload(n))
})

double := Proc.New(func(args ...interface{}) Result.Interface {
	return Result.New(Result.WithPayload(args[0].(int) * 2))
})

pipeline := parse.Then(double)

ok := pipeline.Call("21")            // HasError=false, Payload=42
bad := pipeline.Call("not-a-number") // HasError=true, Error().Message()="not an integer" (double skipped)
```

See [`main.go`](main.go) for a runnable demo of the success path, the short-circuit-on-error path, the nil-function guard, and the Null Proc.

## License

BSD-3-Clause. See [LICENSE](LICENSE).
