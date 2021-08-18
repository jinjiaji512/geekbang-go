/**
 * @Author: jinjiaji
 * @Description:
 * @File:  sliding_window
 * @Version: 1.0.0
 * @Date: 2021/8/16 下午5:48
 */

package module5

import (
	"sync"
	"testing"
	"time"
)

//TestDefaultCounter 模拟1000并行下的4亿次计数以及4亿次统计
func TestDefaultCounter(t *testing.T) {
	qc := NewDefaultQueriesCounter()
	defer qc.Close()
	wg := sync.WaitGroup{}

	t.Log("start:", time.Now().Unix())
	n1, n2 := 1000, 100000
	for i := 0; i < n1; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for j := 0; j < n2; j++ {
				qc.Incr(resultTypeSuccess)
				qc.Incr(resultTypeFailure)
				qc.Incr(resultTypeTimeout)
				qc.Incr(resultTypeRejection)
				qc.Count(resultTypeSuccess, 1)
				qc.Count(resultTypeFailure, 1)
				qc.Count(resultTypeTimeout, 1)
				qc.Count(resultTypeRejection, 1)
			}
		}()
	}
	wg.Wait()
	t.Log("end:", time.Now().Unix())

	//统计前100秒的数据，所以暂停2秒
	time.Sleep(time.Second * 2)
	c1 := qc.Count(resultTypeSuccess, 100)
	c2 := qc.Count(resultTypeFailure, 100)
	c3 := qc.Count(resultTypeTimeout, 100)
	c4 := qc.Count(resultTypeRejection, 100)
	t.Log(c1, c2, c3, c4)
	if c1 != n1*n2 || c2 != n1*n2 || c3 != n1*n2 || c4 != n1*n2 {
		t.Fatal("check err")
	}
}

//BenchmarkDefaultCounter 无锁版计数器benchmark
//BenchmarkDefaultCounter-8   	98804888	        11.6 ns/op
func BenchmarkDefaultCounter(b *testing.B) {
	qc := NewDefaultQueriesCounter()
	defer qc.Close()
	for i := 0; i < b.N; i++ {
		qc.Incr(resultTypeSuccess)
		qc.Count(resultTypeSuccess, 1)
	}
}

//BenchmarkMuxCounter 有锁版计数器benchmark
//BenchmarkMuxCounter-8   	22033450	        51.1 ns/op
func BenchmarkMuxCounter(b *testing.B) {
	qc := NewMuxQueriesCounter()
	defer qc.Close()
	for i := 0; i < b.N; i++ {
		qc.Incr(resultTypeSuccess)
		qc.Count(resultTypeSuccess, 1)
	}
}

//BenchmarkTime time组件benchmark
//BenchmarkTime-8   	13139469	        84.9 ns/op
func BenchmarkTime(b *testing.B) {
	for i := 0; i < b.N; i++ {
		time.Now().Unix()
	}
}
