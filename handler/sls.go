package handler

import (
	"fmt"
	"go-echo-server/datagram"

	sls "github.com/aliyun/aliyun-log-go-sdk"
	"google.golang.org/protobuf/proto"
)

/*
 * https://github.com/aliyun/aliyun-log-go-sdk?spm=a2c4g.11186623.2.16.6ee99951caePq0
 */

// SLS *
type slsLogger struct {
	client       sls.ClientInterface
	projectName  string
	logStoreName string
}

// CreateSLS *
func NewSLS(
	accessKeyID string,
	accessKeySecret string,
	endpoint string,
	projectName string,
	logStoreName string,
) *slsLogger {
	return &slsLogger{
		client:       sls.CreateNormalInterface(endpoint, accessKeyID, accessKeySecret, ""),
		projectName:  projectName,
		logStoreName: logStoreName,
	}
}

// Handle *
func (l *slsLogger) Handle(datagram datagram.Datagram) {
	Contents := []*sls.LogContent{{
		Key:   proto.String("tagName"),
		Value: proto.String(datagram.TagName),
	}, {
		Key:   proto.String("content"),
		Value: proto.String(datagram.Content),
	}}

	if len(datagram.Addr) > 0 {
		Contents = append(Contents, &sls.LogContent{
			Key:   proto.String("addr"),
			Value: proto.String(datagram.Addr),
		})
	}

	if len(datagram.ProjectName) > 0 {
		Contents = append(Contents, &sls.LogContent{
			Key:   proto.String("projectName"),
			Value: proto.String(datagram.ProjectName),
		})
	}

	loggroup := &sls.LogGroup{
		Logs: []*sls.Log{{
			Time:     proto.Uint32(uint32(datagram.Time / 1e9)),
			Contents: Contents,
		}},
	}

	err := l.client.PutLogs(l.projectName, l.logStoreName, loggroup)
	if err != nil {
		fmt.Printf("PutLogs fail, err: %s\n", err)
	}
}
