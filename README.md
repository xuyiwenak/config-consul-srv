<a title="Hits" target="_blank" href="https://github.com/xuyiwenak/consul-config-push"><img src="https://hits.b3log.org/b3log/hits.svg"></a>
# consul-config-push

单独的推送server推送配置到consul配置中心

## 目录结构  
```
├── LICENSE
├── README.md
├── client.go  // 模拟客户端，用户服务开发注册consul服务，读取配置的demo
├── conf
│   └── micro.yml // 静态配置
├── go.mod
├── go.sum
├── loader.go // 向配置中心推送配置
├── main.go
└── vendor
```  

## 使用解释
1. 方便docker环境编译，构建vendor，不实用docker可以跳过
```
go mod tidy
go mod vendor
```
### 1. 调试consul server
2. 需要启动consul
```
consul agent -dev
// 或者从docker环境启动consul镜像
docker run <consul container name>
```
3. 启动consul server，把conf下面的yml配置上传到consul的kv里
```
go run loader.go main.go
```
如果用docker打包可以构建运行镜像
```
docker build -t consul-config-push .
docker run --rm -d consul-config-push
```
### 2. 调试consul client
另外一个终端启动consul cilent  
```
// 可以根据自己的逻辑修改上传的配置，以及获取配置的方式。
go run client.go
```
## 静态配置
conf 目录下的micro.yml 文件变更会自动检测到执行对应的逻辑




