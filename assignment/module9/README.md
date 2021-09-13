### Week09 作业题目：
1. 总结几种 socket 粘包的解包方式。 尝试举例其应用
2. 实现一个从 socket connection 中解码出 goim 协议的解码器。

### 一 应用层分包问题
#### 1 固定长度消息
> 例如每xxx 个字节代表一个整包消息，不足的前面补位。解码器在处理这类定常消息的时候比较简单，每次读到指定长度的字节后再进行解码；
##### 案例：
1. 

#### 2 基于分隔符
> 通过特定的分隔符区分整包消息
##### 案例：
1. http 中将\r\n作为URL method header body 等的分隔符。其次Transfer-Encoding: chunked模式下，\r\n也作为消息块的分隔符。
2. ftp 每个命令最后都以 "\r\n"结尾，FTP命令在控制连接中传输。

#### 3 
> 通过在协议头/消息头中设置长度字段来标识整包消息
##### 案例：
1. http中设置Content-Length字段来标识消息长度
2. websocket协议中设置Payload length 标识位
3. TLV编码格式中的L

#### 4 基于协议号
> 通过协议号来确定解析协议的长度，由一个唯一协议标记就可以确定一个反序列化类，从而实现区分整包
##### 案例：
1. http2中的二进制帧设计
2. 

### 二 实现goim协议的解码器。
> 有参考goim protocol包
>https://github.com/Terry-Mao/goim/tree/master/api/protocol

## 参考

– [1] [知乎：关于应用层解决拆包粘包问题？](https://www.zhihu.com/question/37023914)

– [2] [TLV编码格式详解](https://zhuanlan.zhihu.com/p/62317518)

– [3] [Fundebug：用了这么久HTTP, 你是否了解Content-Length?](https://cloud.tencent.com/developer/article/1501751)

– [3] [HTTP/2协议“多路复用”实现原理](https://segmentfault.com/a/1190000016975064)

> 