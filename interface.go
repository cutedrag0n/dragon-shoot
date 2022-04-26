package dragon

type HandlerFunc func(ctx *Context)

type D map[string]any

type Engine interface {
	Run()

	GET(string, HandlerFunc)
	POST(string, HandlerFunc)
	PUT(string, HandlerFunc)
	DELETE(string, HandlerFunc)
	Choices([]string, string, HandlerFunc)
}

type Response interface {
}
