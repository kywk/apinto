package router_http

import (
	http_service "github.com/eolinker/eosc/http-service"
	"github.com/eolinker/goku/router"
	"github.com/eolinker/goku/service"
)

//IMatcher IMatcher接口实现了Match方法：根据http请求返回服务接口
type IMatcher interface {
	Match(req http_service.IRequestReader) (service.IService, router.IEndPoint, bool)
}

//Matcher Matcher结构体，实现了根据请求返回服务接口的方法
type Matcher struct {
	r        router.IRouter
	services map[string]service.IService
}

//Match 对http请求进行路由匹配，并返回服务
func (m *Matcher) Match(req http_service.IRequestReader) (service.IService, router.IEndPoint, bool) {

	sources := newHTTPSources(req)
	endpoint, has := m.r.Router(sources)
	if !has {
		return nil, nil, false
	}

	s, has := m.services[endpoint.Target()]

	return s, endpoint, has
}

//HTTPSources 封装http请求的结构体
type HTTPSources struct {
	req http_service.IRequestReader
}

func newHTTPSources(req http_service.IRequestReader) *HTTPSources {
	return &HTTPSources{req: req}
}

//Get 由传入的指标key来获取请求中的指标值
func (h *HTTPSources) Get(cmd string) (string, bool) {
	if isHost(cmd) {
		return h.req.Host(), true
	}
	if isMethod(cmd) {
		return h.req.Method(), true
	}

	u := h.req.URL()
	if isLocation(cmd) {
		return u.Path, true
	}

	if hn, yes := headerName(cmd); yes {
		vs := h.req.GetHeader(hn)
		if len(vs) == 0 {
			return "", true
		}
	}

	if qn, yes := queryName(cmd); yes {
		query := u.Query()
		vs := query.Get(qn)
		if len(vs) == 0 {
			return "", true
		}
		return vs, true
	}
	return "", false
}
