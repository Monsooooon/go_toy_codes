package contextimpl

import (
	"math"
	"testing"
	"time"
)

func TestBackgroundNotTODO(t *testing.T) {
	bg := Background()
	todo := TODO()
	if bg == todo {
		t.Errorf("Backgound and TODO are equal: %v vs %v", bg, todo)
	}
}

func TestWithCanccel(t *testing.T) {
	ctx, cancel := WithCancel(Background())

	if err := ctx.Err(); err != nil {
		t.Errorf("error should be nil, got %v", err)
	}

	cancel()
	<-ctx.Done()
	if err := ctx.Err(); err != Canceled {
		t.Errorf("error should be Canceled, got %v", err)
	}
}

func TestWithCancelConcurrent(t *testing.T) {
	ctx, cancel := WithCancel(Background())

	time.AfterFunc(1*time.Second, cancel)

	if err := ctx.Err(); err != nil {
		t.Errorf("error should be nil, got %v", err)
	}

	<-ctx.Done()
	if err := ctx.Err(); err != Canceled {
		t.Errorf("error should be Canceled, got %v", err)
	}
}

func TestWithCancelPropagation(t *testing.T) {
	ctxA, cancelA := WithCancel(Background())
	ctxB, cancelB := WithCancel(ctxA)
	defer cancelB()

	cancelA()
	select {
	case <-ctxB.Done():
	case <-time.After(1 * time.Second):
		t.Error("ctx.Done() time out")
	}

	if err := ctxB.Err(); err != Canceled {
		t.Errorf("error should be Canceled, got %v", err)
	}
}

func TestWithDeadline(t *testing.T) {
	deadline := time.Now().Add(2 * time.Second)

	ctx, cancel := WithDeadline(Background(), deadline)

	if d, ok := ctx.Deadline(); d != deadline || ok != true {
		t.Errorf("expected deadline, ok (%v, true), got (%v, %v)", deadline, d, ok)
	}

	then := time.Now()
	<-ctx.Done()

	if d := time.Since(then); math.Abs(d.Seconds()-2.0) > 0.1 {
		t.Errorf("should have been done after 2.0 seconds, took %v", d)
	}

	if err := ctx.Err(); err != DeadlineExceeded {
		t.Errorf("error should be DeadlineExceeded, got %v", err)
	}

	cancel()
	if err := ctx.Err(); err != DeadlineExceeded {
		t.Errorf("error should be DeadlineExceeded, got %v", err)
	}
}

func TestTimeOut(t *testing.T) {
	timeout := 2 * time.Second
	deadline := time.Now().Add(timeout)

	ctx, cancel := WithTimeout(Background(), timeout)

	if d, ok := ctx.Deadline(); d.Sub(deadline) > 10*time.Microsecond || ok != true {
		t.Errorf("expected deadline, ok (%v, true), got (%v, %v)", deadline, d, ok)
	}

	then := time.Now()
	<-ctx.Done()

	if d := time.Since(then); math.Abs(d.Seconds()-2.0) > 0.1 {
		t.Errorf("should have been done after 2.0 seconds, took %v", d)
	}

	if err := ctx.Err(); err != DeadlineExceeded {
		t.Errorf("error should be DeadlineExceeded, got %v", err)
	}

	cancel()
	if err := ctx.Err(); err != DeadlineExceeded {
		t.Errorf("error should be DeadlineExceeded, got %v", err)
	}
}

func TestWithValue(t *testing.T) {
	testcases := []struct {
		key, val, keyRet, valRet interface{}
		shouldPanic              bool
	}{
		{"a", "b", "a", "b", false},
		{"a", "b", "c", nil, false},
		{42, true, 42, true, false},
		{42, true, int64(42), nil, false},
		{nil, true, nil, nil, true},                       // nil key
		{[]int{1, 2, 3}, true, []int{1, 2, 3}, nil, true}, // incomparable key
	}

	var panicked interface{}
	for _, tc := range testcases {
		func() {
			defer func() {
				panicked = recover()
			}()

			ctx := WithValue(Background(), tc.key, tc.val)
			if val := ctx.Value(tc.keyRet); val != tc.valRet {
				t.Errorf("expected value for key %v is %v, got %v", tc.keyRet, tc.valRet, val)
			}

		}()

		if tc.shouldPanic && panicked == nil {
			t.Errorf("inserting kv pair (%v, %v) should panic, but didn't get it", tc.key, tc.val)
		}
		if !tc.shouldPanic && panicked != nil {
			t.Errorf("inserting kv pair (%v, %v) shouldn't panic, but get it", tc.key, tc.val)
		}
	}

}
