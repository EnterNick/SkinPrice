package skins

import (
	"SkinPrice/skinprice/internal/shared/errx"
	"context"
	"errors"
	"sync"
)

type DefaultRefreshQueue struct {
	UpdateOne SavedSkinPriceUpdater

	mu      sync.Mutex
	closed  bool
	manual  []queuedRefreshTask
	auto    []queuedRefreshTask
	wake    chan struct{}
	stop    chan struct{}
	started sync.Once
}

type queuedRefreshTask struct {
	task RefreshTask
	done chan queuedRefreshResult
}

type queuedRefreshResult struct {
	result UpdateSavedSkinPriceResult
	err    error
}

func NewRefreshQueue(updateOne SavedSkinPriceUpdater) *DefaultRefreshQueue {
	return &DefaultRefreshQueue{
		UpdateOne: updateOne,
		wake:      make(chan struct{}, 1),
		stop:      make(chan struct{}),
	}
}

func (q *DefaultRefreshQueue) Run(ctx context.Context) {
	q.started.Do(func() {
		go q.loop(ctx)
	})
}

func (q *DefaultRefreshQueue) Enqueue(ctx context.Context, task RefreshTask) (UpdateSavedSkinPriceResult, error) {
	if q.UpdateOne == nil {
		return UpdateSavedSkinPriceResult{}, errx.E("skins.refresh_queue.no_updater", errx.CodeInternal, "no price updater configured", nil)
	}
	item := queuedRefreshTask{task: task, done: make(chan queuedRefreshResult, 1)}
	q.mu.Lock()
	if q.closed {
		q.mu.Unlock()
		return UpdateSavedSkinPriceResult{}, errx.E("skins.refresh_queue.closed", errx.CodeUnavailable, "refresh queue is closed", nil)
	}
	if task.Kind == RefreshTaskAuto {
		q.auto = append(q.auto, item)
	} else {
		q.manual = append(q.manual, item)
	}
	q.mu.Unlock()
	q.notify()

	select {
	case <-ctx.Done():
		return UpdateSavedSkinPriceResult{}, ctx.Err()
	case result := <-item.done:
		return result.result, result.err
	}
}

func (q *DefaultRefreshQueue) Shutdown() {
	q.mu.Lock()
	if !q.closed {
		q.closed = true
		close(q.stop)
	}
	q.mu.Unlock()
	q.notify()
}

func (q *DefaultRefreshQueue) loop(ctx context.Context) {
	for {
		item, ok := q.next()
		if !ok {
			select {
			case <-ctx.Done():
				q.failPending(ctx.Err())
				return
			case <-q.stop:
				q.failPending(errors.New("refresh queue stopped"))
				return
			case <-q.wake:
				continue
			}
		}

		result, err := q.UpdateOne.Execute(ctx, UpdateSavedSkinPriceParams{
			MarketHashName: item.task.MarketHashName,
			Currency:       item.task.Currency,
		})
		item.done <- queuedRefreshResult{result: result, err: err}
	}
}

func (q *DefaultRefreshQueue) next() (queuedRefreshTask, bool) {
	q.mu.Lock()
	defer q.mu.Unlock()
	if len(q.manual) > 0 {
		item := q.manual[0]
		q.manual = q.manual[1:]
		return item, true
	}
	if len(q.auto) > 0 {
		item := q.auto[0]
		q.auto = q.auto[1:]
		return item, true
	}
	return queuedRefreshTask{}, false
}

func (q *DefaultRefreshQueue) failPending(err error) {
	q.mu.Lock()
	items := append(append([]queuedRefreshTask(nil), q.manual...), q.auto...)
	q.manual = nil
	q.auto = nil
	q.mu.Unlock()
	for _, item := range items {
		item.done <- queuedRefreshResult{err: err}
	}
}

func (q *DefaultRefreshQueue) notify() {
	select {
	case q.wake <- struct{}{}:
	default:
	}
}

type QueuedSavedSkinPriceUpdater struct {
	Queue RefreshQueue
	Kind  RefreshTaskKind
}

func (u QueuedSavedSkinPriceUpdater) Execute(ctx context.Context, params UpdateSavedSkinPriceParams) (UpdateSavedSkinPriceResult, error) {
	if u.Queue == nil {
		return UpdateSavedSkinPriceResult{}, errx.E("skins.update_one.no_queue", errx.CodeInternal, "no refresh queue configured", nil)
	}
	kind := u.Kind
	if kind == "" {
		kind = RefreshTaskManual
	}
	return u.Queue.Enqueue(ctx, RefreshTask{
		MarketHashName: params.MarketHashName,
		Currency:       params.Currency,
		Kind:           kind,
	})
}
