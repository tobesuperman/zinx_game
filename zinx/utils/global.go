package utils

import (
	"encoding/json"
	"io/ioutil"
	"zinx/ziface"
)

// 存储一切有关Zinx框架的全局参数，供其他模块使用
// 参数可以通过Json由用户进行配置

type GlobalObj struct {
	TcpServer ziface.IServer // 当前Zinx全局的Server对象
	IP        string         `json:"ip"`   // 当前服务器监听的IP
	Port      int            `json:"port"` // 当前服务器监听的端口号
	Name      string         `json:"name"` // 当前服务器名称

	Version           string `json:"version"`              // 当前Zinx的版本号
	MaxConn           int    `json:"max_conn"`             // 当前服务器允许的最大连接数
	MaxPackageSize    uint32 `json:"max_package_size"`     // 当前Zinx数据包的最大值
	WorkerPoolSize    uint32 `json:"worker_pool_size"`     // 当前业务工作Worker池中Goroutine数量
	MaxWorkerPoolSize uint32 `json:"max_worker_pool_size"` // Zinx框架允许用户最多开辟多少个Goroutine
}

// GlobalObject 对外的全局变量
var GlobalObject *GlobalObj

// 初始化当前的GlobalObject
func init() {
	// 如果配置文件没有加载，默认的值
	GlobalObject = &GlobalObj{
		Name:              "ZinServerApp",
		Version:           "V1.0",
		Port:              8999,
		IP:                "0.0.0.0",
		MaxConn:           1000,
		MaxPackageSize:    4096,
		WorkerPoolSize:    10,   // Worker工作池队列的个数
		MaxWorkerPoolSize: 1024, // 每个Worker对应的消息队列的任务数量最大值
	}

	// 应该尝试从配置文件中去加载一些用户自定义的参数
	GlobalObject.Reload()
}

// Reload 从配置文件中加载用户自定义的参数
func (g *GlobalObj) Reload() {
	data, err := ioutil.ReadFile("conf/zinx.json")
	if err != nil {
		panic(err)
	}
	// 解析Json数据
	err = json.Unmarshal(data, GlobalObject)
	if err != nil {
		panic(err)
	}
}
