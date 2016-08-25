package server

import (
	"encoding/json"
	"fmt"
	"github.com/julienschmidt/httprouter"
	"github.com/stretchr/graceful"
	"github.com/zhengyun1112/glrest/logger"
	"net/http"
	"net/url"
	"time"
)

type Handler func(queryParams, pathParam Params) (*Response, error)

type RawHandler func(w http.ResponseWriter, r *http.Request, pathParam Params)

type Server struct {
	srv         *graceful.Server
	router      *httprouter.Router
	namedRoutes map[string]string
	controllers []Controller
	jsonWares   []JsonMiddleware
	preWares    []PreMiddleware
	debug       bool
}

// Run the server listen and server at the addr with graceful shutdown supports.
func (s *Server) Run(addr string) error {
	timeout := 10 * time.Second
	if s.debug {
		timeout = 0
	}
	s.srv = &graceful.Server{
		Timeout: timeout,
		Server: &http.Server{
			Addr:    addr,
			Handler: s.router,
		},
	}
	logger.Info("Server is listening on %s", addr)
	return s.srv.ListenAndServe()
}

func (s *Server) AddController(c Controller) {
	s.controllers = append(s.controllers, c)
	c.setServer(s)
	c.RegisterRoutes()
}

func (s *Server) AddMiddleware(ware JsonMiddleware) {
	s.jsonWares = append(s.jsonWares, ware)
}

func (s *Server) AddPreMiddleware(ware PreMiddleware) {
	s.preWares = append(s.preWares, ware)
}

// Handle: basic interface which register a http request and handler to the router
func (s *Server) HandleJson(method, path string, handle Handler) {
	s.router.Handle(method, path, s.jsonAdapt(handle))
	//s.namedRoutes[name] = path
}

func (s *Server) HandleRaw(method, path string, handle RawHandler) {
	s.router.Handle(method, path, s.rawAdapt(handle))
}

// Get will register a 'GET' request handler to the router.
func (s *Server) Get(path string, handle Handler) {
	s.HandleJson("GET", path, handle)
}

// Post will register a 'POST' request handler to the router.
func (s *Server) Post(path string, handle Handler) {
	s.HandleJson("POST", path, handle)
}

// Put will register a 'PUT' request handler to the router.
func (s *Server) Put(path string, handle Handler) {
	s.HandleJson("PUT", path, handle)
}

// Patch will register a 'PATCH' request handler to the router.
func (s *Server) Patch(path string, handle Handler) {
	s.HandleJson("PATCH", path, handle)
}

// Head will register a 'HEAD' request handler to the router.
func (s *Server) Head(path string, handle Handler) {
	s.HandleJson("HEAD", path, handle)
}

// Delete will register a 'DELETE' request handler to the router.
func (s *Server) Delete(path string, handle Handler) {
	s.HandleJson("DELETE", path, handle)
}

// htAdapt adapts a sweb Handler to the httprouter Handle
func (s *Server) jsonAdapt(fn Handler) httprouter.Handle {
	return s.rawAdapt(s.rawToJsonAdapt(fn))
}

func (s *Server) rawAdapt(fn RawHandler) httprouter.Handle {
	core := func(w http.ResponseWriter, r *http.Request, pathParams Params, next RawHandler) {
		fn(w, r, pathParams)
	}
	return func(w http.ResponseWriter, r *http.Request, params httprouter.Params) {
		start := time.Now()
		pathParam := wrapRouterParam(params)
		mn := buildPreMiddlewareChain(append(s.preWares, PreMiddleFn(core)))
		mn.HandleRaw(w, r, pathParam)
		tm := time.Now().Sub(start).Seconds()
		logger.Info("URL: %s | Method : %s | Remote Addr : %s | Used Time : %.3f ms", r.RequestURI, r.Method, r.RemoteAddr, tm*1000)
	}
}

func (s *Server) rawToJsonAdapt(fn Handler) RawHandler {
	core := func(queryParams, pathParam Params, next Handler) (*Response, error) {
		return fn(queryParams, pathParam)
	}

	return func(w http.ResponseWriter, r *http.Request, pathParams Params) {
		var queryParams url.Values
		if r.Method == "GET" {
			queryParams = r.URL.Query()
		} else {
			r.ParseForm()
			queryParams = r.Form
		}
		qp := Params{Values: queryParams}
		mn := buildMiddlewareChain(append(s.jsonWares, MiddleFn(core)))
		resp, err := mn.HandleJson(qp, pathParams)
		if err != nil {
			serverError(w, err.Error(), http.StatusServiceUnavailable)
			logger.Error("server error, request: %s, err: %s", r.RequestURI, err.Error())
			//s.printHttpRequestLog(req, time.Now().Sub(start), err.Error())
			return
		} else {
			if !s.debug {
				// For non-debug mode, don't show dev message in response
				resp.DevMessage = ""
			}
			jsonStr, err := json.Marshal(resp)
			if err != nil {
				serverError(w, err.Error(), http.StatusServiceUnavailable)
				logger.Error("server error, request: %s, err: %s", r.RequestURI, err.Error())
				//s.printHttpRequestLog(req, time.Now().Sub(start), err.Error())
				return
			}
			w.Header().Set("Content-Type", "application/json;charset=UTF-8")
			respContent := string(jsonStr)
			fmt.Fprint(w, respContent)
		}
	}
}

func serverError(resp http.ResponseWriter, msg string, status int) {
	output := Response{
		Code:       RESPONSE_CODE_INTERNAL_ERROR,
		Message:    RESPONSE_MESSAGE_INTERNAL_ERROR,
		DevMessage: msg,
	}
	jsonStr, _ := json.Marshal(output)
	resp.Header().Set("Content-Type", "application/json;charset=UTF-8")
	respContent := string(jsonStr)
	http.Error(resp, respContent, status)
}

func wrapRouterParam(params httprouter.Params) Params {
	res := Params{Values: make(url.Values)}

	if len(params) > 0 {
		for _, param := range params {
			res.Set(param.Key, param.Value)
		}
	}
	return res
}

func New(isDebug bool) *Server {
	srv := &Server{
		router:    httprouter.New(),
		debug:     isDebug,
		jsonWares: make([]JsonMiddleware, 0),
		preWares:  make([]PreMiddleware, 0),
	}
	return srv
}
