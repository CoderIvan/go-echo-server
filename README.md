# go-echo-server

## 客户端

### HTTP
  * 客户端使用HTTP协议POST方向，发送数据
  * Body内容可以为`application/json`，也可以为`text/plain`，建议使用JSON格式，方便解析与统计
  * tagName标记为'http-server'
  * 在服务地址后可以追加`项目名(projectName)`，即`http://localhost/:projectName`
    * 例如`http://localhost/ivan`，则项目名为`ivan`

### UDP
  * 客户端使用UDP发送数据
  * 服务器接收到的二进制内容均使用`UTF-8`解码成文本，建议文本内容使用JSON格式，方便解析与统计
  * tagName标记为'udp-server'
  * 在UDP包前可以追加`$ + 项目名(projectName) + #`，即`${{projectName}}#{{内容}}`
    * 例如`"$ivan#{\"sn\":\"0000000004\",\"iccid\":\"90000000000000000004\",\"imei\":\"900000000000004\",\"random\":\"Some bytes:0.17507057398090065\"}"`，则项目名为`ivan`

### MQTT
  * 客户端使用MQTT发送数据
  * 服务器接收到的二进制内容均使用`UTF-8`解码成文本，建议文本内容使用JSON格式，方便解析与统计
  * tagName标记为'mqtt-server'
  * `client id`最好填sn之类的唯一标识
  * `username`最好填项目名之类的

## 服务器

### win10
执行
```bash
go env -w GO111MODULE=on
```
```bash
go env -w GOPROXY=https://goproxy.io
```
```bash
go build .
```
```bash
.\go-echo-server.exe
```

### 运行
```bash
go run main
```

### 单元测试
```bash
go test go-echo-server/server -v
go test go-echo-server/handler -v
go test go-echo-server/datagram -v
```

### 覆盖率测试
```bash
go test go-echo-server/server -cover
go test go-echo-server/handler -cover
go test go-echo-server/datagram -cover
```

### 基准测试
```bash
go test -benchmem -run=^$ -bench . go-echo-server/server
```
