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


