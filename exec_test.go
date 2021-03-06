package piper

import (
	"context"
	"sync/atomic"
	"testing"
	"time"
)

type testExec struct {
	startCountPtr *uint64
	stopCountPtr  *uint64
	startFn       execFnType
	stopFn        execFnType
}

func newTestExec() *testExec {
	var startCount, stopCount uint64
	startFn := func(ctx context.Context) {
		atomic.AddUint64(&startCount, 1)
		time.Sleep(50 * time.Millisecond)
	}
	stopFn := func(ctx context.Context) {
		atomic.AddUint64(&stopCount, 1)
		time.Sleep(50 * time.Millisecond)
	}

	return &testExec{
		startCountPtr: &startCount,
		stopCountPtr:  &stopCount,
		startFn:       startFn,
		stopFn:        stopFn,
	}
}

func TestExec_NewExec(t *testing.T) {
	te := newTestExec()
	e := newExec(te.startFn, te.stopFn)
	if e == nil {
		t.Fatal("newExec returned nil")
	}
}

func TestExec_Start(t *testing.T) {
	te := newTestExec()
	e := newExec(te.startFn, te.stopFn)

	e.start(context.TODO())

	got := atomic.LoadUint64(te.startCountPtr)
	if got != 1 {
		t.Fatalf("startCount invalid: wanted: [%d], got [%d]", 1, got)
	}
}

func TestExec_StartTwice(t *testing.T) {
	te := newTestExec()
	e := newExec(te.startFn, te.stopFn)

	ctx := context.TODO()
	e.start(ctx)
	e.start(ctx)

	got := atomic.LoadUint64(te.startCountPtr)
	if got != 1 {
		t.Fatalf("startCount invalid: wanted: [%d], got [%d]", 1, got)
	}
}

func TestExec_Stop(t *testing.T) {
	te := newTestExec()
	e := newExec(te.startFn, te.stopFn)

	ctx := context.TODO()
	e.stop(ctx)

	got := atomic.LoadUint64(te.stopCountPtr)
	if got != 1 {
		t.Fatalf("stopCount invalid: wanted: [%d], got [%d]", 1, got)
	}
}

func TestExec_StopTwice(t *testing.T) {
	te := newTestExec()
	e := newExec(te.startFn, te.stopFn)

	ctx := context.TODO()
	e.stop(ctx)
	e.stop(ctx)

	got := atomic.LoadUint64(te.stopCountPtr)
	if got != 1 {
		t.Fatalf("stopCount invalid: wanted: [%d], got [%d]", 1, got)
	}
}

func TestExec_StartStop(t *testing.T) {
	te := newTestExec()
	e := newExec(te.startFn, te.stopFn)

	ctx := context.TODO()
	e.start(ctx)
	e.stop(ctx)

	got1 := atomic.LoadUint64(te.startCountPtr)
	if got1 != 1 {
		t.Fatalf("startCount invalid: wanted: [%d], got [%d]", 1, got1)
	}

	got2 := atomic.LoadUint64(te.stopCountPtr)
	if got2 != 1 {
		t.Fatalf("stopCount invalid: wanted: [%d], got [%d]", 1, got2)
	}
}

func TestExec_StartStopTwice(t *testing.T) {
	te := newTestExec()
	e := newExec(te.startFn, te.stopFn)

	ctx := context.TODO()
	e.start(ctx)
	e.stop(ctx)
	e.start(ctx)
	e.stop(ctx)

	got1 := atomic.LoadUint64(te.startCountPtr)
	if got1 != 2 {
		t.Fatalf("startCount invalid: wanted: [%d], got [%d]", 2, got1)
	}

	got2 := atomic.LoadUint64(te.stopCountPtr)
	if got2 != 2 {
		t.Fatalf("stopCount invalid: wanted: [%d], got [%d]", 2, got2)
	}
}
