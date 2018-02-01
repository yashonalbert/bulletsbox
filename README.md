# bulletsbox

## 简介（Introduction）

bulletsbox 是一个简单的消息队列，其协议基于二进制编码运行在 tcp 上。客户端连接服务器并发送指令和数据，然后等待响应并关闭连接。对于每个连接，服务器按照接收命令的序列依次处理并响应。

## 安装（Install）

```bash
$ go get github.com/yashonalbert/bulletsbox
```

## 例子（[Example](https://github.com/yashonalbert/bulletsbox/blob/master/example.md)）

## 文档（[Document](https://github.com/yashonalbert/bulletsbox/blob/master/document.md)）

## TODO

- 队列线程安全
- 任务休眠 / 唤醒
- 客户端断线重连
- 服务端队列持久化
- 查看监听队列详细信息
