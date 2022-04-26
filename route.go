package dragon

import "strings"

// router 作为路由匹配的主要核心对象
// 提供了前缀树路由查找（静态、动态、正则）的功能
type router struct {
	roots    map[string]*node
	handlers map[string]HandlerFunc
}

// newRouter 引擎生成时会自动调用该函数生成路由系统
func newRouter() *router {
	return &router{
		// 字典中的默认槽数量为5，默认认为存在四种常用请求方法
		// GET/POST/PUT/DELETE
		roots:    make(map[string]*node, 4),
		handlers: make(map[string]HandlerFunc, 5),
	}
}

// addRoute 该方法实现了对路由地址的添加操作
// 通过调用trie.go中的正则路由解析函数对路径进行解析，然后按照请求方法映射到哈希表中；
// 主要的路径信息录入到路由树中，方便查找时能够动态匹配准确路由
func (r *router) addRoute(method, path string, handler HandlerFunc) {
	path, parts := regexParsePath(path)
	key := stringJoin(method, "-", path)
	if _, ok := r.roots[method]; !ok {
		r.roots[method] = &node{}
	}

	r.roots[method].regexInsert(path, parts, 0)
	r.handlers[key] = handler
}

// getRoute 方法完成了对请求的路径进行路由查找功能
// 方法先对请求的路径按照'/'来进行字符串的切割，然后
func (r *router) getRoute(method, path string) (*node, map[string]string) {
	searchParts := parsePath(path)
	params := make(map[string]string, 3)
	root, ok := r.roots[method]
	if !ok {
		return nil, nil
	}

	// 通过对用户请求的路径地址进行正则路由的搜寻操作
	n := root.regexSearch(searchParts, 0)
	if n != nil {
		_, parts := regexParsePath(n.path)
		for i, part := range parts {
			if part.part[0] == ':' {
				params[part.part[1:]] = searchParts[i]
			}

			if part.part[0] == '*' && len(part.part) > 1 {
				params[part.part[1:]] = strings.Join(searchParts[i:], "/")
				break
			}
		}
		return n, params
	}
	return nil, nil
}
