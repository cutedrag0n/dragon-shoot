package dragon

const (
	MethodGET = 1 << iota
	MethodPOST
	MethodPUT
	MethodDELETE
)

func (e *engine) GET(path string, handler HandlerFunc) {
	e.r.addRoute("GET", path, handler)
}

func (e *engine) POST(path string, handler HandlerFunc) {
	e.r.addRoute("POST", path, handler)
}

func (e *engine) PUT(path string, handler HandlerFunc) {
	e.r.addRoute("PUT", path, handler)
}

func (e *engine) DELETE(path string, handler HandlerFunc) {
	e.r.addRoute("DELETE", path, handler)
}

func (e *engine) Choices(choices any, path string, handler HandlerFunc) {
	var flag uint16
	switch choices.(type) {
	case []string:
		for _, v := range choices.([]string) {
			switch v {
			case "GET":
				if flag&1 == 0 {
					e.r.addRoute("GET", path, handler)
					flag |= 1
				}
			case "POST":
				if flag&2 == 0 {
					e.r.addRoute("POST", path, handler)
					flag |= 2
				}
			case "PUT":
				if flag&4 == 0 {
					e.r.addRoute("PUT", path, handler)
					flag |= 4
				}
			case "DELETE":
				if flag&8 == 0 {
					e.r.addRoute("DELETE", path, handler)
					flag |= 8
				}
			default:
				continue
			}
		}
	case string:
		method := choices.(string)
		switch method {
		case "*":
			e.Any(path, handler)
		case "GET", "POST", "PUT", "DELETE":
			e.r.addRoute(method, path, handler)
		}
	default:
		return
	}
}

func (e *engine) Any(path string, handler HandlerFunc) {
	var methods = []string{"GET", "POST", "PUT", "DELETE"}
	for _, method := range methods {
		e.r.addRoute(method, path, handler)
	}
}
