PROJECT:=gmgo-admin

.PHONY: build
build:
	CGO_ENABLED=0 go build -ldflags="-w -s" -a -installsuffix "" -o gmgo-admin .

# make build-linux
build-linux:
	@docker build -t gmgo-admin:latest .
	@echo "build successful"

build-sqlite:
	go build -tags sqlite3 -ldflags="-w -s" -a -installsuffix -o gmgo-admin .

# make run
run:
    # delete gmgo-admin-api container
	@if [ $(shell docker ps -aq --filter name=gmgo-admin --filter publish=8000) ]; then docker rm -f gmgo-admin; fi

    # 启动方法一 run gmgo-admin-api container  docker-compose 启动方式
    # 进入到项目根目录 执行 make run 命令
	@docker-compose up -d

	# 启动方式二 docker run  这里注意-v挂载的宿主机的地址改为部署时的实际决对路径
    #@docker run --name=gmgo-admin -p 8000:8000 -v /home/code/go/src/gmgo-admin/gmgo-admin/config:/gmgo-admin-api/config  -v /home/code/go/src/go-admin/go-admin-api/static:/go-admin/static -v /home/code/go/src/go-admin/go-admin/temp:/go-admin-api/temp -d --restart=always go-admin:latest

	@echo "gmgo-admin service is running..."

	# delete Tag=<none> 的镜像
	@docker image prune -f
	@docker ps -a | grep "gmgo-admin"

stop:
    # delete go-admin-api container
	@if [ $(shell docker ps -aq --filter name=gmgo-admin --filter publish=8000) ]; then docker-compose down; fi
	#@if [ $(shell docker ps -aq --filter name=gmgo-admin --filter publish=8000) ]; then docker rm -f gmgo-admin; fi
	#@echo "gmgo-admin stop success"


#.PHONY: test
#test:
#	go test -v ./... -cover

#.PHONY: docker
#docker:
#	docker build . -t go-admin:latest

# make deploy
deploy:

	#@git checkout master
	#@git pull origin master
	make build-linux
	make run
