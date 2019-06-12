package main

import (
	"bytes"
	"fmt"
	"github.com/micro/go-config"
	"github.com/micro/go-config/encoder/json"
	"github.com/micro/go-config/source"
	"github.com/micro/go-config/source/file"
	"github.com/micro/go-log"
	"io/ioutil"
	"net/http"
	"path/filepath"
	"strings"
	"sync"
)

var (
	configData []byte
	jsdata     string
	m          sync.RWMutex
	inited     bool
	err        error
)

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
	url := "http://127.0.0.1:8500/v1/kv/micro/config/cluster"

	req, err := http.NewRequest("PUT", url, strings.NewReader(string(conf.Bytes())))
	if err != nil {
		panic(err)
	}
	req.Header.Add("Content-Type", "application/json")
	res, err := http.DefaultClient.Do(req)
	defer func() {
		if res != nil && res.Body != nil {
			res.Body.Close()
		}
	}()
	body, _ := ioutil.ReadAll(res.Body)
	fmt.Println(res)
	fmt.Println(string(body))
	// 侦听文件变动
	watcher, err := conf.Watch()
	if err != nil {
		log.Fatalf("[Init] 开始侦听应用配置文件变动 异常，%s", err)
		panic(err)
	}

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
