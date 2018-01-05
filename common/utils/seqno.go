package common

import (
	"fmt"
	"sync/atomic"
)

type AtomicSeqno int64

func NewAtomicSeqno(v int64) *AtomicSeqno {
	a := AtomicSeqno(v)
	return &a
}

func (a *AtomicSeqno) Get() int64 {
	return int64(*a)
}

func (a *AtomicSeqno) Set(v int64) {
	atomic.StoreInt64((*int64)(a), v)
}

func (a *AtomicSeqno) CompareAndSet(expect, update int64) bool {
	return atomic.CompareAndSwapInt64((*int64)(a), expect, update)
}

func (a *AtomicSeqno) GetAndInc() int64 {
	for {
		cur := a.Get()
		next := cur + 1
		if a.CompareAndSet(cur, next) {
			return cur
		}
	}
}

func (a *AtomicSeqno) GetAndDec() int64 {
	for {
		cur := a.Get()
		next := cur - 1
		if a.CompareAndSet(cur, next) {
			return cur
		}
	}
}

func (a *AtomicSeqno) AddAndGet(delta int64) int64 {
	for {
		cur := a.Get()
		next := cur + delta
		if a.CompareAndSet(cur, next) {
			return next
		}
	}
}

func (a *AtomicSeqno) String() string {
	return fmt.Sprintf("{AtomicSeqno:%d}", a.Get())
}
