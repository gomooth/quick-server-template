package global

import (
	"context"
	"errors"
	"log/slog"
	"sync"
)

// 资源回收注册中心
//
// 设计：各基础设施组件（Redis/Cache/Locker/Producer 等）在初始化成功后
// 调用 RegisterRelease 注册自己的回收函数。进程退出时由 boot 层统一调用
// Release，按 LIFO 顺序释放（后注册的先释放，符合"上层先关"直觉），
// 避免连接泄漏和消息丢失。
//
// 释放函数必须是幂等的（多次调用不 panic），实现时应用 sync.Once 或 nil 置位保护。

var (
	releaseMu    sync.Mutex
	releaseHooks []func(ctx context.Context) error
	released     bool
)

// RegisterRelease 注册一个资源回收钩子。
// 在进程退出时按 LIFO 顺序调用。重复注册同一函数会被追加（调用方自行去重）。
// 在 Release 之后调用属于生命周期倒置的编程错误，将 panic。
func RegisterRelease(fn func(ctx context.Context) error) {
	if fn == nil {
		return
	}
	releaseMu.Lock()
	defer releaseMu.Unlock()
	if released {
		panic("global: RegisterRelease called after Release (lifecycle inversion)")
	}
	releaseHooks = append(releaseHooks, fn)
}

// Release 释放所有已注册的资源，按 LIFO 顺序执行。
// 幂等：多次调用只真正执行一次。返回所有钩子的聚合错误。
func Release(ctx context.Context) error {
	releaseMu.Lock()
	if released {
		releaseMu.Unlock()
		return nil
	}
	released = true
	hooks := releaseHooks
	releaseHooks = nil
	releaseMu.Unlock()

	var errs []error
	// LIFO：后注册的先释放（通常是上层资源，如 producer；底层 db/redis 最后关）
	for i := len(hooks) - 1; i >= 0; i-- {
		if err := hooks[i](ctx); err != nil {
			slog.Error("release hook failed", "err", err)
			errs = append(errs, err)
		}
	}
	return errors.Join(errs...)
}

// resetReleaseForTest 仅供测试重置状态，生产不可用。
func resetReleaseForTest() {
	releaseMu.Lock()
	defer releaseMu.Unlock()
	releaseHooks = nil
	released = false
}
