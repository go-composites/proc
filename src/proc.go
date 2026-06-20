package Proc

import (
	Error "github.com/go-composites/error/src"
	Result "github.com/go-composites/result/src"
)

/*
Interface is a first-class callable composite — the go-composites analogue of
Ruby's Proc/lambda. A Proc wraps a single function as a value you can pass
around, invoke, and chain. Every Proc is a real object: it never returns a bare
nil, it answers IsNull(), and even the Null Proc is safe to Call.
*/
type Interface interface {
	// Call invokes the wrapped function with the given arguments and returns
	// its Result. A Proc whose wrapped function is nil does not panic: it
	// returns an error Result instead.
	Call(args ...interface{}) Result.Interface

	// Then returns a new Proc that calls the receiver and, on a successful
	// Result, feeds that Result's payload as the single argument to next.Call.
	// If the receiver yields an error Result, Then short-circuits and returns
	// that error unchanged — railway-oriented composition — and next is never
	// called.
	Then(next Interface) Interface

	// IsNull reports whether this is the Null Proc (the not-callable
	// null-object).
	IsNull() bool
}

// proc is the ordinary, callable Proc. It wraps a single function value.
type proc struct {
	fn func(args ...interface{}) Result.Interface
}

/*
New wraps fn as a first-class Proc. fn receives the arguments passed to Call and
returns a Result. If fn is nil the resulting Proc is still safe: Call returns an
error Result rather than panicking.
*/
func New(fn func(args ...interface{}) Result.Interface) Interface {
	return proc{fn: fn}
}

// Call invokes the wrapped function, guarding against a nil function value.
func (p proc) Call(args ...interface{}) Result.Interface {
	if p.fn == nil {
		return Result.New(Result.WithError(Error.New("Proc.Call: nil function")))
	}
	return p.fn(args...)
}

// Then composes the receiver with next, railway-style.
func (p proc) Then(next Interface) Interface {
	return New(func(args ...interface{}) Result.Interface {
		result := p.Call(args...)
		if result.HasError() {
			return result
		}
		return next.Call(result.Payload())
	})
}

// IsNull reports that an ordinary Proc is not the null-object.
func (p proc) IsNull() bool {
	return false
}

// null is the Null-Object Proc: a not-callable Proc that is always safe to use.
type null struct{}

/*
Null returns the Null Proc — the null-object of this composite. It is not
callable: Call returns an error Result, Then returns the Null Proc again (so
chains stay on the null), and IsNull reports true.
*/
func Null() Interface {
	return null{}
}

// Call on the Null Proc always fails with a descriptive error Result.
func (n null) Call(args ...interface{}) Result.Interface {
	return Result.New(Result.WithError(Error.New("Proc.Null: not callable")))
}

// Then on the Null Proc stays on the null: it returns the Null Proc unchanged.
func (n null) Then(next Interface) Interface {
	return n
}

// IsNull reports that this is the null-object.
func (n null) IsNull() bool {
	return true
}
