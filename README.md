# Golang IM

使用 Golang 基于 TCP 实现一个多人 IM System，架构如图所示：

<img src="https://s2.loli.net/2023/09/11/cL8eAdMtYNCGXag.png" width="500" alt="design">

## 使用

先开启服务，再开启客户端

```shell
./server
```

```shell
./client
```

> 可能你的电脑执行文件会报错，因为这两个可执行文件是基于 M 系列 macOS 编译创建的，你可以自己打包。

## 编译可执行文件

打包服务部分

```shell
go build -o server server.go user.go  main.go
```

打包客户端部分

```shell
go build -o client client.go
```

## Server 功能描述

完成基本的登录与通信

<img src="https://s2.loli.net/2023/09/11/iwzWsgJjKO6bZVo.png" width="500" alt="message">

查看在线列表

<img src="https://s2.loli.net/2023/09/11/lYhrqXw6PTDceQN.png" width="500" alt="user list">

修改用户名称

<img src="https://s2.loli.net/2023/09/11/BpxeWyM812qlmz5.png" width="500" alt="rename">

不活跃用户，超 5 分钟强踢

<img src="https://s2.loli.net/2023/09/11/DcgZYlh6UOWHRV2.png" width="500" alt="timeout">

私聊用户

<img src="https://s2.loli.net/2023/09/11/KpNoke26ODg9UIa.png" width="500" alt="oneToOne">