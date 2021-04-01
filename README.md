# go-echo-server

## 环境准备

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