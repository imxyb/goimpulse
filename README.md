# goimpulse：高可用，高性能的分布式发号服务

* 通讯协议：http

* 持久化：etcd，强一致性，保证id不重复，不丢失

* 多节点部署，保证只有一个节点running，其余节点standby，running节点挂掉，毫秒级切换其中一个standby节点为running，保证可用性

* node_manager保证内部HA对客户端透明

* 区分不同业务

* 可以提供http basicauth验证

* 单调递增，方便索引存储和查找

* 高性能并且可靠，node启动时候会放置[lastId,batch)个缓存区，当这个缓存区被消费只有batch/2长度时候，会自动扩容到[lastId,batch*2)范围，并把lastid+batch**2持久化到etcd，保证即使节点挂掉，id也不会被重复分发

* 支持优雅重启

* 支持不重启更新配置文件

## 基本架构

![Alt text](http://static.qiziwang.net/8BCBBC07-8E6D-444C-B1C2-AA78CD300E53.png)

## 安装

首先安装etcd，具体参考官网

下载源码或者git clone到你的工作目录

`cd yourproject`

启动node_manger

```./bin node_manager```

启动goimpulse（可以部署多个）

`./bin goimpulse`

## 使用

url地址：

`GET http://node_managerhost/getid`

参数 

type: 业务标识，若传空则为default

返回json 

example:

```json
 {

    "code": 0,

    "id": 19,

    "msg": "success"

}
```

code为0表示成功，id则为结果

## 配置

配置放在`./conf/config.json`中

`etcd.host`：etcd的host和端口，数组形式表示多个

`app.host`：表示goimpulse的host和端口

`node_manager.host`：表示node_manager的host和端口

`type`:业务标识,使用数组表示

`auth.user`：basicauth的用户名

`auth.pass`:basicauth的密码

`auth.enable`：是否开启验证

`batch`：缓存区大小

## 协议

GPL

## 联系

QQ：781028081

