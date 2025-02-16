package sched

import (
	"errors"
	"slices"
	"sync"
)

var (
	ErrInvalidBatch = errors.New("invalid batch")
	ErrEmptyItems   = errors.New("empty items")
)

// Task split target into slices and gather them
type Task[T any] struct {
	items  []T
	sched  func(int, []T) ([]T, error)
	pseudo bool
	single bool
}

// NewTask on items.
//
// - pseudo: pseudo run
// - single: disable parallel
func NewTask[T any](items []T, sched func(int, []T) ([]T, error), pseudo, single bool) *Task[T] {
	return &Task[T]{items: items, sched: sched, pseudo: pseudo, single: single}
}

// Collect divide items by batch and send them to sched func parallelly.
func (sc *Task[T]) Collect(batch int, ignoreoutput, ignoreerror bool) ([]T, error) {
	if batch <= 0 {
		return nil, ErrInvalidBatch
	}
	cnt := len(sc.items)
	if cnt == 0 {
		return nil, ErrEmptyItems
	}
	n := cnt / batch
	if n == 0 || batch >= cnt { // only one batch
		if sc.pseudo {
			return sc.items, nil
		}
		return sc.sched(0, sc.items)
	}
	remain := cnt % batch
	wg := sync.WaitGroup{}
	var (
		itemgroups [][]T
		errorgroup Errors
	)
	if !ignoreoutput {
		if remain == 0 {
			itemgroups = make([][]T, n)
		} else {
			itemgroups = make([][]T, n+1)
		}
	}
	if !ignoreerror {
		if remain == 0 {
			errorgroup = make(Errors, n)
		} else {
			errorgroup = make(Errors, n+1)
		}
	}
	iterfn := func(i int) {
		a := i * batch
		b := (i + 1) * batch
		if sc.pseudo {
			if ignoreoutput {
				return
			}
			itemgroups[i] = sc.items[a:b]
			return
		}
		if ignoreoutput {
			if ignoreerror {
				_, _ = sc.sched(i, sc.items[a:b])
			} else {
				_, errorgroup[i] = sc.sched(i, sc.items[a:b])
			}
		} else {
			if ignoreerror {
				itemgroups[i], _ = sc.sched(i, sc.items[a:b])
			} else {
				itemgroups[i], errorgroup[i] = sc.sched(i, sc.items[a:b])
			}
		}
	}
	if !sc.single {
		wg.Add(n)
	}
	for i := 0; i < n; i++ { // full batches
		if !sc.single {
			go func(i int) {
				defer wg.Done()
				iterfn(i)
			}(i)
		} else {
			iterfn(i)
		}
	}
	if remain > 0 {
		if ignoreoutput {
			if !sc.pseudo {
				if ignoreerror {
					_, _ = sc.sched(n, sc.items[n*batch:])
				} else {
					_, errorgroup[n] = sc.sched(n, sc.items[n*batch:])
				}
			}
		} else {
			if sc.pseudo {
				itemgroups[n] = sc.items[n*batch:]
			} else {
				if ignoreerror {
					itemgroups[n], _ = sc.sched(n, sc.items[n*batch:])
				} else {
					itemgroups[n], errorgroup[n] = sc.sched(n, sc.items[n*batch:])
				}
			}
		}
	}
	if !sc.single {
		wg.Wait()
	}
	if !ignoreerror && errorgroup.Error() != "" {
		if ignoreoutput {
			return nil, errorgroup
		}
		return slices.Concat(itemgroups...), errorgroup
	}
	if ignoreoutput {
		return nil, nil
	}
	return slices.Concat(itemgroups...), nil
}
