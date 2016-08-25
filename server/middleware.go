package server

import "net/http"

// PreMiddleware is an interface defining the middleware for the server in pre-stage, the middleware should call the next handler to pass the
// request down, or just return a HttpRedirect request and etc.
type PreMiddleware interface {
	ServeHttp(w http.ResponseWriter, r *http.Request, pathParams Params, next RawHandler)
}

// JsonMiddleware is an interface defining the middleware in late-stage, which means later than PreMiddleware
// the middleware should call the next handler to pass the request down, or just return a response which will be serve back as JSON,
// like the normal JSON handler does
type JsonMiddleware interface {
	ServeHttp(queryParams, pathParam Params, next Handler) (*Response, error)
}

type MiddleFn func(queryParams, pathParam Params, next Handler) (*Response, error)

func (m MiddleFn) ServeHttp(queryParams, pathParam Params, next Handler) (*Response, error) {
	return m(queryParams, pathParam, next)
}

type PreMiddleFn func(w http.ResponseWriter, r *http.Request, pathParams Params, next RawHandler)

func (m PreMiddleFn) ServeHttp(w http.ResponseWriter, r *http.Request, pathParams Params, next RawHandler) {
	m(w, r, pathParams, next)
}

type middlewareNode struct {
	jsonMiddleware JsonMiddleware
	preMiddleware  PreMiddleware
	next           *middlewareNode
}

func (m middlewareNode) HandleJson(queryParams, pathParam Params) (*Response, error) {
	return m.jsonMiddleware.ServeHttp(queryParams, pathParam, m.next.HandleJson)
}

func (m middlewareNode) HandleRaw(w http.ResponseWriter, r *http.Request, pathParams Params) {
	m.preMiddleware.ServeHttp(w, r, pathParams, m.next.HandleRaw)
}

func buildMiddlewareChain(wares []JsonMiddleware) middlewareNode {
	build := func() middlewareNode {
		fn := func(queryParams, pathParam Params, next Handler) (*Response, error) {
			return nil, nil
		}
		return middlewareNode{
			jsonMiddleware: MiddleFn(fn),
			next:           &middlewareNode{},
		}
	}
	var next middlewareNode
	if len(wares) == 0 {
		return build()
	} else if len(wares) == 1 {
		next = build()
	} else {
		next = buildMiddlewareChain(wares[1:])
	}
	return middlewareNode{
		jsonMiddleware: wares[0],
		next:           &next,
	}
}

func buildPreMiddlewareChain(wares []PreMiddleware) middlewareNode {
	build := func() middlewareNode {
		fn := func(w http.ResponseWriter, r *http.Request, pathParams Params, next RawHandler) {
		}
		return middlewareNode{
			preMiddleware: PreMiddleFn(fn),
			next:          &middlewareNode{},
		}
	}
	var next middlewareNode
	if len(wares) == 0 {
		return build()
	} else if len(wares) == 1 {
		next = build()
	} else {
		next = buildPreMiddlewareChain(wares[1:])
	}
	return middlewareNode{
		preMiddleware: wares[0],
		next:          &next,
	}
}
