package main

import (
	"fmt"
	"strconv"
	"strings"

	Error "github.com/go-composites/error/src"
	Proc "github.com/go-composites/proc/src"
	Result "github.com/go-composites/result/src"
)

func main() {
	// A first-class callable: parse text into an integer, failing with an
	// error Result when the text is not a number.
	parse := Proc.New(func(args ...interface{}) Result.Interface {
		text := strings.TrimSpace(args[0].(string))
		n, err := strconv.Atoi(text)
		if err != nil {
			return Result.New(Result.WithError(
				Error.New("parse: " + strconv.Quote(text) + " is not an integer")))
		}
		return Result.New(Result.WithPayload(n))
	})

	// Another callable: double its single argument. It cannot fail.
	double := Proc.New(func(args ...interface{}) Result.Interface {
		return Result.New(Result.WithPayload(args[0].(int) * 2))
	})

	// Chain them railway-style: parse, then on success feed the payload to double.
	pipeline := parse.Then(double)

	fmt.Println("== success path ==")
	ok := pipeline.Call("21")
	fmt.Printf("Call(%q) -> HasError=%t payload=%v\n", "21", ok.HasError(), ok.Payload())

	fmt.Println("== short-circuit on error ==")
	bad := pipeline.Call("not-a-number")
	fmt.Printf("Call(%q) -> HasError=%t error=%q (double skipped)\n",
		"not-a-number", bad.HasError(), bad.Error().Message())

	fmt.Println("== nil function is safe ==")
	broken := Proc.New(nil)
	out := broken.Call()
	fmt.Printf("Proc.New(nil).Call() -> HasError=%t error=%q\n",
		out.HasError(), out.Error().Message())

	fmt.Println("== Null Proc ==")
	null := Proc.Null()
	nullOut := null.Call("anything")
	fmt.Printf("Proc.Null().Call(...) -> IsNull=%t HasError=%t error=%q\n",
		null.IsNull(), nullOut.HasError(), nullOut.Error().Message())
}
