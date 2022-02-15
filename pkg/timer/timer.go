package timer

import (
	"sync"
	"time"
)

type BookTimer struct {
	ticking bool
	tick    *time.Timer
	lock    sync.Mutex
}

func (b *BookTimer) IsTicking() bool {
	b.lock.Lock()
	defer b.lock.Unlock()
	return b.ticking
}

func (b *BookTimer) StartTimer(d time.Duration) {
	go startOnce.Do(func() {
		b.ticking = true
		b.tick = time.NewTimer(d)
	loop:
		for {
			select {
			case <-b.tick.C:
				break loop
			}
		}
		//将修改队列中做的修改应用到数据库，然后清空缓存
		//models.ApplyModifyToDataBase()
		//models.ClearCache()
		b.ticking = false
		startOnce = sync.Once{}
	})
}

func (b *BookTimer) Reset(d time.Duration) {
	b.lock.Lock()
	defer b.lock.Unlock()
	b.tick.Reset(d)
}

var (
	BTimer    = &BookTimer{}
	startOnce = sync.Once{}
)
