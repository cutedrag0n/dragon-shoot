package dragon

import "net/http"

type response struct {
	_response http.ResponseWriter
}

func newResponser(w http.ResponseWriter) *response {
	return &response{
		_response: w,
	}
}
