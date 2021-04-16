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

func process(data datagram.Datagram) *sls.LogGroup {
	Contents := []*sls.LogContent{{
		Key:   proto.String("tagName"),
		Value: proto.String(data.TagName),
	}}

	if len(data.Addr) > 0 {
		Contents = append(Contents, &sls.LogContent{
			Key:   proto.String("addr"),
			Value: proto.String(data.Addr),
		})
	}

	if len(data.ProjectName) > 0 {
		Contents = append(Contents, &sls.LogContent{
			Key:   proto.String("projectName"),
			Value: proto.String(data.ProjectName),
		})
	}

	if len(data.Content) > 0 {
		Contents = append(Contents, &sls.LogContent{
			Key:   proto.String("content"),
			Value: proto.String(data.Content),
		})
	}

	if len(data.ContextID) > 0 {
		Contents = append(Contents, &sls.LogContent{
			Key:   proto.String("contextID"),
			Value: proto.String(data.ContextID),
		})
	}

	if len(data.ExtraInfo) > 0 {
		Contents = append(Contents, &sls.LogContent{
			Key:   proto.String("extraInfo"),
			Value: proto.String(data.ExtraInfo),
		})
	}

	return &sls.LogGroup{
		Logs: []*sls.Log{{
			Time:     proto.Uint32(uint32(data.Time / 1e9)),
			Contents: Contents,
		}},
	}
}

// Handle *
func (l *slsLogger) Handle(data datagram.Datagram) {
	loggroup := process(data)

	if err := l.client.PutLogs(l.projectName, l.logStoreName, loggroup); err != nil {
		fmt.Printf("PutLogs fail, err: %s\n", err)
	}
}
