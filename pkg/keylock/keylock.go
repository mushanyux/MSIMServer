package keylock

import (
	"sync"
	"sync/atomic"
	"time"
)

const (
	defaultCleanInterval = 24 * time.Hour // 默认24小时清理一次
)

type KeyLock struct {
	locks         map[string]*innerLock // 关键字锁map
	cleanInterval time.Duration         // 定时清除时间间隔
	stopChan      chan struct{}         // 停止信号
	mutex         sync.RWMutex          // 全局读写锁
}

func NewKeyLock() *KeyLock {
	return &KeyLock{
		locks:         make(map[string]*innerLock),
		cleanInterval: defaultCleanInterval,
		stopChan:      make(chan struct{}),
	}
}

func (l *KeyLock) Lock(key string) {
	l.mutex.RLock()
	keyLock, ok := l.locks[key]
	if ok {
		keyLock.add()
	}
	l.mutex.RUnlock()
	if !ok {
		l.mutex.Lock()
		keyLock, ok = l.locks[key]
		if !ok {
			keyLock = newInnerLock()
			l.locks[key] = keyLock
		}
		keyLock.add()
		l.mutex.Unlock()
	}
	keyLock.Lock()
}

func (l *KeyLock) Unlock(key string) {
	l.mutex.RLock()
	keyLock, ok := l.locks[key]
	if ok {
		keyLock.done()
	}
	l.mutex.RUnlock()
	if ok {
		keyLock.Unlock()
	}
}

func (l *KeyLock) Clean() {
	l.mutex.Lock()
	for k, v := range l.locks {
		if v.count == 0 {
			delete(l.locks, k)
		}
	}
	l.mutex.Unlock()
}

func (l *KeyLock) StartCleanLoop() {
	go l.cleanLoop()
}

func (l *KeyLock) StopCleanLoop() {
	close(l.stopChan)
}

func (l *KeyLock) cleanLoop() {
	ticker := time.NewTicker(l.cleanInterval)
	for {
		select {
		case <-ticker.C:
			l.Clean()
		case <-l.stopChan:
			ticker.Stop()
			return
		}
	}
}

type innerLock struct {
	count int64
	sync.Mutex
}

func newInnerLock() *innerLock {
	return &innerLock{}
}

func (l *innerLock) add() {
	atomic.AddInt64(&l.count, 1)
}

func (l *innerLock) done() {
	atomic.AddInt64(&l.count, -1)
}
