# 文档

## 队列（Queue）

一个服务器有一个或者多个 Queue，用来储存统一类型的 Item。每个 Queue 由一个就绪队列与延迟队列组成。每个Item所有的状态迁移在一个 Queue 中完成。Subscriber 订阅者可以监控感兴趣的 Queue，通过发送 watch 指令。Subscriber 订阅者可以取消监控 Queue，通过发送 ignore 命令。

当一个客户端连接上服务器时，客户端提交 Item 之前会使用 use 命令，那么这些 Item 就存于相应名称的 Queue 中。

当客户端获取 Item 时，Item 会从就绪队列迁移至保留对象，服务器发送给所有订阅客户端后废弃此 Item，并继续从就需队列迁移 Item。


## 命令（Cmd）

二进制协议规则请查阅 public/code.go

### Pub指令说明

Use

选择一个队列入口

    Use <nameLen> <name>
    []byte{CmdUse, uint8, []byte(string)...}

- Use 指令名称，实际值为二进制CmdUse（请查阅public/code.go）8位
- nameLen 队列名称长度，范围0-255，8位
- name 队列名称

Send

插入一个 Item 给队列

     Send <score> <delay> <bodyLen> <body>
     []byte{CmdSend, uint32, uint32, uint32, []byte...}

- Send 指令名称，实际值为二进制 CmdSend（请查阅 public/code.go）8位
- score 优先级，范围 0 - 2^32-1，默认 1024， 32位
- delay 延迟，单位秒，范围 0 - 2^32-1，默认 0，32位
- bodyLen 消息主体长度，范围 0 - 2^32-1，32位
- body 消息主体

### Sub指令说明

Watch

订阅队列

    Watch <nameLen> <name>
    []byte{CmdWatch, uint8, []byte(string)...}

- Watch 指令名称，实际值为二进制 CmdWatch（请查阅 public/code.go）8位
- nameLen 队列名称长度，范围 0 - 255，8位
- name 队列名称

Ignore

取消订阅队列

    Ignore <nameLen> <name>
    []byte{CmdIgnore, uint8, []byte(string)...}

- Ignore 指令名称，实际值为二进制 CmdIgnore（请查阅 public/code.go）8位
- nameLen 队列名称长度，范围 0 - 255，8位
- name 队列名称

Receive

接收订阅 Item

    Receive
    []byte{CmdReceive}

- Receive 指令名称，实际值为二进制 CmdReceive（请查阅 public/code.go）8位

正确响应：

    ResSuccess <nameLen> <name> <bodyLen> <body>

错误响应：

    ResCode

- ResCode 指令名称，实际值为二进制，可根据返回的 ResCode 解析相应错误（请查阅 public/code.go）8位
