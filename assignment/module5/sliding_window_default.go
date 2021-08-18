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
	"sync/atomic"
	"time"
)

//DefaultQueriesCounter 默认计数器实现
type DefaultQueriesCounter struct {
	// ori bucket数据
	ori []*[4]int64
	// curUnix 缓存当前秒数，因为time.Now().Unix() 性能差
	curUnix    int64
	curSecData *[4]int64
	// 记录开始unix，与当前unix差值为数组中的位置
	begin int64
	// ctx context，管理生命周期
	ctx    context.Context
	cancel context.CancelFunc
}

//NewDefaultQueriesCounter 新建默认计数器
func NewDefaultQueriesCounter() QueriesCounter {
	d := new(DefaultQueriesCounter)
	d.begin = time.Now().Unix()
	d.ori = make([]*[4]int64, 0)
	d.ctx, d.cancel = context.WithCancel(context.Background())
	d.run()
	return d
}

//Incr implementations of interface QueriesCounter
func (c *DefaultQueriesCounter) Incr(t ResultType) {
	atomic.AddInt64(&(c.curSecData[t]), int64(1))
}

//Count implementations of interface QueriesCounter
func (c *DefaultQueriesCounter) Count(t ResultType, secs int) int {
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
func (c *DefaultQueriesCounter) Close() {
	c.cancel()
}

func (c *DefaultQueriesCounter) run() {
	c.genBucket()
	//每秒钟10000次
	t := time.NewTicker(time.Microsecond * 100)
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

func (c *DefaultQueriesCounter) genBucket() {
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
