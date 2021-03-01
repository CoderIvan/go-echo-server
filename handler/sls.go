package handler

import (
	"fmt"
	"time"

	sls "github.com/aliyun/aliyun-log-go-sdk"
	"google.golang.org/protobuf/proto"
)

/*
 * https://github.com/aliyun/aliyun-log-go-sdk?spm=a2c4g.11186623.2.16.6ee99951caePq0
 */

// SLS *
type SLS struct {
	client       sls.ClientInterface
	projectName  string
	logStoreName string
}

// CreateSLS *
func CreateSLS(accessKeyID string, accessKeySecret string, endpoint string, projectName string, logStoreName string) *SLS {
	return &SLS{
		client:       sls.CreateNormalInterface(endpoint, accessKeyID, accessKeySecret, ""),
		projectName:  projectName,
		logStoreName: logStoreName,
	}
}

// Handle *
func (l *SLS) Handle(protocol string, addr string, content string, projectName string) {
	Contents := []*sls.LogContent{{
		Key:   proto.String("protocol"),
		Value: proto.String(protocol),
	}, {
		Key:   proto.String("content"),
		Value: proto.String(content),
	}}

	if len(addr) > 0 {
		Contents = append(Contents, &sls.LogContent{
			Key:   proto.String("addr"),
			Value: proto.String(addr),
		})
	}

	if len(projectName) > 0 {
		Contents = append(Contents, &sls.LogContent{
			Key:   proto.String("projectName"),
			Value: proto.String(projectName),
		})
	}

	loggroup := &sls.LogGroup{
		Logs: []*sls.Log{{
			Time:     proto.Uint32(uint32(time.Now().Unix())),
			Contents: Contents,
		}},
	}

	err := l.client.PutLogs(l.projectName, l.logStoreName, loggroup)
	if err != nil {
		fmt.Printf("PutLogs fail, err: %s\n", err)
	}
}
