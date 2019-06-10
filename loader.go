package main

import (
	"github.com/micro/go-config"
	"github.com/micro/go-config/source/file"
	"github.com/micro/go-log"
	"os"
	"path/filepath"
	"sync"
)

var (
	configData []byte
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

	pt := filepath.Join(appPath, "conf")
	os.Chdir(appPath)

	// 加载micro.yml文件
	if err = config.Load(file.NewSource(file.WithPath(pt + "/micro.yml"))); err != nil {
		panic(err)
	}
	// 侦听文件变动
	watcher, err := config.Watch()
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
			if err := config.Get(appPath, "micro").Scan(&configData); err != nil {
				panic(err)
			}
			log.Log(string(configData))
			log.Logf("[loadAndWatchConfigFile] 文件变动，%s", string(v.Bytes()))
		}
	}()
	log.Log(string(configData))
	// 标记已经初始化
	inited = true
	return
}
