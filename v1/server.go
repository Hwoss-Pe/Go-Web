package v1

import (
	"fmt"
	"net/http"
)

type Server interface {
	//	对于一个服务器，接口定义来说肯定会有路由和启动的作用
	//handleFunc func(ctx *Context) 这个是为了在服务器内部进行封装上下文
	Route(method string, pattern string, handleFunc func(ctx *Context))

	Start(address string) error
}

// 启动类显然要去实现上面的接口
type sdkHttpServer struct {
	Name    string
	handler *HandlerBaseOnMap
}

func (s *sdkHttpServer) Route(
	method string, pattern string,
	handleFunc func(ctx *Context)) {
	key := s.handler.key(method, pattern)
	s.handler.handlers[key] = handleFunc
}

func (s *sdkHttpServer) Start(address string) error {

	return http.ListenAndServe(address, nil)
}
func NewServer(name string) Server {
	return &sdkHttpServer{
		Name: name,
	}
}

// SignUp 下面可以发现我只是读取json数据 并返回的操作就需要一堆校验，因此还可以抽象出Context
func SignUp(ctx *Context) {
	req := &signUpReq{}
	err := ctx.ReadJson(req)
	if err != nil {
		fmt.Fprintf(ctx.W, "read body failed: %v", err)
		// 要返回掉，不然就会继续执行后面的代码
		return
	}
	err = ctx.writeJson(http.StatusOK, req)
	//记日志
	if err != nil {
		fmt.Printf("写入响应失败 err： %v", err)
	}

}

type signUpReq struct {
	//这里可以用标签来规范传输中json的格式
	Email             string `json:"email"`
	Password          string `json:"password"`
	ConfirmedPassword string `json:"confirmed_password"`
}

type commonResponse struct {
	BizCode int         `json:"biz_code"`
	Msg     string      `json:"msg"`
	Data    interface{} `json:"data"`
}

func NewContext(w http.ResponseWriter, r *http.Request) *Context {
	return &Context{
		W: w,
		R: r,
	}
}

// 工厂模式的创建
//type Factory func() Server
//
//var factory Factory
//
//func Register(f Factory) {
//	factory = f
//}
//func NewServerByFactory() Server {
//	return factory()
//}