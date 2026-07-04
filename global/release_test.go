package global

import (
	"context"
	"errors"
	"sync"
	"testing"
)

func TestRelease_LIFOOrder(t *testing.T) {
	resetReleaseForTest()

	var order []int
	var mu sync.Mutex

	RegisterRelease(func(ctx context.Context) error {
		mu.Lock()
		order = append(order, 1)
		mu.Unlock()
		return nil
	})
	RegisterRelease(func(ctx context.Context) error {
		mu.Lock()
		order = append(order, 2)
		mu.Unlock()
		return nil
	})
	RegisterRelease(func(ctx context.Context) error {
		mu.Lock()
		order = append(order, 3)
		mu.Unlock()
		return nil
	})

	if err := Release(context.Background()); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	want := []int{3, 2, 1}
	for i, v := range want {
		if order[i] != v {
			t.Fatalf("LIFO order mismatch at %d: got %v, want %v", i, order, want)
		}
	}
}

func TestRelease_Idempotent(t *testing.T) {
	resetReleaseForTest()

	calls := 0
	RegisterRelease(func(ctx context.Context) error {
		calls++
		return nil
	})

	_ = Release(context.Background())
	_ = Release(context.Background())
	_ = Release(context.Background())

	if calls != 1 {
		t.Fatalf("hook should run once, got %d", calls)
	}
}

func TestRelease_AggregatesErrors(t *testing.T) {
	resetReleaseForTest()

	errA := errors.New("a")
	errB := errors.New("b")

	RegisterRelease(func(ctx context.Context) error { return errA })
	RegisterRelease(func(ctx context.Context) error { return errB })

	// LIFO: B 后注册先返回
	err := Release(context.Background())
	if err == nil {
		t.Fatal("expected aggregated error, got nil")
	}

	if !errors.Is(err, errA) {
		t.Errorf("error should contain errA: %v", err)
	}
	if !errors.Is(err, errB) {
		t.Errorf("error should contain errB: %v", err)
	}
}

func TestRegisterRelease_NilNoop(t *testing.T) {
	resetReleaseForTest()

	RegisterRelease(nil)

	if err := Release(context.Background()); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestRegisterRelease_AfterReleasePanics(t *testing.T) {
	resetReleaseForTest()

	RegisterRelease(func(ctx context.Context) error { return nil })
	if err := Release(context.Background()); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	defer func() {
		if r := recover(); r == nil {
			t.Error("expected panic when RegisterRelease called after Release")
		}
	}()
	RegisterRelease(func(ctx context.Context) error { return nil })
}
