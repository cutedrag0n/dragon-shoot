package dragon

import "net/http"

type request struct {
	_request *http.Request
	info
}

type info struct {
}

func newRequester(r *http.Request) *request {
	return &request{
		_request: r,
	}
}
