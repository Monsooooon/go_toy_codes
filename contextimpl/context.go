package contextimpl

import (
	"errors"
	"reflect"
	"sync"
	"time"
)

type Context interface {
	Deadline() (deadline time.Time, ok bool)
	Done() <-chan struct{}
	Err() error
	Value(key interface{}) interface{}
}

/* -----------Background & TODO------------- */
type emptyCtx int

var (
	background = new(emptyCtx)
	todo       = new(emptyCtx)
)

func Background() Context {
	return background
}

func TODO() Context {
	return todo
}

func (emptyCtx) Deadline() (deadline time.Time, ok bool) { return time.Time{}, false }
func (emptyCtx) Done() <-chan struct{}                   { return nil }
func (emptyCtx) Err() error                              { return nil }
func (emptyCtx) Value(key interface{}) interface{}       { return nil }

/* -----------WithCancel------------- */

var Canceled = errors.New("context canceled")

type cancelCxt struct {
	Context
	done chan struct{}
	err  error
	mu   sync.Mutex
}

func (ctx *cancelCxt) Done() <-chan struct{} { return ctx.done } //  <- chan means you can only send through it, but not receive
func (ctx *cancelCxt) Err() error {
	ctx.mu.Lock()
	defer ctx.mu.Unlock()
	return ctx.err
}

type CancelFunc func()

func WithCancel(parent Context) (Context, CancelFunc) {
	ctx := &cancelCxt{
		Context: parent,
		done:    make(chan struct{}),
	}

	cancel := func() {
		ctx.cancel(Canceled)
	}

	go func() {
		select {
		case <-parent.Done():
			ctx.cancel(parent.Err())
		case <-ctx.Done():
		}
	}()

	return ctx, cancel
}

func (ctx *cancelCxt) cancel(err error) {
	ctx.mu.Lock()
	defer ctx.mu.Unlock()
	if ctx.err != nil { // ctx.done should be closed once
		return
	}
	ctx.err = err
	close(ctx.done)
}

/* -----------WithDeadline & WithTimeout------------- */

type deadlineExceededError struct{}

func (err deadlineExceededError) Error() string {
	return "Deadline is exceeded"
}

var DeadlineExceeded error = deadlineExceededError{}

type deadlineCtx struct {
	*cancelCxt
	deadline time.Time
}

func (ctx *deadlineCtx) Deadline() (deadline time.Time, ok bool) { return ctx.deadline, true }

func WithDeadline(parent Context, d time.Time) (Context, CancelFunc) {
	cctx, cancel := WithCancel(parent)
	ctx := &deadlineCtx{
		cancelCxt: cctx.(*cancelCxt),
		deadline:  d,
	}
	timer := time.AfterFunc(time.Until(d), func() {
		ctx.cancel(DeadlineExceeded)
	})

	stop := func() {
		timer.Stop()
		cancel()
	}

	return ctx, stop
}

func WithTimeout(parent Context, timeout time.Duration) (Context, CancelFunc) {
	return WithDeadline(parent, time.Now().Add(timeout))
}

/* -----------WithValue------------- */

type valueCtx struct {
	Context
	key, val interface{}
}

func (ctx *valueCtx) Value(key interface{}) interface{} {
	if ctx.key == key {
		return ctx.val
	}
	return ctx.Context.Value(key)
}

func WithValue(parent Context, key, val interface{}) Context {
	if key == nil {
		panic("nil key")
	}
	if !reflect.TypeOf(key).Comparable() {
		panic("incomparable key")
	}
	return &valueCtx{
		Context: parent,
		key:     key,
		val:     val,
	}
}
