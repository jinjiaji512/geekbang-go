### Week05 作业题目：

1. 参考 Hystrix 实现一个滑动窗口计数器。
### 实现方案：
1. 每秒钟一个bucket
2. 缓存time时间，异步生成bucket
3. atomic原子操作进行++
### 总结：
1. 此案例下性能对比，atomic原子锁 > sync锁 > channel同步
2. time.Now()性能差，业务允许的情况下可以异步刷新缓存的方式优化
3. map,++等操作在并行情况下有安全问题
4. 后续可以尝试每秒多个，按随机数分到多个bucket，是否会提升atomic性能