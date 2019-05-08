lemna
=====
lemna是一个用golang开发的，基于grpc的游戏服务器框架，目标是简化游戏服务器的开发
开发环境搭建
-----------
1. modules安装
> lemna提供go.mod，所以
> 将lemna目录放入非go/src文件夹
> 设置环境变量GO111MODULE=auto或者GO111MODULE=yes
> 在lemna目录下直接go build ./...即可
> 在你的项目中提供go.mod模块，将lemna加入

2. 手动安装
> 将lemna目录放入go/src中
> 使用go build ./...
> 根据提示自行添加缺失的模块
