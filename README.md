# go-echo-server

## 客户端

### HTTP
  * 客户端使用HTTP协议POST方向，发送数据
  * Body内容可以为`application/json`，也可以为`text/plain`，建议使用JSON格式，方便解析与统计
  * tagName标记为`http-server`
  * 在服务地址后可以追加`项目名(projectName)`，即`http://localhost/:projectName`
    * 例如`http://localhost/ivan`，则项目名为`ivan`

### UDP
  * 客户端使用UDP发送数据
  * 服务器接收到的二进制内容均使用`UTF-8`解码成文本，建议文本内容使用JSON格式，方便解析与统计
  * tagName标记为`udp-server`
  * 在UDP包前可以追加`$ + 项目名(projectName) + #`，即`${{projectName}}#{{内容}}`
    * 例如`"$ivan#{\"sn\":\"0000000004\",\"iccid\":\"90000000000000000004\",\"imei\":\"900000000000004\",\"random\":\"Some bytes:0.17507057398090065\"}"`，则项目名为`ivan`

### MQTT
  * 客户端使用MQTT发送数据
  * 服务器接收到的二进制内容均使用`UTF-8`解码成文本，建议文本内容使用JSON格式，方便解析与统计
  * tagName标记为`mqtt-server`
  * `client id`最好填sn之类的唯一标识，最终显示为`contextID`
  * `username`最好填项目名之类的，最终显示为`projectName`

### TCP
  * 客户端使用TCP发送数据
  * 报文内容，为`两字节报长度`+`报文内容`
    * 如发送`48 65 6c 6c 6f 20 57 6f 72 6c 64`，则要在前面加上长度11，即`00 0b`，最终为`00 0b 48 65 6c 6c 6f 20 57 6f 72 6c 64`
  * 服务器接收到的二进制内容均使用`UTF-8`解码成文本，建议文本内容使用JSON格式，方便解析与统计
  * tagName标记为`tcp-server`
  * 如果第一个包使用`$` + `项目名(projectName)` + `#` + `上下文ID(contextID)` + `#`
    * 例如：`$BigProject#Ivan#`，则项目名为`BigProject`，上下文ID为`Ivan`
  * `上下文ID`建议使用SN之类的唯一标识

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
```

### 覆盖率测试
```bash
go test go-echo-server/server -cover
go test go-echo-server/handler -cover
```

### 基准测试
```bash
go test -benchmem -run=^$ -bench . go-echo-server/server
```
