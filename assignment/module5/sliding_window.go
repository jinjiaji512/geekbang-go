/**
 * @Author: jinjiaji
 * @Description:
 * @File:  sliding_window
 * @Version: 1.0.0
 * @Date: 2021/8/16 下午5:48
 */

package module5

//ResultType 结果类型
type ResultType int

const (
	resultTypeSuccess ResultType = iota
	resultTypeFailure
	resultTypeTimeout
	resultTypeRejection
)

type (
	//QueriesCounter 请求计数器接口
	QueriesCounter interface {
		//Incr 增加计数
		Incr(t ResultType)
		//Count 统计前secs秒中t类型的请求数
		Count(t ResultType, secs int) int
		//Close 清理资源
		Close()
	}
)
