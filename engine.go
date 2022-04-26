package dragon

import (
	"net/http"
	"sync"
)

var (
	newOnce sync.Once
	e       *engine = &engine{
		// 服务默认端口是7776
		// dragon 6 ^ 5
		host: "127.0.0.1:7776",
	}
)

// engine 是后端服务框架引擎模块，作为服务启动的入口
type engine struct {
	host string
	r    *router
}

// New 方法实现了返回单例模式下的引擎对象接口
func New() Engine {
	newOnce.Do(func() {
		e.r = newRouter()
	})
	return e
}

// 要求引擎对象必须实现ServeHTTP方法，否则编译报错
var _ http.Handler = (*engine)(nil)

// ServeHTTP 方法将请求与响应封装成了上下文对象，方便后续的操作
// 也是作为请求获取，并且进行路由搜寻的抓手入口方法
func (e *engine) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	//context := newContext(w, r)
}

// Run 启动引擎
func (e *engine) Run() {
	err := http.ListenAndServe(e.host, e)
	if err != nil {
		panic(err)
	}
}
