// Author: https://github.com/antigloss

// Package queue offers goroutine-safe Queue implementations such as LockfreeQueue.
package queue

import (
	"sync/atomic"
	"unsafe"
)

// NewLockfreeQueue is the only way to get a new, ready-to-use LockfreeQueue.
//
// Example:
//
//   lfq := queue.NewLockfreeQueue()
//   lfq.Push(100)
//   v := lfq.Pop()
func NewLockfreeQueue() *LockfreeQueue {
	var lfq LockfreeQueue
	lfq.head = unsafe.Pointer(&lfq.dummy)
	lfq.tail = lfq.head
	return &lfq
}

// LockfreeQueue is a goroutine-safe Queue implementation.
type LockfreeQueue struct {
	head  unsafe.Pointer
	tail  unsafe.Pointer
	dummy lfqNode
}

// Pop returns (and removes) an element from the front of the queue, or nil if the queue is empty.
func (lfq *LockfreeQueue) Pop() interface{} {
	for {
		h := atomic.LoadPointer(&lfq.head)
		rh := (*lfqNode)(h)
		n := (*lfqNode)(atomic.LoadPointer(&rh.next))
		if n != nil {
			if atomic.CompareAndSwapPointer(&lfq.head, h, rh.next) {
				return n.val
			} else {
				continue
			}
		} else {
			return nil
		}
	}
}

// Push inserts an element to the back of the queue.
func (lfq *LockfreeQueue) Push(val interface{}) {
	node := unsafe.Pointer(&lfqNode{val: val})
	for {
		t := atomic.LoadPointer(&lfq.tail)
		rt := (*lfqNode)(t)
		if atomic.CompareAndSwapPointer(&rt.next, nil, node) {
			// atomic.StorePointer(&lfq.tail, node)
			atomic.CompareAndSwapPointer(&lfq.tail, t, node)
			return
		} else {
			continue
		}
	}
}

type lfqNode struct {
	val  interface{}
	next unsafe.Pointer
}
