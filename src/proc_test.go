package Proc_test

import (
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	Error "github.com/go-composites/error/src"
	Proc "github.com/go-composites/proc/src"
	Result "github.com/go-composites/result/src"
)

// ok wraps a value in a successful Result.
func ok(value interface{}) Result.Interface {
	return Result.New(Result.WithPayload(value))
}

// fail wraps a message in an error Result.
func fail(message string) Result.Interface {
	return Result.New(Result.WithError(Error.New(message)))
}

var _ = Describe("Proc", func() {
	Describe("New / Call", func() {
		It("invokes the wrapped function and returns its Result on success", func() {
			double := Proc.New(func(args ...interface{}) Result.Interface {
				return ok(args[0].(int) * 2)
			})
			out := double.Call(21)
			Expect(out).NotTo(BeNil())
			Expect(out.HasError()).To(BeFalse())
			Expect(out.Payload()).To(Equal(42))
		})

		It("returns the error Result the wrapped function produces", func() {
			boom := Proc.New(func(args ...interface{}) Result.Interface {
				return fail("boom")
			})
			out := boom.Call()
			Expect(out.HasError()).To(BeTrue())
			Expect(out.Error().Message()).To(Equal("boom"))
		})

		It("does not panic on a nil function — it returns an error Result", func() {
			nilProc := Proc.New(nil)
			out := nilProc.Call(1, 2, 3)
			Expect(out).NotTo(BeNil())
			Expect(out.HasError()).To(BeTrue())
			Expect(out.Error().Message()).To(Equal("Proc.Call: nil function"))
		})

		It("reports IsNull false for an ordinary Proc", func() {
			Expect(Proc.New(nil).IsNull()).To(BeFalse())
		})
	})

	Describe("Then", func() {
		It("feeds the first Proc's payload to the next on the success path", func() {
			incr := Proc.New(func(args ...interface{}) Result.Interface {
				return ok(args[0].(int) + 1)
			})
			triple := Proc.New(func(args ...interface{}) Result.Interface {
				return ok(args[0].(int) * 3)
			})
			out := incr.Then(triple).Call(4)
			Expect(out.HasError()).To(BeFalse())
			Expect(out.Payload()).To(Equal(15)) // (4+1)*3
		})

		It("short-circuits when the first Proc errors and never calls next", func() {
			boom := Proc.New(func(args ...interface{}) Result.Interface {
				return fail("first failed")
			})
			nextRan := false
			next := Proc.New(func(args ...interface{}) Result.Interface {
				nextRan = true
				return ok("unreachable")
			})
			out := boom.Then(next).Call("x")
			Expect(out.HasError()).To(BeTrue())
			Expect(out.Error().Message()).To(Equal("first failed"))
			Expect(nextRan).To(BeFalse(), "next must not run when the first Proc errors")
		})

		It("propagates an error produced by the next Proc", func() {
			first := Proc.New(func(args ...interface{}) Result.Interface {
				return ok(args[0])
			})
			next := Proc.New(func(args ...interface{}) Result.Interface {
				return fail("second failed")
			})
			out := first.Then(next).Call(7)
			Expect(out.HasError()).To(BeTrue())
			Expect(out.Error().Message()).To(Equal("second failed"))
		})
	})

	Describe("Null", func() {
		It("reports IsNull true", func() {
			Expect(Proc.Null().IsNull()).To(BeTrue())
		})

		It("Call returns a not-callable error Result", func() {
			out := Proc.Null().Call(1, 2)
			Expect(out).NotTo(BeNil())
			Expect(out.HasError()).To(BeTrue())
			Expect(out.Error().Message()).To(Equal("Proc.Null: not callable"))
		})

		It("Then stays on the Null Proc", func() {
			next := Proc.New(func(args ...interface{}) Result.Interface {
				return ok("ran")
			})
			chained := Proc.Null().Then(next)
			Expect(chained.IsNull()).To(BeTrue())
			out := chained.Call("x")
			Expect(out.HasError()).To(BeTrue())
			Expect(out.Error().Message()).To(Equal("Proc.Null: not callable"))
		})
	})
})
