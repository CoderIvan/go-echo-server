FROM golang:alpine as build-stage
ENV GO111MODULE=on
ENV GOPROXY=https://goproxy.io
WORKDIR /app
COPY . .
RUN go build .

FROM alpine
WORKDIR /app
RUN apk add tzdata --update --no-cache \
	&& cp /usr/share/zoneinfo/Asia/Shanghai /etc/localtime \
	&& echo "Asia/Shanghai" /etc/localtime \
	&& apk del tzdata
COPY --from=build-stage /app/go-echo-server /app
COPY config.yaml /app
EXPOSE 80

ENTRYPOINT ["./go-echo-server"]