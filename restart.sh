#!/bin/bash
echo "go build"
go mod tidy
go build -o gmgo-admin main.go
chmod +x ./gmgo-admin
echo "kill gmgo-admin service"
killall gmgo-admin # kill go-admin service
nohup ./gmgo-admin server -c=config/settings.dev.yml >> access.log 2>&1 & #后台启动服务将日志写入access.log文件
echo "run gmgo-admin success"
ps -aux | grep gmgo-admin
