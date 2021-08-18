/**
 * @Author: jinjiaji
 * @Description:
 * @File:  sliding_window
 * @Version: 1.0.0
 * @Date: 2021/8/16 下午5:48
 */

package module5

import (
	"context"
	"sync"
	"sync/atomic"
	"time"
)

//MuxQueriesCounter 默认计数器实现
type MuxQueriesCounter struct {
	// ori bucket数据
	ori        []*[4]int64
	curUnix    int64
	curSecData *[4]int64
	begin      int64
	// ctx context，管理生命周期
	ctx    context.Context
	cancel context.CancelFunc
	mux    sync.RWMutex
}

//NewMuxQueriesCounter 新建默认计数器
func NewMuxQueriesCounter() QueriesCounter {
	d := new(MuxQueriesCounter)
	d.begin = time.Now().Unix()
	d.ori = make([]*[4]int64, 0)
	d.ctx, d.cancel = context.WithCancel(context.Background())
	d.mux = sync.RWMutex{}
	d.run()
	return d
}

//Incr implementations of interface QueriesCounter
func (c *MuxQueriesCounter) Incr(t ResultType) {
	c.mux.Lock()
	defer c.mux.Unlock()
	c.ori[c.curUnix-c.begin][t]++
}

//Count implementations of interface QueriesCounter
func (c *MuxQueriesCounter) Count(t ResultType, secs int) int {
	c.mux.RLock()
	defer c.mux.RUnlock()
	ret := 0
	r := c.curUnix - 1
	l := c.curUnix - int64(secs)
	for l <= r {
		if l-c.begin < 0 {
			l++
			continue
		}
		r := c.ori[l-c.begin]
		ret += int(r[t])
		l++
	}
	return ret
}

//Close implementations of interface QueriesCounter
func (c *MuxQueriesCounter) Close() {
	c.cancel()
}

func (c *MuxQueriesCounter) run() {
	c.genBucket()
	t := time.NewTicker(time.Microsecond * 10000)
	go func() {
		for {
			select {
			case <-c.ctx.Done():
				t.Stop()
				return
			case <-t.C:
				c.genBucket()
			}
		}
	}()
}

func (c *MuxQueriesCounter) genBucket() {
	c.mux.Lock()
	defer c.mux.Unlock()

	nowUnix := time.Now().Unix()
	if c.curUnix != nowUnix || c.curUnix == 0 {
		atomic.StoreInt64(&c.curUnix, nowUnix)
		//初始化未来5秒的缓存
		for i := 0; i < 5; i++ {
			if len(c.ori)-5 < int(c.curUnix-c.begin) {
				c.ori = append(c.ori, &[4]int64{0, 0, 0, 0})
			}
		}
		c.curSecData = c.ori[c.curUnix-c.begin]
	}
}
