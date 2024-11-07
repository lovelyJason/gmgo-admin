# gmgo-admin

  <img align="right" width="320" src="http://cdn.qdovo.com/assets/gmgo.svg">


[![Build Status](http://cdn.qdovo.com/assets/gmgo.svg)](https://github.com/go-admin-team/go-admin)
[![Release](https://img.shields.io/github/release/go-admin-team/go-admin.svg?style=flat-square)](https://github.com/go-admin-team/go-admin/releases)
[![License](https://img.shields.io/github/license/mashape/apistatus.svg)](https://github.com/go-admin-team/go-admin)

基于Gin + MongoDB + Vue + Ant Design的前后端分离权限管理系统,系统初始化极度简单，只需要配置文件中，修改数据库连接，系统支持多指令操作，迁移指令可以让初始化数据库信息变得更简单，服务指令可以很简单的启动api服务

![image](https://github.com/user-attachments/assets/4f376450-1ba8-4443-850d-98449e6d61c8)

![image](https://github.com/user-attachments/assets/d92c77db-9e42-49ab-8c25-3d6eab41a700)



[在线文档](https://www.go-admin.pro): 这是go-admin项目的文档，本项目基于它，但是更换了数据库驱动，所以部分地方会略有不同，请看后文的数据库初始化

[前端项目](https://github.com/go-admin-team/go-admin-ui)

## 写在前面

感谢 [go-admin](https://github.com/go-admin-team/go-admin)项目提供支持，本项目基于go-admin二次开发，是mongodb驱动的后台管理系统

## 🎬 在线体验

antd体验：[https://antd.go-admin.pro](https://antd.go-admin.pro/)
> ⚠️⚠️⚠️ 账号 / 密码： admin / 123456

## ✨ 特性

- 基于 GIN WEB API 框架，提供了丰富的中间件支持（用户认证、跨域、访问日志、追踪ID等）

- 基于Casbin的 RBAC 访问控制模型

- JWT 认证

- 支持 Swagger 文档(基于swaggo)

- 基于 mongo-driver的数据库存储

- 配置文件简单的模型映射，快速能够得到想要的配置

- 代码生成工具

- 表单构建工具

- 多指令模式

- 多租户的支持

- TODO: 单元测试

## 🎁 内置

1. 多租户：系统默认支持多租户，按库分离，一个库一个租户。
1. 用户管理：用户是系统操作者，该功能主要完成系统用户配置。
2. 部门管理：配置系统组织机构（公司、部门、小组），树结构展现支持数据权限。
3. 岗位管理：配置系统用户所属担任职务。
4. 菜单管理：配置系统菜单，操作权限，按钮权限标识，接口权限等。
5. 角色管理：角色菜单权限分配、设置角色按机构进行数据范围权限划分。
6. 字典管理：对系统中经常使用的一些较为固定的数据进行维护。
7. 参数管理：对系统动态配置常用参数。
8. 操作日志：系统正常操作日志记录和查询；系统异常信息日志记录和查询。
9. 登录日志：系统登录日志记录查询包含登录异常。
1. 接口文档：根据业务代码自动生成相关的api接口文档。
1. 代码生成：根据数据表结构生成对应的增删改查相对应业务，全程可视化操作，让基本业务可以零代码实现。
1. 表单构建：自定义页面样式，拖拉拽实现页面布局。
1. 服务监控：查看一些服务器的基本信息。
1. 内容管理：demo功能，下设分类管理、内容管理。可以参考使用方便快速入门。
1. 定时任务：自动化任务，目前支持接口调用和函数调用。

## 准备工作

你需要在本地安装 [go] [gin] [node](http://nodejs.org/) 和 [git](https://git-scm.com/) 

同时配套了系列教程包含视频和文档，如何从下载完成到熟练使用，强烈建议大家先看完这些教程再来实践本项目！！！

### 轻松实现go-admin写出第一个应用 - 文档教程

[步骤一 - 基础内容介绍](https://doc.zhangwj.com/guide/intro/tutorial01.html)

[步骤二 - 实际应用 - 编写增删改查](https://doc.zhangwj.com/guide/intro/tutorial02.html)

### 手把手教你从入门到放弃 - 视频教程

[如何启动go-admin](https://www.bilibili.com/video/BV1z5411x7JG)

[使用生成工具轻松实现业务](https://www.bilibili.com/video/BV1Dg4y1i79D)

[v1.1.0版本代码生成工具-释放双手](https://www.bilibili.com/video/BV1N54y1i71P) [进阶]

[多命令启动方式讲解以及IDE配置](https://www.bilibili.com/video/BV1Fg4y1q7ph)

[go-admin菜单的配置说明](https://www.bilibili.com/video/BV1Wp4y1D715) [必看]

[如何配置菜单信息以及接口信息](https://www.bilibili.com/video/BV1zv411B7nG) [必看]

[go-admin权限配置使用说明](https://www.bilibili.com/video/BV1rt4y197d3) [必看]

[go-admin数据权限使用说明](https://www.bilibili.com/video/BV1LK4y1s71e) [必看]

**如有问题请先看上述使用文档和文章，若不能满足，欢迎 issue 和 pr ，视频教程和文档持续更新中**

## 📦 本地开发

### 环境要求

go 1.21

node版本: v14.16.0

npm版本: 6.14.11

### 开发目录创建

```bash

# 创建开发目录
mkdir goadmin
cd goadmin
```

### 获取代码

```bash
# 获取后端代码
git clone https://github.com/lovelyJason/gmgo-admin.git
# 获取前端代码
git clone https://github.com/lovelyJason/gmgo-antd-ui.git

```

### 启动说明

#### 服务端启动说明

```bash
# 进入 go-admin 后端项目
cd ./go-admin

# 更新整理依赖
go mod tidy

# 编译项目
go build

# 修改配置 
# 文件路径  go-admin/config/settings.yml, 也可以新建settings.dev.yml，这是本地开发配置，不会推送到远程仓库
vi ./config/settings.yml

# 1. 配置文件中修改数据库信息 
# 注意: settings.database 下对应的配置数据
# 2. 确认log路径
```

#### 初始化数据库，以及服务启动

``` bash
# 首次配置需要初始化数据库资源信息
复制config/db.js里面的代码到mongo的shell环境内执行，或者数据库连接工具如navicat里面新建查询，执行数据库初始化操作


# 启动项目，也可以用IDE进行调试
# macOS or linux 下使用
$ ./go-admin server -c config/settings.yml # 或者settings.dev.yml


# ⚠️注意:windows 下使用
$ go-admin.exe server -c config/settings.yml
```

![image](https://github.com/user-attachments/assets/ddf9bff3-9ceb-44f8-89f1-77cbcfc54dd8)


#### sys_api 表的数据如何添加

在项目启动时，使用`-a true` 系统会自动添加缺少的接口数据
```bash
./go-admin server -c config/settings.yml -a true
```

#### 使用docker 编译启动

```shell
# 编译镜像
docker build -t go-admin .

# 启动容器，第一个go-admin是容器名字，第二个go-admin是镜像名称
# -v 映射配置文件 本地路径：容器路径
docker run --name go-admin -p 8000:8000 -v /config/settings.yml:/config/settings.yml -d go-admin-server
```

#### 文档生成

```bash
go generate
```

#### 交叉编译

```bash
# windows
env GOOS=windows GOARCH=amd64 go build main.go

# or
# linux
env GOOS=linux GOARCH=amd64 go build main.go
```

### UI交互端启动说明

这是一个标准的vben5.2的项目，参考vben官网即可：https://doc.vben.pro/guide/introduction/quick-start.html

```bash
# 安装依赖
pnpm i

# 启动服务
pnpm run dev:antd
```


## 🤝 链接

[Go开发者成长线路图](http://www.golangroadmap.com/)

## 🔑 License

[MIT](https://github.com/go-admin-team/go-admin/blob/master/LICENSE.md)

Copyright (c) 2024 wenjianzhang
