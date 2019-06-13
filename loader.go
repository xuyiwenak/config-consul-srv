package main

import (
	"bytes"
	"fmt"
	"github.com/micro/go-micro/config"
	"github.com/micro/go-micro/config/encoder/json"
	"github.com/micro/go-micro/config/source"
	"github.com/micro/go-micro/config/source/file"
	"io/ioutil"
	"net/http"
	"path/filepath"
	"strconv"
	"sync"
)

var (
	m                      sync.RWMutex
	inited                 bool
	err                    error
	consulAddr             consulConfig
	consulConfigCenterAddr string
)

// consulConfig 配置结构
type consulConfig struct {
	Enabled    bool   `json:"enabled"`
	Host       string `json:"host"`
	Port       int    `json:"port"`
	KVLocation string `json:"kv_location"`
}

// Init 初始化配置
func Init() {

	m.Lock()
	defer m.Unlock()
	if inited {
		log.Logf("[Init] 配置已经初始化过")
		return
	}

	// 加载yml默认配置
	// 先加载基础配置
	appPath, _ := filepath.Abs(filepath.Dir(filepath.Join("./", string(filepath.Separator))))

	e := json.NewEncoder()
	fmt.Println(appPath)
	fileSource := file.NewSource(
		file.WithPath(appPath+"/conf/micro.yml"),
		source.WithEncoder(e),
	)
	conf := config.NewConfig()
	// 加载micro.yml文件
	if err = conf.Load(fileSource); err != nil {
		panic(err)
	}
	fmt.Println(string(conf.Bytes()))

	// 读取连接的配置中心
	configMap := conf.Map()
	fmt.Println(configMap)
	if err := conf.Get("micro", "consul").Scan(&consulAddr); err != nil {
		panic(err)
	}
	// 拼接配置的地址和 KVcenter 存储路径
	consulConfigCenterAddr = consulAddr.Host + ":" + strconv.Itoa(consulAddr.Port)
	url := fmt.Sprintf("http://%s/v1/kv/%s", consulConfigCenterAddr, consulAddr.KVLocation)
	_, err, _ := PutJson(url, string(conf.Bytes()))
	if err != nil {
		log.Fatalf("http 发送模块异常，%s", err)
		panic(err)
	}
	// 侦听文件变动
	watcher, err := conf.Watch()
	if err != nil {
		log.Fatalf("[Init] 开始侦听应用配置文件变动 异常，%s", err)
		panic(err)
	}

	fmt.Println(consulConfigCenterAddr)
	go func() {
		for {
			v, err := watcher.Next()
			if err != nil {
				log.Fatalf("[loadAndWatchConfigFile] 侦听应用配置文件变动 异常， %s", err)
				return
			}
			if err = conf.Load(fileSource); err != nil {
				panic(err)
			}
			log.Logf("[loadAndWatchConfigFile] 文件变动，%s", string(v.Bytes()))
		}
	}()
	// 标记已经初始化
	inited = true
	return
}
func PutJson(url, body string) (ret string, err error, resp *http.Response) {
	buf := bytes.NewBufferString(body)
	req, err := http.NewRequest("PUT", url, buf)
	if err != nil {
		panic(err)
		return
	}
	req.Header.Add("Content-Type", "application/json")
	resp, err = http.DefaultClient.Do(req)
	defer func() {
		if resp != nil && resp.Body != nil {
			resp.Body.Close()
		}
	}()
	if err != nil {
		log.Log(err.Error())
		return "", err, resp
	}
	data, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", err, resp
	}

	return string(data), nil, resp
}
