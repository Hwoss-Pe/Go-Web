package v1

import "net/http"

type HandlerBaseOnMap struct {
	//key是对应的请求方式加路径
	handlers map[string]func(ctx *Context)
}

func (h *HandlerBaseOnMap) ServeHTTP(writer http.ResponseWriter, request *http.Request) {
	key := h.key(request.Method, request.URL.Path)
	if handler, ok := h.handlers[key]; ok {
		//		注册过
		handler(NewContext(writer, request))
	} else {
		writer.WriteHeader(http.StatusNotFound)
		writer.Write([]byte("Not Found"))
	}
	h.handlers[key] = func(ctx *Context) {

	}
}
func (h *HandlerBaseOnMap) key(method string, pattern string) string {
	return method + "#" + pattern
}