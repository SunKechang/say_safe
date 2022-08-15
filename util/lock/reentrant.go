package lock

import "sync/atomic"

var (
	REENTRANT_NONE uint32 = 0
	REENTRANT_LOCK uint32 = 1
)

type ReentrantLock struct {
	state uint32
}

func (rl *ReentrantLock) Init() {
	atomic.StoreUint32(&rl.state, REENTRANT_NONE)
}

func (rl *ReentrantLock) Lock() bool {
	return atomic.CompareAndSwapUint32(&rl.state, REENTRANT_NONE, REENTRANT_LOCK)
}

func (rl *ReentrantLock) Free() bool {
	return atomic.CompareAndSwapUint32(&rl.state, REENTRANT_LOCK, REENTRANT_NONE)
}
