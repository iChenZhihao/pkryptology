## 基于GG20门限签名算法的签名系统

#### 使用方法：

在service文件夹下操作

```bash
# go build -o <file_name>
# 例如：
go build -o gg20-node

# 然后在config.yaml中的zookeeper.servers下配置zookeeper的地址
# 在server.port中配置web服务端口号（也可以在启动时指定）
# 配置完成后即可启动多个节点（以三个节点，门限为2为例）
./gg20-node --server.port=8080 --config=./config.yaml
./gg20-node --server.port=8081 --config=./config.yaml 
./gg20-node --server.port=8082 --config=./config.yaml

# 启动后当节点个数稳定时，会开启分布式密钥生成（Distributed Key Generation, DKG）流程
# DKG有生成大素数等耗时较长的计算步骤，整体流程需耗时约3分钟，请耐心等待
```

DKG完成后，会打印日志信息：`DKG完成，可以开始签名`

```shell
# 之后对任意节点访问接口：
POST http://<node_endpoint>/sign/signMsg
# 请求体（application/json）：
{
    "message" : "<to be signed msg>"
}

# 若能签名成功，则返回信息的响应体如下：
{
    "code": 200,
    "success": true,
    "message": "",
    "data": {
        "V": 1,
        "R": 69428976848040351739625009223318748319646362879949458736725298105720595054233,
        "S": 52800983695523711144751917559406104729347155896511033506091291480583264854450
    }
}
# 基于R与S，还有发送数据中的message，利用DKG结束后生成的公钥即可进行验签操作
# 要获取公钥，可以在Dkg结束后将公钥信息发送出来保存，以便后续验签时直接使用
```



在DKG过程中，向节点发送message请求进行时，会返回服务降级信息：

```json
{
    "code": 400,
    "success": false,
    "message": "签名失败: 签名节点集群暂不可用",
    "data": 0
}
```

新添加了一键启动/关闭脚本
```shell
# 同样以编译出的文件为gg20-node
# 执行：
./nodes.sh start 7
# 即可后台批量启动7个节点，也可以不传节点数，不传时默认节点数为5
# 启动后的端口号为9080, 9081, ... , 9080+node_count-1, 
#      可以自行修改nodes.sh中的DEFAULT_START_PORT参数以调整端口起始值
# 相关日志将输出在./logs/之下

# 执行：
./nodes.sh stop
# 即可一键关闭启动的以gg20-node为名的进程节点
```

