package dragon

import "net/http"

// Context 请求与响应的上下文对象
type Context struct {
	Request  *request
	Response Response
	handlers []HandlerFunc     // 每一个请求经过的中间件列表集合
	index    int               // 当前执行到的视图函数下标记录
	tag      map[string]string // 保存流转于各个中间件中需要记录的值
}

// newContext 返回一个上下文对象，对请求与响应参数进行封装
// 对当前执行中间件下标进行初始化
func newContext(w http.ResponseWriter, r *http.Request) *Context {
	return &Context{
		Request:  newRequester(r),
		Response: newResponser(w),
		index:    -1,
	}
}
