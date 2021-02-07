package handler

import (
	"fmt"
	"strings"
	"time"

	sls "github.com/aliyun/aliyun-log-go-sdk"
	"google.golang.org/protobuf/proto"
)

/*
 * https://github.com/aliyun/aliyun-log-go-sdk?spm=a2c4g.11186623.2.16.6ee99951caePq0
 */

var (
	accessKeyID     = ""
	accessKeySecret = ""
	endpoint        = "cn-qingdao.log.aliyuncs.com"
	client          = sls.CreateNormalInterface(endpoint, accessKeyID, accessKeySecret, "")
)

// Sls *
type Sls struct {
}

// Handle *
func (l *Sls) Handle(protocol string, addr string, content string, projectName string) {
	logs := []*sls.Log{{
		Time: proto.Uint32(uint32(time.Now().Unix())),
		Contents: []*sls.LogContent{{
			Key:   proto.String("protocol"),
			Value: proto.String(protocol),
		}, {
			Key:   proto.String("addr"),
			Value: proto.String(addr),
		}, {
			Key:   proto.String("content"),
			Value: proto.String(content),
		}, {
			Key:   proto.String("projectName"),
			Value: proto.String(projectName),
		}},
	}}

	/**
	* Source和Topic是有什么用？
	 */
	loggroup := &sls.LogGroup{
		// Topic:  proto.String(""),
		// Source: proto.String("10.230.201.117"),
		Logs: logs,
	}

	for retryTimes := 0; retryTimes < 10; retryTimes++ {
		err := client.PutLogs("docker-alpha", "go-echo-server", loggroup)
		if err == nil {
			fmt.Printf("PutLogs success, retry: %d\n", retryTimes)
			break
		} else {
			fmt.Printf("PutLogs fail, retry: %d, err: %s\n", retryTimes, err)
			//handle exception here, you can add retryable erorrCode, set appropriate put_retry
			if strings.Contains(err.Error(), sls.WRITE_QUOTA_EXCEED) || strings.Contains(err.Error(), sls.PROJECT_QUOTA_EXCEED) || strings.Contains(err.Error(), sls.SHARD_WRITE_QUOTA_EXCEED) {
				//mayby you should split shard
				time.Sleep(1000 * time.Millisecond)
			} else if strings.Contains(err.Error(), sls.INTERNAL_SERVER_ERROR) || strings.Contains(err.Error(), sls.SERVER_BUSY) {
				time.Sleep(200 * time.Millisecond)
			}
		}
	}
}
